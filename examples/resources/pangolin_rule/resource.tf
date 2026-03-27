resource "pangolin_rule" "example" {
  resource_id = pangolin_resource.example.id
  action      = "ALLOW"
  match       = "CIDR"
  value       = "100.0.0.0/24"
  priority    = 1
}

resource "pangolin_resource" "example" {
  org_id    = "your-org-id"
  name      = "Example App Resource"
  protocol  = "tcp"
  http      = true
  subdomain = "example-app"
  domain_id = "your-domain-id"
}
