# Dotcom-Monitor Provider
This is a Terraform provider for [Dotcom-Monitor](https://www.dotcom-monitor.com). You can read about their API [here](https://wiki.dotcom-monitor.com/knowledge-base/getting-started-with-the-api/).

->This provider only supports UID authentication, not legacy username/password authentication.

## Example Usage
Terraform 0.13 and later:
```hcl
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
  uid = "XXXXXXXXXXXXXXXXXXXXXXXXXXX"
}
```

## Argument Reference
* `uid` - **(Required, string)** The Dotcom-Monitor customer UID token. Can be specified via env variable `DOTCOM_MONITOR_UID`.
