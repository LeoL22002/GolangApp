terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}


# Configuración de proveedor
provider "digitalocean" {
  token = "dop_v1_bccd496d59a375becfce9808c850cfe76b1581e7eaeceef943f825f7fd88f76b"
}

resource "digitalocean_droplet" "mi_servidor" {
  name  = "mi-servidor"
  region  = "nyc1"
  image  = "ubuntu-20-04-x64"
  size   = "s-1vcpu-1gb"
}

# Salida
output "ip_servidor" {
  value = digitalocean_droplet.mi_servidor.ipv4_address
}