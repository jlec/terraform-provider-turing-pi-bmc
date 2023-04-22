
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
  endpoint = "10.100.100.231"
}

data "turing-pi-bmc_nodeinfo" "example" {
}

output "nodeinfo" {
  value = data.turing-pi-bmc_nodeinfo.example
}
