---
page_title: "Scheduler Data Source"
subcategory: "Scheduler"
---
# Data Source: dotcommonitor_scheduler
Represents a Dotcom-Monitor scheduler

## Example usage
```hcl
data "dotcommonitor_scheduler" "example" {
  name = "example-scheduler"
}
```

## Argument Reference
* `id` - **(Optional, int)** The ID of the scheduler. Must provide exactly one of `id`, `name`.
* `name` - **(Optional, string)** The exact name of the scheduler. This will fail if there exists more than one scheduler with the same name. Must provide exactly one of `id`, `name`.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `description` - The description of the scheduler.
* `weekly_interval` - List of weekly interval schedules.
* `excluded_time_interval` - List of excluded time interval schedules.
