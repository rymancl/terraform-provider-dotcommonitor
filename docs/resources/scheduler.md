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
  weekly_interval {
    days        = ["Mo", "Tu"]
    from_minute = 5
    to_minute   = 90
    enabled     = true
  }
  weekly_interval {
    days        = ["Th"]
    from_minute = 30
    to_minute   = 60
    enabled     = false
  }
  excluded_time_interval {
    from_unix = 1358712000000
    to_unix   = 1358798400000
  } 
}

resource "dotcommonitor_device" "example" {
  name         = "example-device"
  scheduler_id = dotcommonitor_scheduler.example.id
}
```

## Argument Reference
* `name` - **(Required, string)** The name of the scheduler.
* `description` - **(Optional, string)** The description of the scheduler.
* `weekly_interval` - **(Optional, object)** Configuration block for a weekly interval schedule. Can be specified multiple times for each weekly interval. Each block supports the fields documented below.
* `excluded_time_interval` - **(Optional, object)** Configuration block for an excluded time interval schedule. Can be specified multiple times for each excluded time interval. Each block supports the fields documented below.

### weekly_interval
* `days` - **(Required, list{string})** The days the scheduler is active. Can be a list of any of "Su", "Mo", "Tu", "We", "Th", "Fr", "Sa".
* `from_minute` - **(Optional, int)** The minute of the day when the scheduler becomes active. Can be any int between 0 and 1439. Defaults to 0.
* `to_minute` - **(Optional, int)** The minute of the day when the scheduler turns inactive. Can be any int between 1 and 1440. Defaults to 1440.
* `enabled` - **(Optional, bool)** Indicates if the scheduler is enabled. Defaults to `true`.

### excluded_time_interval
* `from_unix` - **(Required, int)** The starting date/time during which monitoring should be excluded. Per the API, this must be a valid [Unix timestamp](https://en.wikipedia.org/wiki/Unix_time).
* `to_unix` - **(Required, int)** The ending date/time during which monitoring should be excluded. Per the API, this must be a valid [Unix timestamp](https://en.wikipedia.org/wiki/Unix_time).


## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the scheduler.

## Import
`dotcommonitor_scheduler` can be imported using the ID of the secheduler, e.g.

```
$ terraform import dotcommonitor_scheduler.example 6789
```
