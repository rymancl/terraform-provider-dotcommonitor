package dotcommonitor

import (
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
							ValidateFunc: validation.IntAtLeast(0),
						},
						"address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"number": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 16),
								validateGroupAddressNumber(),
							),
						},
						"code": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(3, 3),
								validateGroupAddressCode(),
							),
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

	addresses := expandGroupAddresses(d.Get("addresses").([]interface{}))

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

	group := &client.Group{}
	group.ID = groupID

	api := meta.(*client.APIClient)
	err := api.GetGroup(group)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get group: %s", err)
	}

	// Check if group exists before trying to read it
	if !(group.ID > 0) {
		log.Printf("[Dotcom-Monitor] [WARNING] Group does not exist, removing ID %v from state", group.ID)
		d.SetId("")
		return nil
	}

	// set state to detect drift
	d.Set("name", group.Name)
	d.Set("scheduler_id", group.SchedulerID)

	if group.Addresses != nil {
		d.Set("addresses", group.Addresses)
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	d.Partial(true)

	// Pull group ID from state
	groupID, _ := strconv.Atoi(d.Id())

	addresses := expandGroupAddresses(d.Get("addresses").([]interface{}))

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
	d.Partial(false)
	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull group ID from state
	groupID, _ := strconv.Atoi(d.Id())

	group := &client.Group{
		ID: groupID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteGroup(group)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete group: %s", err)
	}

	d.SetId("")

	return nil
}

//////////////////////////////
// Group helpers
//////////////////////////////

// expandGroupAddresses .. constructs a list of client.Addresses structs based on the addresses schema in the TF configuration
func expandGroupAddresses(schemaAddresses []interface{}) []client.Addresses {
	addressList := make([]client.Addresses, len(schemaAddresses))

	for i, item := range schemaAddresses {
		var schemaMap = item.(map[string]interface{})

		addressList[i] = client.Addresses{
			Type:       schemaMap["type"].(string),
			TemplateID: schemaMap["template_id"].(int),
		}

		// Populate rest of the struct with the appropriate fields
		switch addressList[i].Type {
		case "Email":
			addressList[i].Address = schemaMap["address"].(string)
		case "Phone":
			addressList[i].Number = schemaMap["number"].(string)
			addressList[i].Code = schemaMap["code"].(string)
		case "Pager":
			addressList[i].Number = schemaMap["number"].(string)
			addressList[i].Code = schemaMap["code"].(string)
			addressList[i].Message = schemaMap["message"].(string)
		case "Sms":
			addressList[i].Number = schemaMap["number"].(string)
		case "PagerDuty":
			addressList[i].IntegrationKey = schemaMap["integration_key"].(string)
		case "Script":
			addressList[i].Message = schemaMap["message"].(string)
		}
	}

	return addressList
}
