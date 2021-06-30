---
page_title: "Task Data Source"
subcategory: "Task"
---
# Data Source: dotcommonitor_task
Represents a Dotcom-Monitor task

!>Please note that this data source cannot be used if there exists more than one resource with the same name under the specified device! The Dotcom-Monitor API supports `n` resources with the same name, but this becomes problematic when trying to target a specific resouce.

## Example usage
```hcl
data "dotcommonitor_task" "example" {
  name      = "example-task"
  device_id = dotcommonitor_device.example.id
}

resource "dotcommonitor_device" "example" {
  name      = "example-device"
  postpone  = true
  frequency = 60
}
```

## Argument Reference
* `name` - **(Required, string)** The exact name of the task.
* `device_id` - **(Required, int)** The ID of the device under which the task resides.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the device.
