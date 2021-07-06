package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func dataScheduler() *schema.Resource {
	return &schema.Resource{
		Read: dataSchedulerRead,

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

func dataSchedulerRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	var schedulers []client.Scheduler
	api := meta.(*client.APIClient)

	id := d.Get("id").(int)
	name := d.Get("name").(string)

	// check which agrument was provided and make the appropriate API call
	if id != 0 {
		var scheduler client.Scheduler
		scheduler.ID = id
		err := api.GetScheduler(&scheduler)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get scheduler: %s", err)
		}
		schedulers = append(schedulers, scheduler)
	} else if name != "" {
		err := api.GetSchedulersByName(name, &schedulers)
		if err != nil {
			return fmt.Errorf("[Dotcom-Monitor] Failed to get schedulers by name: %s", err)
		}

		// We cannot process a situation where there is more than one scheduler with the same name
		if len(schedulers) > 1 {
			// Get the list of ID's
			ids := make([]int, len(schedulers))
			for i, item := range schedulers {
				ids[i] = item.ID
			}

			return fmt.Errorf("[Dotcom-Monitor] Query returned %v schedulers from API for name %s - "+
				"Scheduler ID's returned: %v - "+
				"Scheduler names must be unique in order to use this data source", len(schedulers), name, ids)
		}
	}

	// No schedulers found for this name
	if len(schedulers) < 1 {
		return fmt.Errorf("[Dotcom-Monitor] Query did not return any matching schedulers from the API")
	}

	// If we get this far, we know we only got one scheduler back from the API
	scheduler := schedulers[0]
	log.Printf("[Dotcom-Monitor] Single scheduler found: %v", &scheduler.Name)

	if err1 := populateSchedulerAttributes(d, scheduler); err1 != nil {
		log.Printf("[Dotcom-Monitor] Error setting scheduler attributes: %v", err1)
	}

	return nil
}

// populateSchedulerAttributes ... fills in necessary schema attributes of the data source
func populateSchedulerAttributes(d *schema.ResourceData, scheduler client.Scheduler) error {
	strID := fmt.Sprint(scheduler.ID)
	d.SetId(strID)
	d.Set("name", scheduler.Name)
	// TODO: Output more attributes?

	return nil
}
