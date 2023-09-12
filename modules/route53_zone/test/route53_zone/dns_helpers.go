package test

import (
	"github.com/miekg/dns"
	"net"
	"testing"
)

type LookupResult struct {
	fqdn          string
	responseA     *dns.Msg
	responseCNAME *dns.Msg
	t             *testing.T
}

func FetchDNSRecords(t *testing.T, fqdn string, dnsServer string) *LookupResult {
	client := new(dns.Client)
	serverAddress := net.JoinHostPort(dnsServer, "53")

	responseA := performLookup(t, fqdn, dns.TypeA, client, serverAddress)
	if responseA == nil {
		return nil
	}

	responseCNAME := performLookup(t, fqdn, dns.TypeCNAME, client, serverAddress)
	if responseCNAME == nil {
		return nil
	}

	return &LookupResult{
		fqdn:          fqdn,
		responseA:     responseA,
		responseCNAME: responseCNAME,
		t:             t}
}

func performLookup(t *testing.T, fqdn string, recordType uint16, client *dns.Client, serverAddress string) *dns.Msg {
	query := new(dns.Msg)
	query.SetQuestion(dns.Fqdn(fqdn), recordType)
	response, _, err := client.Exchange(query, serverAddress)

	if err != nil {
		t.Errorf("DNS %s record lookup failed: %s", dns.TypeToString[recordType], err.Error())
		return nil
	}

	return response
}

func (lookup *LookupResult) AssertHasARecord(expectedARecord string) {
	if len(lookup.responseA.Answer) <= 0 {
		lookup.t.Errorf("DNS assertion failed: no answers found")
		return
	}

	found := false
	recordsString := ""

	for _, answer := range lookup.responseA.Answer {
		recordsString = recordsString + "\t" + answer.String() + "\n"

		if answer.(*dns.A).A.String() == expectedARecord {
			found = true
		}
	}

	if !found {
		lookup.t.Errorf(
			"DNS asserting failed: No A record with value %s found for %s.\nRecords Found:\n%s",
			expectedARecord,
			lookup.fqdn,
			recordsString,
		)
	}
}

func (lookup *LookupResult) AssertHasCNAMERecord(expectedTarget string) {
	if len(lookup.responseCNAME.Answer) <= 0 {
		lookup.t.Errorf("DNS assertion failed: no answers found")
		return
	}

	found := false
	recordsString := ""

	for _, answer := range lookup.responseCNAME.Answer {
		recordsString = recordsString + "\t" + answer.String() + "\n"

		if answer.(*dns.CNAME).Target == expectedTarget {
			found = true
		}
	}

	if !found {
		lookup.t.Errorf(
			"DNS asserting failed: No CNAME record with value %s found for %s.\nRecords Found:\n%s",
			expectedTarget,
			lookup.fqdn,
			recordsString,
		)
	}
}
