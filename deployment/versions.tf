terraform {
  required_version = ">= 1.4.0"
  cloud {
    organization = "jlec"

    workspaces {
      name = "terraform_provider_turing_pi_bmc"
    }
  }
}
