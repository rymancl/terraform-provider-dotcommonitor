---
page_title: "Task Resource"
subcategory: "Task"
---
# Resource: dotcommonitor_task
Represents a Dotcom-Monitor task

## Example usage
```hcl
resource "dotcommonitor_task" "example" {
  device_id    = dotcommonitor_device.example.id
  request_type = "GET"
  url          = "https://www.google.com"
  name         = "example-task"
  
  get_params {
    name  = "paramname1"
    value = "paramvalue1"
  }
  get_params {
    name  = "paramname2"
    value = "paramvalue2"
  }

  custom_dns_hosts {
    ip_address = "1.1.1.1"
    host       = "myhost"
  }
  custom_dns_hosts {
    ip_address = "2.2.2.2"
    host       = "myhost2"
  }
}

resource "dotcommonitor_device" "example" {
  name = "example-device"
  # other arguments
}
```

## Argument Reference
* `url` - **(Required, string)** The url of the request of the task. Note: if using `get_params`, the API will append them to the end of the `url`. You may see a plan diff on the first plan after initial creation but this will not cause misconfiguration or drift.
* `name` - **(Required, string)** The name of the task.
* `device_id` - **(Required, int)** The valid ID of a device which to add the task to.
* `request_type` - **(Optional, string)** The type of request of the task. Can be one of "GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS", "TRACE", "PATCH". Defaults to "GET".
* `keyword1` - **(Optional, string)** The words or phrases that you wish to search for in the web page content. See [keyword validation](https://wiki.dotcom-monitor.com/knowledge-base/keyword-content-validation/) for more info.
* `keyword2` - **(Optional, string)** The words or phrases that you wish to search for in the web page content. See [keyword validation](https://wiki.dotcom-monitor.com/knowledge-base/keyword-content-validation/) for more info.
* `keyword3` - **(Optional, string)** The words or phrases that you wish to search for in the web page content. See [keyword validation](https://wiki.dotcom-monitor.com/knowledge-base/keyword-content-validation/) for more info.
* `username` - **(Optional, string)** The username to use for basic authentication.
* `userpass` - **(Optional, string)** The user password to use for basic authentication.
* `full_page_download` - **(Optional, bool)** Indicates if the task should download the full web page.
* `download_html` - **(Optional, bool)** Indicates if the task should download HTML.
* `download_frames` - **(Optional, bool)** Indicates if the task should download frames.
* `download_style_sheets` - **(Optional, bool)** Indicates if the task should download style sheets.
* `download_scripts` - **(Optional, bool)** Indicates if the task should download scripts.
* `download_images` - **(Optional, bool)** Indicates if the task should download images.
* `download_objects` - **(Optional, bool)** Indicates if the task should download objects.
* `download_applets` - **(Optional, bool)** Indicates if the task should download applets.
* `download_additional` - **(Optional, bool)** Indicates if the task should download additional content.
* `ssl_check_certificate_authority` - **(Optional, bool)** Indicates if the task should check the SSL certificate authority.
* `ssl_check_certificate_cn` - **(Optional, bool)** Indicates if the task should check the SSL certificate CN.
* `ssl_check_certificate_date` - **(Optional, bool)** Indicates if the task should check the SSL certificate date.
* `ssl_check_certificate_revocation` - **(Optional, bool)** Indicates if the task should check the SSL certificate revocation.
* `ssl_check_certificate_usage` - **(Optional, bool)** Indicates if the task should check the SSL certificate usage.
* `ssl_expiration_reminder_in_days` - **(Optional, int)** Sends an expiration alert X number of days prior to certificate expiration. Defaults to 0, meaning no expiration alert.
* `ssl_client_certificate` - **(Optional, string)** The name of the client certificate needed to access the site.
* `get_params` **(Optional, set{object})** Configuration block for GET request parameter. Can be specified multiple times for each parameter. Each block supports the fields documented below. Conflicts with `post_params`.
* `post_params` **(Optional, set{object})** Configuration block for POST request parameter. Can be specified multiple times for each parameter. Each block supports the fields documented below. Conflicts with `get_params`.
* `headers_params` **(Optional, set{object})** Configuration block for request header parameter. Can be specified multiple times for each parameter. Each block supports the fields documented below.
* `prepare_script` **(Optional, string)** The script contents to execute.
* `dns_resolve_mode` **(Optional, string)** The DNS resolve mode of the task. Can be one of "Device Cached", "Non Cached", "TTL Cached", "External DNS Server". Defaults to "Device Cached".
* `dns_server_ip` **(Optional, string)** The IP of a DNS server to use for the task.
* `custom_dns_hosts` **(Optional, list{object})** Configuration block for a custom DNS host. Can be specified multiple times for each custom DNS host. Each block supports the fields documented below.
* `task_type_id` **(Optional, int)** The ID of the task type to use for the task. See [ServerView documentation](https://wiki.dotcom-monitor.com/knowledge-base/serverview/) for valid task type ID's. Defaults to 2 (which is, HTTPS).
* `timeout` **(Optional, int)** The timeout value to use for the task, in seconds.

### get_params
Note: the API will append these parameters to the end of the `url`. You may see a plan diff on the first plan after initial creation but this will not cause misconfiguration or drift.

* `name` **(Required, string)** The name of the param.
* `value` **(Required, string)** The value of the param.

### post_params
* `name` **(Required, string)** The name of the param.
* `value` **(Required, string)** The value of the param.

### header_params
* `name` **(Required, string)** The name of the param.
* `value` **(Required, string)** The value of the param.

### `custom_dns_hosts` Configuration Block
* `ip_address` **(Required, string)** The IP address.
* `host` **(Required, string)** The host name.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the task.


## Import
`dotcommonitor_task` can be imported using the ID of the task, e.g.

```
$ terraform import dotcommonitor_task.example 12345
```
