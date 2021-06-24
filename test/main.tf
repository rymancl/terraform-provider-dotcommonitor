### DEVICE ###
resource "dotcommonitor_device" "device" {
  name     = "test-device"
  postpone = true
  #locations = "2,3,4,6,11,13,14,15,17,18,19,23,43,68,71,72,73,74,97,113,118,138,153,181,184,233"

  notifications_group {
      id = dotcommonitor_group.group.id 
      time_shift_min = 10
  }
}

# resource "dotcommonitor_device" "device" {
#     name = "test-device-renamed"
#     postpone = true
#     #locations = "2,3,4,6,11,13,14,15,17,18,19,23,43,68,71,72,73,74,97,113,118,138,153,181,184,233"

#     notifications_group {
#         id = data.dotcommonitor_group.dgroup.id
#         time_shift_min = 10
#     }
# }

# data "dotcommonitor_device" "ddevice" {
#     name = "test-device-renamed"
# }



### TASK ###
resource "dotcommonitor_task" "task" {
    device_id       = dotcommonitor_device.device.id
    request_type    = "GET"
    url             = "https://www.google.com"
    name            = "test-task"
    timeout         = 30

    # custom_dns_hosts = [
    #     {
    #     ip_address = "1.1.1.1"
    #     host = "myhost"
    # },
    #  {
    #     ip_address = "2.2.2.2"
    #     host = "myhost2"
    # }
    # ]
    # ssl_check_certificate_authority     = true
    # ssl_check_certificate_date          = true
    # ssl_check_certificate_cn            = true
    # ssl_check_certificate_usage         = true
    # ssl_check_certificate_revocation    = true
    # ssl_expiration_reminder_in_days     = 0
    # ssl_client_certificate              = "cert-name"

    # get_params = [
    #     {
    #         name = "paramname"
    #         value = "paramvalue"
    #     },
    #     {
    #         name = "paramname2"
    #         value = "paramvalue2"
    #     }
    # ]

    dns_resolve_mode = "External DNS Server"
    dns_server_ip = "1.1.1.1"
}



### GROUP ###
resource "dotcommonitor_group" "group" {
    name = "test-group"
    # addresses {
    #     type = "Email"
    #     address = "example@example.com"
    # }
    # addresses {
    #     type = "Phone"
    #     number = "5"
    #     code = "123"
    # }
    # addresses {
    #     type = "Pager"
    #     number = "9999999999999999"
    #     code = "123"
    #     message = "123456"
    # }
    # addresses {
    #     type = "Sms"
    #     number = "1235555555"
    # }
    # addresses {
    #     type = "Sms"
    #     number = "1239999999"
    # }
    # addresses {
    #     type = "PagerDuty"
    #     integration_key = "keygoeshere"
    # }
}

# data "dotcommonitor_group" "dgroup" {
#     name = "test-group"
# }