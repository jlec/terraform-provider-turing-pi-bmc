terraform {
  required_version = ">= 1.6.0"
  cloud {
    organization = "jlec-devops"

    workspaces {
      name = "terraform_provider_turing_pi_bmc"
    }
  }
}
