package test

import (
  "github.com/gruntwork-io/terratest/modules/terraform"
  "testing"
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

func TestTerraformAwsRoute53Zone(t *testing.T) {
  /* ARRANGE */
  options, err := toTerraformOptions("../../examples/route53_zone", &ModuleVars{
    ZoneName: "terratest.website.test.co.uk.",
    Records: []RecordVar{
      {"one", "A", []string{"10.0.0.0", "192.0.0.0"}, 60},
      {"two", "CNAME", []string{"dummy.armakuni.co.uk"}, 60},
    },
  })
  if err != nil {
    t.Fatalf("Failed to create terraform options: %s", err.Error())
  }
  terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

  ///* ACTION */
  terraform.InitAndPlan(t, terraformOptions)
  defer terraform.Destroy(t, terraformOptions)
  terraform.InitAndApply(t, terraformOptions)

  /* ASSERTIONS */
  zoneID := terraform.Output(t, terraformOptions, "zone_id")

  nameServers := GetRoute53HostedZoneNameServers(t, zoneID)
  if len(nameServers) < 1 {
    t.Errorf("No nameservers return for hosted zone")
    return
  }

  dnsServer := nameServers[0]

  lookupOne := FetchDNSRecords(t, "one.terratest.website.test.co.uk", dnsServer)
  lookupOne.AssertHasARecord("10.0.0.0")
  lookupOne.AssertHasARecord("192.0.0.0")

  lookupTwo := FetchDNSRecords(t, "two.terratest.website.test.co.uk", dnsServer)
  lookupTwo.AssertHasCNAMERecord("dummy.armakuni.co.uk.")
}
