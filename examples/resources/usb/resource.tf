terraform {
  required_version = "~>1.4"

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
