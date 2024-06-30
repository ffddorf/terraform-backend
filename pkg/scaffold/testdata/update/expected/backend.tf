terraform {
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
  backend "http" {
    address        = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend/default"
    lock_address   = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend/default"
    unlock_address = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend/default"
    username       = "github_pat"
  }
}
