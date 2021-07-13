package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/hashstructure"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func dataLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataLocationsRead,

		Schema: map[string]*schema.Schema{
			"all_locations": {
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{"all_locations", "all_public_locations", "all_private_locations", "ids", "names"},
			},
			"all_public_locations": {
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{"all_locations", "all_public_locations", "all_private_locations", "ids", "names"},
			},
			"all_private_locations": {
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{"all_locations", "all_public_locations", "all_private_locations", "ids", "names"},
			},
			"ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Elem: 		  &schema.Schema{Type: schema.TypeInt},
				ExactlyOneOf: []string{"all_locations", "all_public_locations", "all_private_locations", "ids", "names"},
			},
			"names": {
				Type:         schema.TypeSet,
				Optional:     true,
				Elem: 		  &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"all_locations", "all_public_locations", "all_private_locations", "ids", "names"},
			},
			"platform_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,                               // ServerView, currently only valid other than UserView which isn't supported by the API
				ValidateFunc: validation.IntInSlice([]int{1}), // 1=ServerView
			},
			"include_unavailable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"include_restrictive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func dataLocationsRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	var locations []client.Location
	api := meta.(*client.APIClient)

	all := d.Get("all_locations").(bool)
	allPublic := d.Get("all_public_locations").(bool)
	allPrivate := d.Get("all_private_locations").(bool)
	ids := expandIntSet(d.Get("ids").(*schema.Set))
	names := expandStringSet(d.Get("names").(*schema.Set))
	platformID := d.Get("platform_id").(int)
	includeUnavailable := d.Get("include_unavailable").(bool)
	includeRestrictive := d.Get("include_restrictive").(bool)

	// check which agrument was provided and make the appropriate API call
	if (all) {
		err := api.GetLocations(platformID, includeUnavailable, &locations)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get all locations: %s", err)
		}
	} else if (allPublic) {
		err := api.GetPublicLocations(platformID, includeUnavailable, &locations)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get all public locations: %s", err)
		}
	} else if (allPrivate) {
		err := api.GetPrivateLocations(platformID, includeUnavailable, &locations)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get all public locations: %s", err)
		}
	} else if (len(ids) > 0) {
		var allTemp []client.Location
		err := api.GetLocations(platformID, includeUnavailable, &allTemp)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get all locations: %s", err)
		}

		// check all locations to ensure ID is valid
		for _, item := range allTemp {
			if (!locationListContainsLocationID(allTemp, item.ID)) {
				return fmt.Errorf("[Dotcom-Monitor] No valid location return from API for ID: %v", item.ID)
			}
			locations = append(locations, item)
		}
	} else if (len(names) > 0) {
		var allTemp []client.Location
		err := api.GetLocations(platformID, includeUnavailable, &allTemp)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get all locations: %s", err)
		}

		// check all locations to ensure name is valid
		for _, item := range allTemp {
			if (!locationListContainsLocationName(allTemp, item.Name)) {
				return fmt.Errorf("[Dotcom-Monitor] No valid location return from API for name: %v", item.ID)
			}
			locations = append(locations, item)
		}
	}

	// No locations found  on the platform
	if len(locations) < 1 {
		return fmt.Errorf("[Dotcom-Monitor] Query did not return any locations from API")
	}

	// remove restrictive locations if requested
	if (!includeRestrictive) {
		locations = removeRestrictiveLocations(locations)
	}

	if err1 := populateLocationsAttributes(d, locations); err1 != nil {
		log.Printf("[Dotcom-Monitor] Error setting location attributes: %v", err1)
	}

	return nil
}

// populateLocationAttributes ... fills in necessary schema attributes of the data source
func populateLocationsAttributes(d *schema.ResourceData, locations []client.Location) error {
	hash, err := hashstructure.Hash(locations, nil)  // this may not generate a unique ID, but it is fine for data sources
	if err != nil {
		panic("[Dotcom-Monitor] Error hashing location data to create ID")
	}
	strHash := fmt.Sprint(hash)
	d.SetId(strHash)

	// fill ids and names
	ids := []int{}
	names := []string{}
	for _, item := range(locations) {
		ids = append(ids, item.ID)
		names = append(names, item.Name)
	}
	d.Set("ids", ids)
	d.Set("names", names)

	return nil
}


//////////////////////////////
// Location helpers
//////////////////////////////

// locationListContainsLocationID .. checks if provided location ID is valid in the list of locations
func locationListContainsLocationID(locations []client.Location, id int) bool {
    for _, item := range locations {
        if item.ID == id {
            return true
        }
    }
    return false
}

// locationListContainsLocationName .. checks if provided location name is valid in the list of locations
func locationListContainsLocationName(locations []client.Location, name string) bool {
    for _, item := range locations {
        if item.Name == name {
            return true
        }
    }
    return false
}

// removeRestrictiveLocations .. removes any locations that may be considered restrictive by
//  country-wide firewalls, government regulations, restrictions, etc.
//  This list can be updated as appropriate. 
func removeRestrictiveLocations(locations []client.Location) []client.Location {
	var restrictiveLocationIds = []int {11, 72, 184, 445, 446, 447, 448}
	// 11  = Hong Kong
	// 72  = Shanghai
	// 184 = Beijing
	// 445 = Chengdu
	// 446 = Guangzhou
	// 447 = Qingdao
	// 448 = Shenzhen

	var trimmedLocationList []client.Location
	for _, item := range locations { // iterate selected locations
		if !intInList(restrictiveLocationIds, item.ID) && !locationListContainsLocationID(trimmedLocationList, item.ID) {
			trimmedLocationList = append(trimmedLocationList, item)
		}
	}
	return trimmedLocationList
}
