---
page_title: "Scheduler Resource"
subcategory: "Scheduler"
---
# Resource: dotcommonitor_scheduler
Represents a Dotcom-Monitor scheduler

## Example usage
```hcl
resource "dotcommonitor_scheduler" "example" {
  name = "example-scheduler"
  weekly_intervals {
    days    = ["Mo", "Tu"]
    from    = "1h00m"
    to      = "12h30m"
    enabled = true
  }
  weekly_intervals {
    days    = ["Th"]
    enabled = false
  }
  excluded_time_intervals {
    from = "2021-07-10T00:00:00Z"
    to   = "2021-07-11T24:00:00Z"
  } 
}

resource "dotcommonitor_device" "example" {
  name         = "example-device"
  scheduler_id = dotcommonitor_scheduler.example.id
  # other arguments
}
```

## Argument Reference
* `name` - **(Required, string)** The name of the scheduler.
* `description` - **(Optional, string)** The description of the scheduler.
* `weekly_intervals` - **(Optional, set{object})** Configuration block for a weekly interval schedule. Can be specified multiple times for each weekly interval. Each block supports the fields documented below.
* `excluded_time_intervals` - **(Optional, set{object})** Configuration block for an excluded time interval schedule. Can be specified multiple times for each excluded time interval. Each block supports the fields documented below.

### weekly_intervals
* `days` - **(Required, list{string})** The days the scheduler is active. Can be a list of any of "Su", "Mo", "Tu", "We", "Th", "Fr", "Sa".
* `from` - **(Optional, string)** The time of day when the scheduler becomes active. Must be in the format of `##h##m`. The input gets convered to minutes before being passed to the API. Defaults to "0h0m" (start of day).
* `to` - **(Optional, string)** The time of day when the scheduler turns inactive. Must be in the format of `##h##m`. The input gets convered to minutes before being passed to the API. Defaults to "24h0m" (end of day).
* `enabled` - **(Optional, bool)** Indicates if the scheduler is enabled.

### excluded_time_intervals
* `from` - **(Required, string)** The starting date/time during which monitoring should be excluded. Must be in "YYYY-MM-DDThh:mmZ" format in UTC/GMT only (for example, "2014-06-01T00:00Z"). The input gets converted to [Unix epoch](https://en.wikipedia.org/wiki/Unix_time) time before being passed to the API.
* `to` - **(Required, string)** The ending date/time during which monitoring should be excluded. Must be in "YYYY-MM-DDThh:mmZ" format in UTC/GMT only (for example, "2014-06-01T00:00Z"). The input gets converted to [Unix epoch](https://en.wikipedia.org/wiki/Unix_time) time before being passed to the API.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the scheduler.

## Import
`dotcommonitor_scheduler` can be imported using the ID of the scheduler, e.g.

```
$ terraform import dotcommonitor_scheduler.example 12345
```
