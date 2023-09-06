package test

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	awsTerratest "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAwsS3WebsiteBucket(t *testing.T) {
	/* ARRANGE */
	t.Parallel()
	// Give this S3 Bucket a unique ID for a name tag so we can distinguish it from any other Buckets provisioned
	// in your AWS account
	expectedBucketName := fmt.Sprintf("terratest-website-bucket-test-%s", strings.ToLower(random.UniqueId()))
	// AWS region set in provider.tf or versions.tf
	expectedAwsRegion := "eu-west-2"
	expectBucketPublicBlock := "{\n  PublicAccessBlockConfiguration: {\n    BlockPublicAcls: false,\n    BlockPublicPolicy: false,\n    IgnorePublicAcls: false,\n    RestrictPublicBuckets: false\n  }\n}"
	expectedCanonicalUserPermissions := []string{"WRITE_ACP", "READ_ACP", "READ", "WRITE"}
	expectedGroupPermissions := []string{"READ", "READ_ACP", "WRITE"}
	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/website_bucket",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"name":   expectedBucketName,
			"region": expectedAwsRegion,
		},
	})

	/* ACTION */
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)
	// This will run `terraform init` and `terraform plan` and fail the test if there are any errors
	terraform.InitAndPlan(t, terraformOptions)
	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	/* ASSERTIONS */
	// Run `terraform output` to get the value of an output variable
	bucketID := terraform.Output(t, terraformOptions, "bucket_id")
	// Verify that our Bucket has versioning enabled
	actualStatus := awsTerratest.GetS3BucketVersioning(t, expectedAwsRegion, bucketID)
	expectedStatus := "Enabled"
	assert.Equal(t, expectedStatus, actualStatus)
	sess, err := session.NewSession(&aws.Config{
		Region: &expectedAwsRegion},
	)
	// Create S3 service client
	svc := s3.New(sess)
	//Verify that our Bucket have ACL and Permission matches the expected ones
	actualBucketACL, err := svc.GetBucketAcl(&s3.GetBucketAclInput{Bucket: aws.String(expectedBucketName)})
	if err != nil {
		exitErrorf("Unable to GetBucketAclInput, %v", err)
	}
	assert.NotEmpty(t, actualBucketACL)
	for _, b := range actualBucketACL.Grants {
		if *b.Grantee.Type == "CanonicalUser" {
			fmt.Println("Asserting CanonicalUser Permissions")
			assert.True(t, slices.Contains(expectedCanonicalUserPermissions, *b.Permission))
		} else {
			fmt.Println("Asserting Group Permissions")
			assert.True(t, slices.Contains(expectedGroupPermissions, *b.Permission))
		}
	}
	//Verify that our Bucket is Publicly Accessible
	actualPublicAccessBlock, err := svc.GetPublicAccessBlock(&s3.GetPublicAccessBlockInput{Bucket: aws.String(expectedBucketName)})
	if err != nil {
		exitErrorf("Unable to GetPublicAccessBlock, %v", err)
	}
	assert.EqualValues(t, expectBucketPublicBlock, actualPublicAccessBlock.String())
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
