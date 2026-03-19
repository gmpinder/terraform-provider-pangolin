terraform {
  required_providers {
    pangolin = {
      source = "registry.terraform.io/gmpinder/pangolin"
    }
  }
}

provider "pangolin" {
  token = "YOUR_API_TOKEN"
}

resource "pangolin_role" "admin" {
  org_id      = "your-org-id"
  name        = "Terraform Admin"
  description = "Managed by Terraform"
}

resource "pangolin_site_resource" "web_app" {
  org_id      = "your-org-id"
  name        = "Web App"
  mode        = "host"
  site_id     = 123
  destination = "webapp.internal"
  alias       = "webapp.your-domain.com"
  enabled     = true
  user_ids    = []
  role_ids    = [pangolin_role.admin.id]
  client_ids  = []
}

resource "pangolin_target" "web_app_target" {
  resource_id = pangolin_site_resource.web_app.id
  site_id     = 123
  ip          = "10.0.0.1"
  port        = 80
  enabled     = true
}
