---
page_title: "Device Resource"
subcategory: "Device"
---
# Resource: dotcommonitor_device
Represents a Dotcom-Monitor device

## Example usage
```hcl
resource "dotcommonitor_device" "example" {
  name      = "example-device"
  postpone  = true
  frequency = 60
}
```

## Argument Reference
* `name` - **(Required, string)** The name of the device.
* `platform_id` - **(Optional, int)**  The ID of the platform of the device. See [Monitoring Platforms](https://wiki.dotcom-monitor.com/knowledge-base-category/monitoring-platforms/) for more info. Note that [UserView is not supported](https://wiki.dotcom-monitor.com/knowledge-base/get-device-list-by-platform/) by API v.1. Can be one of 1 (ServerView), 3 (MetricsView), 7 (BrowserView). Defaults to 1.
* `frequency` - **(Optional, int)** The frequency that that the device checks at, in seconds. Can be one of 60, 180, 300, 600, 900, 1800, 2700, 3600, 7200, 10800. Defaults to 300.
* `locations` - **(Optional, list{int})** The list of location ID's for monitoring agents. Defined below. Note that if you do not provide this argument upon creation, the API will assign all locations and drift will be detected on subsequent plans.
* `avoid_simultaneous_checks` - **(Optional, bool)** Indicates if the device should avoid simultaneous checks.
* `alert_silence_min` - **(Optional, int)** The length of time alerts should be silenced, in minutes.
* `false_positive_check` - **(Optional, bool)** Indicates if the device should check for false positives (brief hiccup / network glitch). Dotcom-Monitor recommends having this enabled.
* `send_uptime_alert` - **(Optional, bool)** Indicates if uptime alerts should be sent when a device begins successfully completing tasks after a failure.
* `postpone` - **(Optional, bool)** Indicates if the device should be postponed/disabled.
* `owner_device_id` - **(Optional, int)** The valid device ID of the device that owns this device.
* `filter_id` - **(Optional, int)** The valid filter ID to use for the device.
* `scheduler_id` - **(Optional, int)** The valid scheduler ID to use for the device.
* `notifications_group` - **(Optional, object)** Configuration block for a notifications group. Can be specified multiple times for each notifications group. Note that groups can only be assigned to a device, you cannot assign a device to a group. Each block supports the fields documented below.

### locations
Can be any combination of valid public or private location ID's. This argument can be used in combination with the [locations data source](https://registry.terraform.io/providers/rymancl/dotcommonitor/latest/docs/data-sources/locations) or defined by providing ID's manully.

Public location list mapping:

Location ID | Location Name
--- | ---
1 	| Minneapolis 
2 	| New York 
3 	| London 
4 	| San Francisco 
6 	| Miami 
11 	| Hong Kong 
13 	| Montreal 
14 	| Frankfurt 
15 	| Denver 
17 	| Brisbane (premium)
18 	| Dallas 
19 	| Amsterdam 
23 	| Tel-Aviv 
43 	| Washington DC 
68 	| N. Virginia 
71 	| Tokyo (premium)
72 	| Shanghai 
73 	| Buenos Aires (premium)
74 	| Johannesburg 
97 	| Paris 
113 |	Warsaw 
118 |	Mumbai 
125 |	IPv6 San Franciso 
138 |	Seattle 
153 |	Copenhagen 
181 |	Sydney (premium)
184 |	Beijing 
233 |	Madrid 
445 |	Chengdu 
446 |	Guangzhou 
447 |	Qingdao 
448 |	Shenzhen 

_(last updated: July 2021)_

### notifications_group
* `id` - **(Required, int)** The ID of the alert group.
* `time_shift_min` - **(Optional, int)** The escalation time for the alert, in minutes. Can be one of 0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the device.

## Import
`dotcommonitor_device` can be imported using the ID of the device, e.g.

```
$ terraform import dotcommonitor_device.example device12345
```
