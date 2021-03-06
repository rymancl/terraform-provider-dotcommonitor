package dotcommonitor

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeviceCreate,
		Read:   resourceDeviceRead,
		Update: resourceDeviceUpdate,
		Delete: resourceDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"locations": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"platform_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 7}), // 1=ServerView, 3=MetricsView, 7=BrowserView
			},
			"frequency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntInSlice([]int{60, 180, 300, 600, 900, 1800, 2700, 3600, 7200, 10800}), // in seconds
			},
			"avoid_simultaneous_checks": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"alert_silence_min": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"false_positive_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_uptime_alert": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"postpone": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"owner_device_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"filter_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scheduler_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"notifications_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"time_shift_min": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180}), // 0=immediate ... 10 mins - 3 hours
						},
					},
				},
			},
		},
	}
}

func resourceDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	api := meta.(*client.APIClient)

	notifications := &client.DeviceNotificationsBlock{
		NotificationGroups: expandNotificationsNotificationGroupList(d.Get("notifications_groups").(*schema.Set)),
	}

	device := &client.Device{
		Name:                    d.Get("name").(string),
		PlatformID:              d.Get("platform_id").(int),
		Frequency:               d.Get("frequency").(int),
		Locations:               expandIntSet(d.Get("locations").(*schema.Set)),
		AvoidSimultaneousChecks: d.Get("avoid_simultaneous_checks").(bool),
		AlertSilenceMin:         d.Get("alert_silence_min").(int),
		FalsePositiveCheck:      d.Get("false_positive_check").(bool),
		SendUptimeAlert:         d.Get("send_uptime_alert").(bool),
		Postpone:                d.Get("postpone").(bool),
		OwnerDeviceID:           d.Get("owner_device_id").(int),
		FilterID:                d.Get("filter_id").(int),
		SchedulerID:             d.Get("scheduler_id").(int),
		Notifications:           notifications,
	}
	log.Printf("[Dotcom-Monitor] device create configuration: %v", device)

	// create the device
	err := api.CreateDevice(device)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to create device: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Device successfully created - ID: %v", fmt.Sprint(device.ID))

	// Set ID
	strID := fmt.Sprint(device.ID)
	d.SetId(strID)

	mutex.Unlock()
	return resourceDeviceRead(d, meta)
}

func resourceDeviceRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull device ID from state
	deviceID, _ := strconv.Atoi(d.Id())

	device := &client.Device{}
	device.ID = deviceID

	api := meta.(*client.APIClient)
	err := api.GetDevice(device)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get device: %s", err)
	}

	// Check if device exists before trying to read it
	if !(device.ID > 0) {
		log.Printf("[Dotcom-Monitor] [WARNING] Device does not exist, removing ID %v from state", device.ID)
		d.SetId("")
		return nil
	}

	// set state to detect drift
	d.Set("name", device.Name)
	d.Set("platform_id", device.PlatformID)
	d.Set("frequency", device.Frequency)
	d.Set("locations", device.Locations)
	d.Set("avoid_simultaneous_checks", device.AvoidSimultaneousChecks)
	d.Set("alert_silence_min", device.AlertSilenceMin)
	d.Set("false_positive_check", device.FalsePositiveCheck)
	d.Set("send_uptime_alert", device.SendUptimeAlert)
	d.Set("postpone", device.Postpone)
	d.Set("owner_device_id", device.OwnerDeviceID)
	d.Set("filter_id", device.FilterID)
	d.Set("scheduler_id", device.SchedulerID)

	if device.Notifications != nil && device.Notifications.NotificationGroups != nil {
		d.Set("notifications_groups", device.Notifications.NotificationGroups)
	}

	return nil
}

func resourceDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	d.Partial(true)

	// Pull device ID from state
	deviceID, _ := strconv.Atoi(d.Id())

	notifications := &client.DeviceNotificationsBlock{
		NotificationGroups: expandNotificationsNotificationGroupList(d.Get("notifications_groups").(*schema.Set)),
	}

	device := &client.Device{
		ID:                      deviceID,
		Name:                    d.Get("name").(string),
		PlatformID:              d.Get("platform_id").(int),
		Frequency:               d.Get("frequency").(int),
		Locations:               expandIntSet(d.Get("locations").(*schema.Set)),
		AvoidSimultaneousChecks: d.Get("avoid_simultaneous_checks").(bool),
		AlertSilenceMin:         d.Get("alert_silence_min").(int),
		FalsePositiveCheck:      d.Get("false_positive_check").(bool),
		SendUptimeAlert:         d.Get("send_uptime_alert").(bool),
		Postpone:                d.Get("postpone").(bool),
		OwnerDeviceID:           d.Get("owner_device_id").(int),
		FilterID:                d.Get("filter_id").(int),
		SchedulerID:             d.Get("scheduler_id").(int),
		Notifications:           notifications,
	}

	log.Printf("[Dotcom-Monitor] Attempting to update device ID: %v", fmt.Sprint(device.ID))

	api := meta.(*client.APIClient)
	err := api.UpdateDevice(device)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to update device: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Device ID: %v successfully updated", fmt.Sprint(device.ID))

	mutex.Unlock()
	d.Partial(false)
	return resourceDeviceRead(d, meta)
}

func resourceDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull device ID from state
	deviceID, _ := strconv.Atoi(d.Id())

	device := &client.Device{
		ID: deviceID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteDevice(device)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete device: %s", err)
	}

	d.SetId("")

	return nil
}

//////////////////////////////
// Device helpers
//////////////////////////////

// expandNotificationsNotificationGroupList ... constructs a list of dotcommonitor.NotificationsNotificationGroups structs based on the set of notifications_group in the TF configuration
func expandNotificationsNotificationGroupList(notificationGroups *schema.Set) []client.NotificationsNotificationGroups {

	nnGroupList := make([]client.NotificationsNotificationGroups, len(notificationGroups.List()))

	for i, item := range notificationGroups.List() {
		var schemaMap = item.(map[string]interface{})

		nnGroupList[i] = client.NotificationsNotificationGroups{
			ID:           schemaMap["id"].(int),
			TimeShiftMin: schemaMap["time_shift_min"].(int),
		}
	}

	return nnGroupList
}
