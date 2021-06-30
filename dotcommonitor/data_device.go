package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func dataDevice() *schema.Resource {
	return &schema.Resource{
		Read: dataDeviceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"platform_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,                                         // ServerView
				ValidateFunc: validation.IntInSlice([]int{1, 3, 7, 12}), // 1=ServerView, 3=MetricsView, 7=BrowserView, 12=WebView
			},
		},
	}
}

func dataDeviceRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	var devices []client.Device
	api := meta.(*client.APIClient)

	platformID := d.Get("platform_id").(int)
	name := d.Get("name").(string)
	err := api.GetDevicesByName(platformID, name, &devices)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get devices by name: %s", err)
	}

	// No devices found for the given name on the platform
	if len(devices) < 1 {
		return fmt.Errorf("[Dotcom-Monitor] Query did not return any devices from API")
	}

	// We cannot process a situation where there is more than one device with the same name
	if len(devices) > 1 {
		// Get the list of ID's
		ids := make([]int, len(devices))
		for i, item := range devices {
			ids[i] = item.ID
		}

		return fmt.Errorf("[Dotcom-Monitor] Query returned %v devices from API for name %s on platform ID %v - "+
			"Device ID's returned: %v - "+
			"Devices must be updated to be unique in order to use this data source", len(devices), name, platformID, ids)
	}

	// If we get this far, we know we only got one device back from the API
	device := devices[0]
	log.Printf("[Dotcom-Monitor] Single device found: %v", &device.Name)

	if err1 := populateDeviceAttributes(d, device); err1 != nil {
		log.Printf("[Dotcom-Monitor] Error setting device attributes: %v", err1)
	}

	return nil
}

// populateDeviceAttributes ... fills in necessary schema attributes of the data source
func populateDeviceAttributes(d *schema.ResourceData, device client.Device) error {
	// Set ID
	strID := fmt.Sprint(device.ID)
	d.SetId(strID)

	return nil
}
