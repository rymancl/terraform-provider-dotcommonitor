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
			"platform_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,                                         // ServerView
				ValidateFunc: validation.IntInSlice([]int{1, 3, 7, 12}), // 1=ServerView, 3=MetricsView, 7=BrowserView, 12=WebView
			},
			"frequency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,                                                                                 // 5 minutes
				ValidateFunc: validation.IntInSlice([]int{60, 180, 300, 600, 900, 1800, 2700, 3600, 7200, 10800}), // in seconds
			},
			"locations": {
				Type:         schema.TypeList,
				Optional:     true,
				Elem: 		  &schema.Schema{Type: schema.TypeInt},
			},
			"avoid_simultaneous_checks": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"alert_silence_min": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"false_positive_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"send_uptime_alert": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"postpone": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"owner_device_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"filter_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scheduler_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"notifications_group": {
				Type:     schema.TypeList,
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
							Default:      0,
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
		NotificationGroups: constructNotificationsNotificationGroupList(d.Get("notifications_group").([]interface{})),
	}

	device := &client.Device{
		Name:                    d.Get("name").(string),
		PlatformID:              d.Get("platform_id").(int),
		Frequency:               d.Get("frequency").(int),
		Locations:               convertInterfaceListToIntList(d.Get("locations").([]interface{})),
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

	// Check if device exists before trying to read it
	if !doesDeviceExist(deviceID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Device does not exist, removing ID %v from state", deviceID)
		d.SetId("")
		return nil
	}

	device := &client.Device{}

	api := meta.(*client.APIClient)
	err := api.GetDevice(device)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get device: %s", err)
	}

	return nil
}

func resourceDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	// Pull device ID from state
	deviceID, _ := strconv.Atoi(d.Id())

	// Check if device exists before trying to update it
	if !doesDeviceExist(deviceID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Device does not exist, removing ID %v from state", deviceID)
		d.SetId("")
		return nil
	}

	notifications := &client.DeviceNotificationsBlock{
		NotificationGroups: constructNotificationsNotificationGroupList(d.Get("notifications_group").([]interface{})),
	}

	device := &client.Device{
		ID:                      deviceID,
		Name:                    d.Get("name").(string),
		PlatformID:              d.Get("platform_id").(int),
		Frequency:               d.Get("frequency").(int),
		Locations:               convertInterfaceListToIntList(d.Get("locations").([]interface{})),
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
	return resourceDeviceRead(d, meta)
}

func resourceDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull device ID from state
	deviceID, _ := strconv.Atoi(d.Id())

	// Check if device exists before trying to remove it
	if !doesDeviceExist(deviceID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Device does not exist, removing ID %v from state", deviceID)
		d.SetId("")
		return nil
	}

	device := &client.Device{
		ID: deviceID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteDevice(device)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete device: %s", err)
	}

	return nil
}

// doesDeviceExist ... determintes if a device with the given deviceID exists
func doesDeviceExist(deviceID int, meta interface{}) bool {
	log.Printf("[Dotcom-Monitor] [DEBUG] Checking if device exists with ID: %v", deviceID)
	device := &client.Device{
		ID: deviceID,
	}

	// Since an empty HTTP response is a valid 200 from the API, we will determine if
	//  the device exists by comparing the hash of the struct before and after the HTTP call.
	//  If the has does not change, it means nothing else was added, therefore it does not exist.
	//  If the hash changes, the API found the device and added the rest of the fields.
	h := sha256.New()
	t := fmt.Sprintf("%v", device)
	sum := h.Sum([]byte(t))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash before: %x", sum)

	// Try to get device from API
	api := meta.(*client.APIClient)
	err := api.GetDevice(device)

	t2 := fmt.Sprintf("%v", device)
	sum2 := h.Sum([]byte(t2))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash after: %x", sum2)

	// Compare the hashes, and if there was an error from the API we will assume the device exists
	//  to be safe that we do not improperly remove an existing device from state
	if bytes.Equal(sum, sum2) && err == nil {
		log.Println("[Dotcom-Monitor] [DEBUG] No new fields added to the device, therefore the device did not exist")
		return false
	}

	// If we get here, we can assume the device does exist
	return true
}
