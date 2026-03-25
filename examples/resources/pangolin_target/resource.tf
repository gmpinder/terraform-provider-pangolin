resource "pangolin_target" "example" {
  resource_id = pangolin_resource.example.id
  site_id     = pangolin_site.example.id
  ip          = "10.0.0.1"
  port        = 8080

  health_check = {
    enabled  = true
    hostname = "10.0.0.1"
    port     = 8080
  }
}

resource "pangolin_resource" "example" {
  org_id    = "your-org-id"
  name      = "Example App Resource"
  protocol  = "tcp"
  http      = true
  subdomain = "example-app"
  domain_id = "your-domain-id"
}

resource "pangolin_site" "example" {
  org_id = "your-org-id"
  name   = "example-site"
}
