terraform {
  backend "remote" {
    organization = "ffddorf"

    workspaces {
      name = "terraform-backend"
    }
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~>3.1"
    }
  }
}
