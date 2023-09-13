package test

import (
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
