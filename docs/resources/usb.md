---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "turing-pi-bmc_usb Resource - terraform-provider-turing-pi-bmc"
subcategory: ""
description: |-
  Interface to the USB CM4 port in the Turing PI 2 board. It allows switching between the host and device mode as well as mapping to one of the 4 nodes.
---

# turing-pi-bmc_usb (Resource)

Interface to the USB CM4 port in the Turing PI 2 board. It allows switching between the host and device mode as well as mapping to one of the 4 nodes.

## Example Usage

```terraform
terraform {
  required_providers {
    turing-pi-bmc = {
      source  = "jlec.de/dev/turing-pi-bmc"
      version = ">0"
    }
  }
}

provider "turing-pi-bmc" {
  endpoint = "turingpi"
}

resource "turing-pi-bmc_usb" "example" {
  node = 4
  mode = 0
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `mode` (String) USB port mode ('host' or 'device')
- `node` (Number) Node which USB port is mapped to

### Read-Only

- `id` (String) Unique identifier for this resource