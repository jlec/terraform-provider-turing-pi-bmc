package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccSDCardDataSourceConfig = `
data "turing-pi-bmc_sdcard" "test" {
}
`

func TestAccSDCardDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: GetAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccSDCardDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.turing-pi-bmc_sdcard.test", "total", "0"),
				),
			},
		},
	})
}
