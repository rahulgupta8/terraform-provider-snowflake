# minimal
resource "snowflake_contact" "minimal" {
  name  = "example_contact"
  email = "admin@example.com"
}

# with all attributes set
resource "snowflake_contact" "complete" {
  name    = "production_alerts"
  email   = "alerts@example.com"
  comment = "Contact for production system alerts"
}
