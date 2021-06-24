terraform {
  required_version = ">= 0.13"
  required_providers {
    dotcommonitor = {
      source  = "rymancl/dotcommonitor"
      version = "~> 0.1"
    }
  }
}

provider "dotcommonitor" {
  uid = var.doctommonitor_uid
}