package dotcommonitor

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

const schedulerExcludedTimeIntervalLayout = "2006-01-02T15:04Z"

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
						"from": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "0h0m",
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile("^([0-9]|[1][0-9]|[2][0-3])h([0-5]?[0-9])m$"), "must be in the format of [0-23]h[0-59]m"),
								validateWeeklyIntervalFrom(),
							),
						},
						"to": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "23h59m", // end of day
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile("^([0-9]|[1][0-9]|[2][0-3])h([0-5]?[0-9])m$"), "must be in the format of [0-23]h[0-59]m"),
								validateWeeklyIntervalTo(),
							),
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
						"from": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateExcludedTimeIntervalTimestamp,
						},
						"to": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateExcludedTimeIntervalTimestamp,
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
		WeeklyIntervals:       expandSchedulerWeeklyIntervalsList(d.Get("weekly_intervals").(*schema.Set)),
		ExcludedTimeIntervals: expandSchedulerExcludedTimeIntervalsList(d.Get("excluded_time_intervals").(*schema.Set)),
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
		d.Set("weekly_intervals", flattenSchedulerWeeklyIntervalsList(&scheduler.WeeklyIntervals))
	}
	if scheduler.ExcludedTimeIntervals != nil {
		d.Set("excluded_time_intervals", flattenSchedulerExcludedTimeIntervalsList(&scheduler.ExcludedTimeIntervals))
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
		WeeklyIntervals:       expandSchedulerWeeklyIntervalsList(d.Get("weekly_intervals").(*schema.Set)),
		ExcludedTimeIntervals: expandSchedulerExcludedTimeIntervalsList(d.Get("excluded_time_intervals").(*schema.Set)),
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

//////////////////////////////
// Scheduler helpers
//////////////////////////////

// expandSchedulerWeeklyIntervalsList ... constructs a list of dotcommonitor.WeeklyInterval structs based on the set of weekly_intervals in the TF configuration
func expandSchedulerWeeklyIntervalsList(weeklyIntervals *schema.Set) []client.WeeklyInterval {
	//log.Printf("[Dotcom-Monitor] Converting weekly_intervals list to dotcommonitor.WeeklyIntervals list")

	wiList := make([]client.WeeklyInterval, len(weeklyIntervals.List()))

	for i, item := range weeklyIntervals.List() {
		var schemaMap = item.(map[string]interface{})

		wiList[i] = client.WeeklyInterval{
			Days:       expandStringSet(schemaMap["days"].(*schema.Set)),
			FromMinute: convertDurationStringToMinutes(schemaMap["from"].(string)),
			ToMinute:   convertDurationStringToMinutes(schemaMap["to"].(string)),
			Enabled:    schemaMap["enabled"].(bool),
		}
	}

	return wiList
}

// flattenSchedulerWeeklyIntervalsList ... flattens weekly interval objects to generic interface for state
func flattenSchedulerWeeklyIntervalsList(weeklyIntervals *[]client.WeeklyInterval) []map[string]interface{} {
	l := make([]map[string]interface{}, 0)

	for _, item := range *weeklyIntervals {
		m := make(map[string]interface{})
		m["days"] = item.Days
		m["from"] = convertMinutesToDurationString(item.FromMinute)
		m["to"] = convertMinutesToDurationString(item.ToMinute)
		m["enabled"] = item.Enabled

		l = append(l, m)
	}

	return l
}

// expandSchedulerExcludedTimeIntervalsList ... constructs a list of dotcommonitor.DateTimeInterval structs based on the set of excluded_time_intervals in the TF configuration
func expandSchedulerExcludedTimeIntervalsList(excludedTimeIntervals *schema.Set) []client.DateTimeInterval {
	etList := make([]client.DateTimeInterval, len(excludedTimeIntervals.List()))

	for i, item := range excludedTimeIntervals.List() {
		var schemaMap = item.(map[string]interface{})

		etList[i] = client.DateTimeInterval{
			From: convertExcludedTimeIntervalFormatToUnix(schemaMap["from"].(string)),
			To:   convertExcludedTimeIntervalFormatToUnix(schemaMap["to"].(string)),
		}
	}

	return etList
}

// flattenSchedulerExcludedTimeIntervalsList ... flattens datetime interval objects to generic interface for state
func flattenSchedulerExcludedTimeIntervalsList(excludedTimeIntervals *[]client.DateTimeInterval) []map[string]interface{} {
	l := make([]map[string]interface{}, 0)

	for _, item := range *excludedTimeIntervals {
		m := make(map[string]interface{})
		m["from"] = convertUnixToExcludedTimeIntervalFormat(item.From)
		m["to"] = convertUnixToExcludedTimeIntervalFormat(item.To)
		l = append(l, m)
	}

	return l
}

// convertDurationStringToMinutes ... converts time duration string into minutes
func convertDurationStringToMinutes(s string) int {
	d, _ := time.ParseDuration(s)
	return int(d.Minutes())
}

// convertMinutesToDurationString ... converts minutes into time duration string
func convertMinutesToDurationString(i int) string {
	hours := i / 60
	mins := i % 60
	return fmt.Sprintf("%vh%vm", hours, mins)
}

// convertUnixToExcludedTimeIntervalFormat ... converts Unix epoch time to time format string
func convertUnixToExcludedTimeIntervalFormat(i int64) string {
	t := time.Unix(i/1000, 0).UTC() // API time is in milliseconds
	tf := t.Format(schedulerExcludedTimeIntervalLayout)
	return tf
}

// convertExcludedTimeIntervalFormatToUnix ... converts time format string to Unix epoch time
func convertExcludedTimeIntervalFormatToUnix(s string) int64 {
	u, _ := time.Parse(schedulerExcludedTimeIntervalLayout, s)
	return u.Unix() * 1000 // API time is in milliseconds
}
