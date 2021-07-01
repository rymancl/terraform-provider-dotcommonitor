package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func dataLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataLocationRead,

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
			"platform_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,                               // ServerView, currently only valid other than UserView which isn't supported by the API
				ValidateFunc: validation.IntInSlice([]int{1}), // 1=ServerView
			},
			"private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataLocationRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	var locations []client.Location
	api := meta.(*client.APIClient)

	platformID := d.Get("platform_id").(int)
	id := d.Get("id").(int)
	name := d.Get("name").(string)

	// check which agrument was provided and make the appropriate API call
	if (id != 0) {
		var location client.Location
		err := api.GetLocation(platformID, id, &location)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get location by ID: %s", err)
		}
		locations = append(locations, location)
	} else if (name != "") {
		err := api.GetLocationsByName(platformID, name, true, &locations)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get location by name: %s", err)
		}

		// We cannot process a situation where there is more than one location with the same name
		if len(locations) > 1 {
			// Get the list of ID's
			ids := make([]int, len(locations))
			for i, item := range locations {
				ids[i] = item.ID
			}

			return fmt.Errorf("[Dotcom-Monitor] Query returned %v locations from API for name %s on platform ID %v - "+
				"Location ID's returned: %v - "+
				"Locations must be unique in order to use this data source", len(locations), name, platformID, ids)
		}
	}

	// No locations found  on the platform
	if len(locations) < 1 {
		return fmt.Errorf("[Dotcom-Monitor] Query did not return any locations from API")
	}

	// If we get this far, we know we only got one location back from the API
	location := locations[0]
	log.Printf("[Dotcom-Monitor] Single location found: %v", &location.Name)

	if err1 := populateLocationAttributes(d, location); err1 != nil {
		log.Printf("[Dotcom-Monitor] Error setting location attributes: %v", err1)
	}

	return nil
}

// populateLocationAttributes ... fills in necessary schema attributes of the data source
func populateLocationAttributes(d *schema.ResourceData, location client.Location) error {
	strID := fmt.Sprint(location.ID)
	d.SetId(strID)
	d.Set("name", location.Name)
	d.Set("private", location.IsPrivate)
	d.Set("available", location.Available)
	d.Set("deleted", location.IsDeleted)  // we don't return deleted locations from the client, but we will set this for transparency

	return nil
}
