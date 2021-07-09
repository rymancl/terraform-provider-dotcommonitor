package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func dataGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataGroupRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	var groups []client.Group
	api := meta.(*client.APIClient)

	name := d.Get("name").(string)
	err := api.GetGroupsByName(name, &groups)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get groups by name: %s", err)
	}

	// No groups found for the given name on the platform
	if len(groups) < 1 {
		return fmt.Errorf("[Dotcom-Monitor] Query did not return any groups from API")
	}

	// We cannot process a situation where there is more than one group with the same name
	if len(groups) > 1 {
		// Get the list of ID's
		ids := make([]int, len(groups))
		for i, item := range groups {
			ids[i] = item.ID
		}

		return fmt.Errorf("[Dotcom-Monitor] Query returned %v groups from API for name %s - "+
			"Group ID's returned: %v - "+
			"Groups must be updated to be unique in order to use this data source", len(groups), name, ids)
	}

	// If we get this far, we know we only got one group back from the API
	group := groups[0]
	log.Printf("[Dotcom-Monitor] Single group found: %v", &group.Name)

	if err1 := populateGroupAttributes(d, group); err1 != nil {
		log.Printf("[Dotcom-Monitor] Error setting group attributes: %v", err1)
	}

	return nil
}

// populateGroupAttributes ... fills in necessary schema attributes of the data source
func populateGroupAttributes(d *schema.ResourceData, group client.Group) error {
	// Set ID
	strID := fmt.Sprint(group.ID)
	d.SetId(strID)

	return nil
}
