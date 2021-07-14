---
page_title: "Alert Group Resource"
subcategory: "Alert Group"
---
# Resource: dotcommonitor_group
Represents a Dotcom-Monitor alert group

## Example usage
```hcl
resource "dotcommonitor_group" "example" {
  name = "example-group"
  addresses {
    type    = "Email"
    address = "test@test.com"
  }
  addresses {
    type   = "Phone"
    number = "5555555"
    code   = "123"
  }
  addresses {
    type   = "Sms"
    number = "1235555555"
  }
  addresses {
    type   = "Sms"
    number = "1239999999"
  }
  addresses {
    type            = "PagerDuty"
    integration_key = "key_goes_here"
  }
}
```

## Argument Reference
* `name` - **(Required, string)** The name of the alert group.
* `scheduler_id` - **(Optional, int)** The valid scheduler ID to use for the group.
* `addresses` - **(Optional, set{object})** Configuration block for an address. Can be specified multiple times for each address. Each block supports the fields documented below.

### addresses
* `type` - **(Required, string)** The type of address. Can be one of "Email", "Phone", "Sms", "PagerDuty".
* `template_id` - **(Optional, int)** The valid ID of the group template. Defaults to 0 (default template).
* `address` - **(Optional, string)** The address. Valid for "Email" `type` argument.
* `number` - **(Optional, string)** The number. Valid for "Phone", and "Sms" `type` argument.
* `code` - **(Optional, string)** The number code. Valid for "Phone" and `type` argument.
* `integration_key` - **(Optional, string)** The PagerDuty integration key. Valid for "PagerDuty" `type` argument.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the alert group


## Import
`dotcommonitor_group` can be imported using the ID of the alert group, e.g.

```
$ terraform import dotcommonitor_group.example 12345
```