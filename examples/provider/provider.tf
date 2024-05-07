terraform {
  required_providers {
    aidbox = {
      source  = "petalmd/aidbox"
      version = "~> 0.0.1"
    }
  }
}

provider "aidbox" {
  endpoint = "https://aidbox.app/rpc"
}

