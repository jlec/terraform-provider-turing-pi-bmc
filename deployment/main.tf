locals {
  tfe_oauth_name = "TF/Github/jlec"
  tfe_org_name   = "jlec"
}

data "tfe_oauth_client" "gh_oauth_token" {
  name         = local.tfe_oauth_name
  organization = local.tfe_org_name
}

resource "tfe_registry_module" "terraform-module" {
  vcs_repo {
    display_identifier = "jlec/terraform-provider-turing-pi-bmc"
    identifier         = "jlec/terraform-provider-turing-pi-bmc"
    oauth_token_id     = data.tfe_oauth_client.gh_oauth_token.oauth_token_id
  }
}
