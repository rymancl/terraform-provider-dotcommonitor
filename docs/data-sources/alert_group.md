---
page_title: "Alert Group Data Source"
subcategory: "Alert Group"
---
# Data Source: dotcommonitor_group
Represents a Dotcom-Monitor alert group

!>Please note that this data source cannot be used if there exists more than one resource with the same name! The Dotcom-Monitor API supports `n` resources with the same name, but this becomes problematic when trying to target a specific resouce.

## Example usage
```hcl
data "dotcommonitor_group" "example" {
  name = "example-group"
}
```

## Argument Reference
* `name` - **(Required, string)** The exact name of the alert group.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the alert group.
