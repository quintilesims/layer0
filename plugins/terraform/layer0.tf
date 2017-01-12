
# Configure the Layer0 Provider
provider "layer0" {
  endpoint        = ""
  token           = ""
  skip_ssl_verify = true
}

# Create an environment named "demo"
resource "layer0_environment" "demo" {
  name = "demo"
}

# Create a load balancer named "guestbook" with port 80 exposed
resource "layer0_load_balancer" "guestbook" {
  name        = "guestbook"
  environment = "${layer0_environment.demo.id}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

# Show the load balancer's url as output
output "guestbook_url" {
  value = "${layer0_load_balancer.guestbook.url}"
}
