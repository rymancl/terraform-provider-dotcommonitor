package dotcommonitor

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var mutex = &sync.Mutex{}

// Provider main
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOTCOM_MONITOR_UID", nil),
				Description: "Customer UID token",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dotcommonitor_task":      resourceTask(),
			"dotcommonitor_device":    resourceDevice(),
			"dotcommonitor_group":     resourceGroup(),
			"dotcommonitor_scheduler": resourceScheduler(),
			"dotcommonitor_filter":    resourceFilter(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"dotcommonitor_task":      dataTask(),
			"dotcommonitor_device":    dataDevice(),
			"dotcommonitor_group":     dataGroup(),
			"dotcommonitor_location":  dataLocation(),
			"dotcommonitor_locations": dataLocations(),
			"dotcommonitor_scheduler": dataScheduler(),
			"dotcommonitor_filter":    dataFilter(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		UID: d.Get("uid").(string),
	}

	return config.Client()
}
