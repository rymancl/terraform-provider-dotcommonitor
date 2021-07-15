package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func dataFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataFilterRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataFilterRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	var filters []client.Filter
	api := meta.(*client.APIClient)

	id := d.Get("id").(int)
	name := d.Get("name").(string)

	// check which agrument was provided and make the appropriate API call
	if id != 0 {
		var filter client.Filter
		filter.ID = id
		err := api.GetFilter(&filter)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get filter: %s", err)
		}
		filters = append(filters, filter)
	} else if name != "" {
		err := api.GetFiltersByName(name, &filters)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get filters by name: %s", err)
		}

		// We cannot process a situation where there is more than one filter with the same name
		if len(filters) > 1 {
			// Get the list of ID's
			ids := make([]int, len(filters))
			for i, item := range filters {
				ids[i] = item.ID
			}

			return fmt.Errorf("[Dotcom-Monitor] Query returned %v filters from API for name %s - "+
				"Filter ID's returned: %v - "+
				"Filter names must be unique in order to use this data source", len(filters), name, ids)
		}
	}

	// No filters found for this name
	if len(filters) < 1 {
		return fmt.Errorf("[Dotcom-Monitor] Query did not return any matching filters from the API")
	}

	// If we get this far, we know we only got one filter back from the API
	filter := filters[0]
	log.Printf("[Dotcom-Monitor] Single filter found: %v", &filter.Name)

	if err1 := populateFilterAttributes(d, filter); err1 != nil {
		log.Printf("[Dotcom-Monitor] Error setting filter attributes: %v", err1)
	}

	return nil
}

// populateFilterAttributes ... fills in necessary schema attributes of the data source
func populateFilterAttributes(d *schema.ResourceData, filter client.Filter) error {
	strID := fmt.Sprint(filter.ID)
	d.SetId(strID)
	d.Set("name", filter.Name)
	// TODO: Output more attributes?

	return nil
}
