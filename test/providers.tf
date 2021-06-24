terraform {
  required_version = "~> 1.0"
  required_providers {
    dotcommonitor = {
      source  = "github.com/rymancl/dotcommonitor"
      version = "~> 0.0.1"
    }
  }
}

provider "dotcommonitor" {
  uid = var.doctommonitor_uid
}