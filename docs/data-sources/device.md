---
page_title: "Device Data Source"
subcategory: "Device"
---
# Data Source: dotcommonitor_device
Represents a Dotcom-Monitor device

!>Please note that this data source cannot be used if there exists more than one resource with the same name! The Dotcom-Monitor API supports `n` resources with the same name, but this becomes problematic when trying to target a specific resouce.

## Example usage
```hcl
data "dotcommonitor_device" "example" {
  name = "example-device"
}

resource "dotcommonitor_task" "example" {
  device_id    = data.dotcommonitor_device.example.id
  request_type = "GET"
  url          = "https://www.google.com"
  name         = "example-task"
}
```

## Argument Reference
* `name` - **(Required, string)** The exact name of the device.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the device.
