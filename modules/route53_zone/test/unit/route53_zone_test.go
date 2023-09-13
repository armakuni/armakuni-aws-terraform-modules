package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type RecordVar struct {
	Name       string   `json:"name"`
	RecordType string   `json:"type"`
	Records    []string `json:"records"`
	Ttl        int      `json:"ttl"`
}

type ModuleVars struct {
	ZoneName string      `json:"zone_name"`
	Records  []RecordVar `json:"records"`
}

func toTerraformOptions(path string, options interface{}) (terraform.Options, error) {
	vars, err := structToMap(options)
	if err != nil {
		return terraform.Options{}, err
	}
	return terraform.Options{
		TerraformDir: path,
		Vars:         vars,
	}, nil
}

func TestRoute53ZoneWhenInvalidRecordTypeIsPassed(t *testing.T) {
	/* ARRANGE */
	options, err := toTerraformOptions("../../examples/route53_zone", &ModuleVars{
		ZoneName: "terratest.website.test.co.uk.",
		Records: []RecordVar{
			{"one", "ALPHA", []string{"10.0.0.0", "192.0.0.0"}, 60},
			{"two", "CNAME", []string{"dummy.armakuni.co.uk"}, 60},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create terraform options: %s", err.Error())
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	_, err = InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)

	/* ASSERTION */
	assert.ErrorContains(t, err, `Only valid types permitted (A, CNAME, MX, NS, TXT, SOA, SPF)`)
}

func structToMap(input interface{}) (map[string]interface{}, error) {
	var output map[string]interface{}
	jsonStr, err := json.Marshal(input)
	err = json.Unmarshal(jsonStr, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %s", err.Error())
	}
	return output, nil
}
