terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~>2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

variable "do_token" {
  type        = string
  description = "DigitalOcean API Token"
}

variable "ssh_key_path" {
  type        = string
  description = "Path to SSH key"
  default     = ""
}

output "ip_address" {
  value       = digitalocean_droplet.app.ipv4_address
  description = "IP address of droplet"
}

data "digitalocean_sizes" "main" {
  filter {
    key    = "vcpus"
    values = [1]
  }

  filter {
    key    = "regions"
    values = ["nyc1"]
  }

  sort {
    key       = "price_monthly"
    direction = "asc"
  }
}

resource "digitalocean_ssh_key" "main" {
  name       = "webserver"
  public_key = file(var.ssh_key_path)
}

resource "digitalocean_droplet" "app" {
  image      = "ubuntu-20-04-x64"
  name       = "webserver"
  region     = "nyc1"
  size       = element(data.digitalocean_sizes.main.sizes, 0).slug
  monitoring = true
  ssh_keys   = [digitalocean_ssh_key.main.fingerprint]
  user_data         = file("./cloud-config")
  graceful_shutdown = true
}
