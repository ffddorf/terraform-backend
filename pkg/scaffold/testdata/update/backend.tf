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
      version = "~>2.17"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~>3.1"
    }
  }
}
