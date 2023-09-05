package test

import (
	"fmt"
	"os"
	// "slices"
	// "strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAwsRoute53Zone(t *testing.T) {
	t.Parallel()

	expectedRoute53ZoneName := fmt.Sprintf("terratest.website.test.co.uk.")
	record1 := map[string]interface{}{
		"name":    "one",
		"records": []string{"10.0.0.0", "192.0.0.0"},
		"ttl":     60,
		"type":    "A",
	}

	record2 := map[string]interface{}{
		"name":    "two",
		"records": []string{"dummy.armakuni.co.uk"},
		"ttl":     60,
		"type":    "CNAME",
	}
	expectedRecords := []map[string]interface{}{record1, record2}

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/route53_zone",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"zone_name": expectedRoute53ZoneName,
			"records":   expectedRecords,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform plan` and fail the test if there are any errors
	terraform.InitAndPlan(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	zoneID := terraform.Output(t, terraformOptions, "zone_id")

	mySession := session.Must(session.NewSession())

	// Create a Route53 client from just a session.
	svc := route53.New(mySession)
	actualHostedZone, err := svc.GetHostedZone(&route53.GetHostedZoneInput{Id: aws.String(zoneID)})
	if err != nil {
		exitErrorf("Unable to GetHostedZone, %v", err)
	}
	// Assert that Route 53 Hosted Zone created with same Name
	assert.EqualValues(t, expectedRoute53ZoneName, *actualHostedZone.HostedZone.Name)
	// Assert that Route 53 Hosted Zone contains 4 Records (2 NS, 1 CNAME & 1 A)
	assert.EqualValues(t, 4, *actualHostedZone.HostedZone.ResourceRecordSetCount)
	listResourceRecordSets, err := svc.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(zoneID)})
	if err != nil {
		exitErrorf("Unable to ListResourceRecordSetsInput, %v", err)
	}
	for _, b := range listResourceRecordSets.ResourceRecordSets {
		if *b.Type == record1["type"] {
			expectedRecordName := fmt.Sprintf("%s%s%s", record1["name"],
				".", expectedRoute53ZoneName)
			assert.EqualValues(t, expectedRecordName, *b.Name)
			assert.EqualValues(t, record1["ttl"], *b.TTL)
		} else if *b.Type == record2["type"] {
			expectedRecordName := fmt.Sprintf("%s%s%s", record2["name"],
				".", expectedRoute53ZoneName)
			assert.EqualValues(t, expectedRecordName, *b.Name)
			assert.EqualValues(t, record2["ttl"], *b.TTL)
		} else {
			assert.True(t, (*b.Type == "NS" || *b.Type == "SOA"))
		}
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
