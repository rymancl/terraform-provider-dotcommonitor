---
page_title: "Locations Data Source"
subcategory: "Location"
---
# Data Source: dotcommonitor_locations
Retreives information for Dotcom-Monitor locations

->This data source does not return locations marked "IsDeleted" from the API.

## Example usage
### Basic example
```hcl
data "dotcommonitor_locations" "example" {
  all_locations = true
}

resource "dotcommonitor_device" "monitor_all" {
  name      = "example-device"
  postpone  = true
  frequency = 60
  locations = data.dotcommonitor_locations.example.ids
}
```

### Selecting locations by name
```hcl
data "dotcommonitor_locations" "example" {
  names = ["Seattle", "New York", "London", "Sydney"]
}

resource "dotcommonitor_device" "example" {
  name      = "example-device"
  postpone  = true
  frequency = 60
  locations = data.dotcommonitor_locations.example.ids
}
```

### Selecting only private agent locations
```hcl
data "dotcommonitor_locations" "example" {
  all_private_locations = true
}

resource "dotcommonitor_device" "monitor_private" {
  name      = "example-device"
  postpone  = true
  frequency = 60
  locations = data.dotcommonitor_locations.example.ids
}
```


## Argument Reference
* `all_locations` - **(Optional, bool)** Select all locations on the account, public and private. Must provide exactly one of `all_locations`, `all_public_locations`, `all_private_locations`, `ids`, `names`.
* `all_public_locations` - **(Optional, bool)** Select all public locations. Must provide exactly one of `all_locations`, `all_public_locations`, `all_private_locations`, `ids`, `names`.
* `all_private_locations` - **(Optional, bool)** Select all private agent locations. Must provide exactly one of `all_locations`, `all_public_locations`, `all_private_locations`, `ids`, `names`.
* `ids` - **(Optional, set{int})** List of location ID's to select. Must provide exactly one of `all_locations`, `all_public_locations`, `all_private_locations`, `ids`, `names`.
* `names` - **(Optional, set{string})** List of location names to select. Must provide exactly one of `all_locations`, `all_public_locations`, `all_private_locations`, `ids`, `names`.
* `platform_id` - **(Optional, int)** The ID of the platform. Locations are only supported in ServerView (1) and UserView, but since the API does not support UserView, only 1 is valid here.
* `include_unavailable` - **(Optional, bool)** Indicates whether or not to include locations not marked "Avilable" from the API. Defaults to `false`.
* `include_restrictive` - **(Optional, bool)** Indicates whether or not to include locations the provider restrictive by country-wide firewalls, government regulations, restrictions, etc. Defaults to `true`. Defined below.

### include_restrictive
Indicates whether or not to include locations the provider considers restrictive by country-wide firewalls, government regulations, restrictions, etc. Defaults to `true`.

The current locations on that list are:

Location ID | Location Name
--- | ---
11 	| Hong Kong 
72 	| Shanghai 
184 |	Beijing 
445 |	Chengdu 
446 |	Guangzhou 
447 |	Qingdao 
448 |	Shenzhen 

_(last updated: July 2021)_

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - Hash of the returned locations object. This should not be used.
