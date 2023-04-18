terraform {
  required_version = ">= 1.4.0"
  cloud {
    organization = "jlec"

    workspaces {
      name = "terraform_provider_turing_pi_bmc"
    }
  }
  required_providers {
    tfe = {
      source  = "hashicorp/tfe"
      version = "~>0.40"
    }
  }
}
provider "tfe" {
  # Configuration options
  token = var.tfe_token
}
