package dotcommonitor

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func resourceScheduler() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchedulerCreate,
		Read:   resourceSchedulerRead,
		Update: resourceSchedulerUpdate,
		Delete: resourceSchedulerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"weekly_intervals": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"from_minute": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 1439),
						},
						"to_minute": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1440, // end of day
							ValidateFunc: validation.IntBetween(1, 1440),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"excluded_time_intervals": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_unix": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"to_unix": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSchedulerCreate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	api := meta.(*client.APIClient)

	scheduler := &client.Scheduler{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		WeeklyIntervals:       constructSchedulerWeeklyIntervalsList(d.Get("weekly_intervals").([]interface{})),
		ExcludedTimeIntervals: constructSchedulerExcludedTimeIntervalsList(d.Get("excluded_time_intervals").([]interface{})),
	}
	log.Printf("[Dotcom-Monitor] Scheduler create configuration: %v", scheduler)

	// validate weekly interval days strings
	// this is done here since the provider plugin does not support TypeList validation
	for _, item := range scheduler.WeeklyIntervals {
		invalidDays := detectInvalidSchedulerWeeklyIntervalDays(item.Days)
		if len(invalidDays) > 0 {
			return fmt.Errorf("[Dotcom-Monitor] Invalid WeeklyInterval Days provided: %v", invalidDays)
		}
	}

	// create the scheduler
	err := api.CreateScheduler(scheduler)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to create scheduler: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Scheduler successfully created - ID: %v", fmt.Sprint(scheduler.ID))

	// Set ID
	strID := fmt.Sprint(scheduler.ID)
	d.SetId(strID)

	mutex.Unlock()
	return resourceSchedulerRead(d, meta)
}

func resourceSchedulerRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull scheduler ID from state
	schedulerID, _ := strconv.Atoi(d.Id())

	// Check if scheduler exists before trying to read it
	if !doesSchedulerExist(schedulerID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Scheduler does not exist, removing ID %v from state", schedulerID)
		d.SetId("")
		return nil
	}

	scheduler := &client.Scheduler{}
	api := meta.(*client.APIClient)
	err := api.GetScheduler(scheduler)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get scheduler: %s", err)
	}

	return nil
}

func resourceSchedulerUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	d.Partial(true)

	// Pull scheduler ID from state
	schedulerID, _ := strconv.Atoi(d.Id())

	// Check if scheduler exists before trying to remove it
	if !doesSchedulerExist(schedulerID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Scheduler does not exist, removing ID %v from state", schedulerID)
		d.SetId("")
		return nil
	}

	scheduler := &client.Scheduler{
		ID:                    schedulerID,
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		WeeklyIntervals:       constructSchedulerWeeklyIntervalsList(d.Get("weekly_intervals").([]interface{})),
		ExcludedTimeIntervals: constructSchedulerExcludedTimeIntervalsList(d.Get("excluded_time_intervals").([]interface{})),
	}

	// validate weekly interval days strings
	for _, item := range scheduler.WeeklyIntervals {
		invalidDays := detectInvalidSchedulerWeeklyIntervalDays(item.Days)
		if len(invalidDays) > 0 {
			return fmt.Errorf("[Dotcom-Monitor] Invalid WeeklyInterval Days provided: %v", invalidDays)
		}
	}

	log.Printf("[Dotcom-Monitor] Attempting to update scheduler ID: %v", fmt.Sprint(scheduler.ID))

	api := meta.(*client.APIClient)
	err := api.UpdateScheduler(scheduler)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to update scheduler: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Scheduler ID: %v successfully updated", fmt.Sprint(scheduler.ID))

	mutex.Unlock()
	d.Partial(false)
	return resourceSchedulerRead(d, meta)
}

func resourceSchedulerDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull scheduler ID from state
	schedulerID, _ := strconv.Atoi(d.Id())

	// Check if scheduler exists before trying to remove it
	if !doesSchedulerExist(schedulerID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Scheduler does not exist, removing ID %v from state", schedulerID)
		d.SetId("")
		return nil
	}

	scheduler := &client.Scheduler{
		ID: schedulerID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteScheduler(scheduler)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete scheduler: %s", err)
	}

	return nil
}

// doesSchedulerExist ... determintes if a scheduler with the given schedulerID exists
func doesSchedulerExist(schedulerID int, meta interface{}) bool {
	log.Printf("[Dotcom-Monitor] [DEBUG] Checking if scheduler exists with ID: %v", schedulerID)
	scheduler := &client.Scheduler{
		ID: schedulerID,
	}

	// Since an empty HTTP response is a valid 200 from the API, we will determine if
	//  the scheduler exists by comparing the hash of the struct before and after the HTTP call.
	//  If the has does not change, it means nothing else was added, therefore it does not exist.
	//  If the hash changes, the API found the scheduler and added the rest of the fields.
	h := sha256.New()
	t := fmt.Sprintf("%v", scheduler)
	sum := h.Sum([]byte(t))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash before: %x", sum)

	// Try to get scheduler from API
	api := meta.(*client.APIClient)
	err := api.GetScheduler(scheduler)

	t2 := fmt.Sprintf("%v", scheduler)
	sum2 := h.Sum([]byte(t2))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash after: %x", sum2)

	// Compare the hashes, and if there was an error from the API we will assume the scheduler exists
	//  to be safe that we do not improperly remove an existing scheduler from state
	if bytes.Equal(sum, sum2) && err == nil {
		log.Println("[Dotcom-Monitor] [DEBUG] No new fields added to the scheduler, therefore the scheduler did not exist")
		return false
	}

	// If we get here, we can assume the scheduler does exist
	return true
}
