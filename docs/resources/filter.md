---
page_title: "Filter Resource"
subcategory: "Filter"
---
# Resource: dotcommonitor_filter
Represents a Dotcom-Monitor filter

## Example usage
```hcl
resource "dotcommonitor_filter" "example" {
  name = "example-filter"
  rules {
    num_locations     = 3
    num_tasks         = 2
    owner_device_down = true
  }
  ignore_errors {
    type  = "http"
    codes = [301, 302]
  }
}

resource "dotcommonitor_device" "example" {
  name      = "example-device"
  filter_id = dotcommonitor_filter.example.id
  # other arguments
}
```

## Argument Reference
* `name` - **(Required, string)** The name of the filter.
* `rules` - **(Required, set{object})** Configuration block for a filter rule. Can be specified a maximum of one time. Each block supports the fields documented below.
* `description` - **(Optional, string)** The description of the filter.
* `ignore_errors` - **(Optional, set{object})** Configuration block for filter ignored errors. Can be specified multiple times for each ignored error. Each block supports the fields documented below.

### rules
* `num_locations` - **(Required, int)** The number of monitoring locations which are sending error responses. Must be at least 1.
* `num_tasks` - **(Required, int)** The number of failed taks. Must be at least 1.
* `num_minutes` - **(Optional, int)** The duration in minutes of the reported error. Defaults to 0.
* `owner_device_down` - **(Optional, bool)** Indicates if verification is enabled for if an owner device is down. Defaults to false.

### ignore_errors
* `type` - **(Required, string)** The ignored error type. Can be one of "Validation", "Runtime", "CustomScript", "Certificate", "Cryptographic", "Tcp", "Dns", "Udp", "Http", "Ftp", "Sftp", "Smtp", "Pop3", "Imap", "Icmp", "IcmpV6", "DnsBL", "Media", "Sip".
* `codes` - **(Required, set{int})** The ignored error codes. Must contain at least 1 code. Note: the web console suggests it supports ranges in addition to individual values. According to support (as of July 2021), this isn't actually true: only a list of single values are supported.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the filter.

## Import
`dotcommonitor_filter` can be imported using the ID of the filter, e.g.

```
$ terraform import dotcommonitor_filter.example 12345
```
