---
page_title: "Location Data Source"
subcategory: "Location"
---
# Data Source: dotcommonitor_location
Retreives information for a Dotcom-Monitor location

->This data source will not return a location marked "IsDeleted" from the API.

## Example usage
```hcl
data "dotcommonitor_location" "ny" {
  name = "New York"
}

output "location_ny" {
  value = data.dotcommonitor_location.ny
}
```
Output:
```
location_ny = {
  available   = true
  deleted     = false
  id          = 2
  name        = "New York"
  platform_id = 1
  private     = false
}
```

## Argument Reference
* `id` - **(Optional, int)** The ID of the location. Must provide exactly one of `id`, `name`.
* `name` - **(Optional, string)** The exact name of the location. Must provide exactly one of `id`, `name`.
* `platform_id` - **(Optional, int)** The ID of the platform. Locations are only supported in ServerView (1) and UserView, but since the API does not support UserView, only 1 is valid here.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the device.
* `private` - Indicates if the location is a private agent.
* `available` - Indicates if the location is marked as available.
* `deleted` - Indicates if the location is deleted. Currently the client does not return deleted locations.
