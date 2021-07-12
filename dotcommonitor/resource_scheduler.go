package dotcommonitor

import (
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"from_minute": {
							Type:         schema.TypeInt,
							Optional:     true,
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
						},
					},
				},
			},
			"excluded_time_intervals": {
				Type:     schema.TypeSet,
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
		WeeklyIntervals:       constructSchedulerWeeklyIntervalsList(d.Get("weekly_intervals").(*schema.Set)),
		ExcludedTimeIntervals: constructSchedulerExcludedTimeIntervalsList(d.Get("excluded_time_intervals").(*schema.Set)),
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

	scheduler := &client.Scheduler{}
	scheduler.ID = schedulerID

	api := meta.(*client.APIClient)
	err := api.GetScheduler(scheduler)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get scheduler: %s", err)
	}

	// Check if scheduler exists before trying to read it
	if !(scheduler.ID > 0) {
		log.Printf("[Dotcom-Monitor] [WARNING] Scheduler does not exist, removing ID %v from state", scheduler.ID)
		d.SetId("")
		return nil
	}

	// set state to detect drift
	d.Set("name", scheduler.Name)
	d.Set("description", scheduler.Description)
	if scheduler.WeeklyIntervals != nil {
		d.Set("weekly_intervals", scheduler.WeeklyIntervals)
	}
	if scheduler.ExcludedTimeIntervals != nil {
		d.Set("excluded_time_intervals", scheduler.ExcludedTimeIntervals)
	}

	return nil
}

func resourceSchedulerUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	d.Partial(true)

	// Pull scheduler ID from state
	schedulerID, _ := strconv.Atoi(d.Id())

	scheduler := &client.Scheduler{
		ID:                    schedulerID,
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		WeeklyIntervals:       constructSchedulerWeeklyIntervalsList(d.Get("weekly_intervals").(*schema.Set)),
		ExcludedTimeIntervals: constructSchedulerExcludedTimeIntervalsList(d.Get("excluded_time_intervals").(*schema.Set)),
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

	scheduler := &client.Scheduler{
		ID: schedulerID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteScheduler(scheduler)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete scheduler: %s", err)
	}

	d.SetId("")

	return nil
}
