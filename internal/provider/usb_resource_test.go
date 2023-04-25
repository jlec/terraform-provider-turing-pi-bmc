package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("device"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("TuringPiBMC_usb.test", "mode", "host"),
					resource.TestCheckResourceAttr("TuringPiBMC_usb.test", "node", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "TuringPiBMC_usb.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// usb code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"id", "usb"},
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("TuringPiBMC_usb.test", "node", "4"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "TuringPiBMC_usb" "test" {
  mode = %[1]q
}
`, configurableAttribute)
}
