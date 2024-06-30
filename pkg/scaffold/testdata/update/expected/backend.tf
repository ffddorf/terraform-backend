terraform {
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
  backend "http" {
    address        = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend.git/default"
    lock_address   = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend.git/default"
    unlock_address = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend.git/default"
    username       = "github_pat"
  }
}
