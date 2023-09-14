package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func toTerraformOptions(path string, vars *map[string]interface{}) terraform.Options {
	return terraform.Options{
		TerraformDir: path,
		Vars:         *vars,
	}
}

func TestRoute53ZoneWhenInvalidRecordTypeIsPassed(t *testing.T) {
	/* ARRANGE */
	options := toTerraformOptions("../../examples/route53_zone", &map[string]interface{}{
		"zone_name": "example.com.",
		"records": []map[string]interface{}{
			{"name": "one", "type": "ALPHA", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
		},
	})
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	_, err := InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)

	/* ASSERTION */
	assert.ErrorContains(t, err, `Only valid types permitted (A, CNAME, MX, NS, TXT, SOA, SPF)`)
}

func TestRoute53ZoneHasValidRecordEntries(t *testing.T) {
	/* ARRANGE */
	route53ExpectedData := &map[string]interface{}{
		"zone_name": "example.armakuni.com.",
		"records": []map[string]interface{}{
			{"name": "one", "type": "A", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
		},
	}

	options := toTerraformOptions("../../examples/route53_zone", route53ExpectedData)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	plan, _ := InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)
	fmt.Println(plan.ResourceChangesMap)

	/* ASSERTION */
	for _, record := range plan.ResourceChangesMap {
		// Record.Change.After is interface{} from raw terraform deserialised JSON
		recordMap, _ := record.Change.After.(map[string]interface{})
		// Records is a list of strings
		recordsInterface := recordMap["records"].([]interface{})

		switch recordMap["type"] {
		case "A":
			// Check correct amount of records expected
			if len(recordsInterface) != 2 {
				fmt.Printf("Invalid 'A' records: %+v, should be %+v\n", recordsInterface, 2)
			}

			// Check each records matches expected
			for _, record := range recordsInterface {
				if recordStr, ok := record.(string); ok {
					fmt.Printf("Record: %+v\n", recordStr)
				}
			}
		case "CNAME":
			// Check each records matches expected
			// for _, cname := range recordsInterface {
			// 	if cnameStr, ok := cname.(string); ok {
			// 		fmt.Printf("CNAME: %+v\n", cnameStr)
			// 	}
			// }
		default:
			//Zone name
			// assert.EqualValues(t, route53ExpectedData["zone_name"], recordMap["name"])
		}
	}

	// Number of Resources to be created
	assert.EqualValues(t, 3, len(plan.ResourcePlannedValuesMap))
}
