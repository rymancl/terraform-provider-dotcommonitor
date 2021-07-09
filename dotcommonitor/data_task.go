package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func dataTask() *schema.Resource {
	return &schema.Resource{
		Read: dataTaskRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device_id": {
				Type:         schema.TypeInt,
				Required:     true,
			},
		},
	}
}

func dataTaskRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	var tasks []client.Task
	api := meta.(*client.APIClient)

	deviceID := d.Get("device_id").(int)
	name := d.Get("name").(string)
	err := api.GetDeviceTasksByName(deviceID, name, &tasks)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get task by name: %s", err)
	}

	// No tasks found for the given name on the device
	if len(tasks) < 1 {
		return fmt.Errorf("[Dotcom-Monitor] Query did not return any tasks from API")
	}

	// We cannot process a situation where there is more than one task with the same name
	if len(tasks) > 1 {
		// Get the list of ID's
		ids := make([]int, len(tasks))
		for i, item := range tasks {
			ids[i] = item.ID
		}

		return fmt.Errorf("[Dotcom-Monitor] Query returned %v tasks from API for name %s on device ID %v - "+
			"Task ID's returned: %v - "+
			"Tasks must be updated to be unique in order to use this data source", len(tasks), name, deviceID, ids)
	}

	// If we get this far, we know we only got one device back from the API
	task := tasks[0]
	log.Printf("[Dotcom-Monitor] Single task found: %v", &task.Name)

	if err1 := populateTaskAttributes(d, task); err1 != nil {
		log.Printf("[Dotcom-Monitor] Error setting task attributes: %v", err1)
	}

	return nil
}

// populateTaskAttributes ... fills in necessary schema attributes of the data source
func populateTaskAttributes(d *schema.ResourceData, task client.Task) error {
	// Set ID
	strID := fmt.Sprint(task.ID)
	d.SetId(strID)

	return nil
}
