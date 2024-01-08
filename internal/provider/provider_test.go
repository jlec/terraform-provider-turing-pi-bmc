package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jlec/terraform-provider-turing-pi-bmc/internal/provider"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
//
//	var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
//		"TuringPiBMC": providerserver.NewProtocol6WithError(provider.New("test")()),
//	}
//
// Instead of a global variable, use a function.
func GetAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"TuringPiBMC": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

func testAccPreCheck(_ *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
