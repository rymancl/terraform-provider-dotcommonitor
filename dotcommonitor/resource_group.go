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

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"scheduler_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Email", "Phone", "Pager", "Sms", "PagerDuty"}, true),
						},
						"template_id": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"number": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: groupAddressNumberIsValid(),
						},
						"code": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: groupAddressCodeIsValid(),
						},
						"message": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"integration_key": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	api := meta.(*client.APIClient)

	addresses := constructGroupAddresses(d.Get("addresses").([]interface{}))

	group := &client.Group{
		Name:        d.Get("name").(string),
		SchedulerID: d.Get("scheduler_id").(int),
		Addresses:   addresses,
	}
	log.Printf("[Dotcom-Monitor] group create configuration: %v", group)

	// create the group
	err := api.CreateGroup(group)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to create group: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Group successfully created - ID: %v", fmt.Sprint(group.ID))

	// Set ID
	strID := fmt.Sprint(group.ID)
	d.SetId(strID)

	mutex.Unlock()
	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull group ID from state
	groupID, _ := strconv.Atoi(d.Id())

	// Check if group exists before trying to read it
	if !doesGroupExist(groupID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Group does not exist, removing ID %v from state", groupID)
		d.SetId("")
		return nil
	}

	group := &client.Group{}

	api := meta.(*client.APIClient)
	err := api.GetGroup(group)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get group: %s", err)
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	// Pull group ID from state
	groupID, _ := strconv.Atoi(d.Id())

	// Check if group exists before trying to read it
	if !doesGroupExist(groupID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Group does not exist, removing ID %v from state", groupID)
		d.SetId("")
		return nil
	}

	addresses := constructGroupAddresses(d.Get("addresses").([]interface{}))

	group := &client.Group{
		ID:          groupID,
		Name:        d.Get("name").(string),
		SchedulerID: d.Get("scheduler_id").(int),
		Addresses:   addresses,
	}

	log.Printf("[Dotcom-Monitor] Attempting to update group ID: %v", fmt.Sprint(group.ID))

	api := meta.(*client.APIClient)
	err := api.UpdateGroup(group)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to update group: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Group ID: %v successfully updated", fmt.Sprint(group.ID))

	mutex.Unlock()
	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull group ID from state
	groupID, _ := strconv.Atoi(d.Id())

	// Check if group exists before trying to read it
	if !doesGroupExist(groupID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Group does not exist, removing ID %v from state", groupID)
		d.SetId("")
		return nil
	}

	group := &client.Group{
		ID: groupID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteGroup(group)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete group: %s", err)
	}

	return nil
}

// doesGroupExists ... determintes if a group with the given groupID exists
func doesGroupExist(groupID int, meta interface{}) bool {
	log.Printf("[Dotcom-Monitor] [DEBUG] Checking if group exists with ID: %v", groupID)
	group := &client.Group{
		ID: groupID,
	}

	// Since an empty HTTP response is a valid 200 from the API, we will determine if
	//  the group exists by comparing the hash of the struct before and after the HTTP call.
	//  If the has does not change, it means nothing else was added, therefore it does not exist.
	//  If the hash changes, the API found the group and added the rest of the fields.
	h := sha256.New()
	t := fmt.Sprintf("%v", group)
	sum := h.Sum([]byte(t))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash before: %x", sum)

	// Try to get group from API
	api := meta.(*client.APIClient)
	err := api.GetGroup(group)

	t2 := fmt.Sprintf("%v", group)
	sum2 := h.Sum([]byte(t2))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash after: %x", sum2)

	// Compare the hashes, and if there was an error from the API we will assume the group exists
	//  to be safe that we do not improperly remove an existing group from state
	if bytes.Equal(sum, sum2) && err == nil {
		log.Println("[Dotcom-Monitor] [DEBUG] No new fields added to the group, therefore the group did not exist")
		return false
	}

	// If we get here, we can assume the group does exist
	return true
}
