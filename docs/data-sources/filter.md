---
page_title: "Filter Data Source"
subcategory: "Filter"
---
# Data Source: dotcommonitor_filter
Represents a Dotcom-Monitor filter

## Example usage
```hcl
data "dotcommonitor_filter" "example" {
  name = "example-filter"
}
```

## Argument Reference
* `id` - **(Optional, int)** The ID of the filter. Must provide exactly one of `id`, `name`.
* `name` - **(Optional, string)** The exact name of the filter. This will fail if there exists more than one filter with the same name. Must provide exactly one of `id`, `name`.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* _N/A_

