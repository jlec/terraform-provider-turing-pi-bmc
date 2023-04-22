package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccNodeInfoDataSourceConfig = `
data "turing-pi-bmc_nodeinfo" "test" {
}
`

func TestAccNodeInfoDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccNodeInfoDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.turing-pi-bmc_nodeinfo.test", "node1", "unknown"),
				),
			},
		},
	})
}
