package provider_test

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the HASHICUPS_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
terraform {
	required_providers {
		turing-pi-bmc = {
			source = "jlec.de/dev/turing-pi-bmc"
		}
	}
}
provider "turing-pi-bmc" {
	endpoint = "turingpi"
}
`
)
