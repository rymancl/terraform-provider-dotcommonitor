package dotcommonitor

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

const IgnoreErrorsCodesSeparator = ";"
const IgnoreErrorsCodesRangeSeparator = "-"

func resourceFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceFilterCreate,
		Read:   resourceFilterRead,
		Update: resourceFilterUpdate,
		Delete: resourceFilterDelete,
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
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_locations": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"num_tasks": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"num_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Default:      0,
						},
						"owner_device_down": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"ignore_errors": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							StateFunc:    StateToLower,
							ValidateFunc: validation.StringInSlice([]string{"Validation", "Runtime", "CustomScript", "Certificate", "Cryptographic", "Tcp", "Dns", "Udp", "Http", "Ftp", "Sftp", "Smtp", "Pop3", "Imap", "Icmp", "IcmpV6", "DnsBL", "Media", "Sip"}, true),
						},
						"codes": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIgnoreErrorsCodes,
						},
					},
				},
			},
		},
	}
}

func resourceFilterCreate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	api := meta.(*client.APIClient)

	filter := &client.Filter{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       expandFilterRules(d.Get("rules").(*schema.Set)),
		Items:       expandIgnoreErrors(d.Get("ignore_errors").(*schema.Set)),
	}
	log.Printf("[Dotcom-Monitor] Filter create configuration: %v", filter)

	// create the filter
	err := api.CreateFilter(filter)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to create filter: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Filter successfully created - ID: %v", fmt.Sprint(filter.ID))

	// Set ID
	strID := fmt.Sprint(filter.ID)
	d.SetId(strID)

	mutex.Unlock()
	return resourceFilterRead(d, meta)
}

func resourceFilterRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull filter ID from state
	filterID, _ := strconv.Atoi(d.Id())

	filter := &client.Filter{}
	filter.ID = filterID

	api := meta.(*client.APIClient)
	err := api.GetFilter(filter)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get filter: %s", err)
	}

	// Check if filter exists before trying to read it
	if !(filter.ID > 0) {
		log.Printf("[Dotcom-Monitor] [WARNING] Filter does not exist, removing ID %v from state", filter.ID)
		d.SetId("")
		return nil
	}

	// set state to detect drift
	d.Set("name", filter.Name)
	d.Set("description", filter.Description)
	d.Set("rules", flattenFilterRules(&filter.Rules))
	if filter.Items != nil {
		d.Set("ignore_errors", flattenIgnoreErrors(&filter.Items))
	}

	return nil
}

func resourceFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	d.Partial(true)

	// Pull filter ID from state
	filterID, _ := strconv.Atoi(d.Id())

	filter := &client.Filter{
		ID:          filterID,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       expandFilterRules(d.Get("rules").(*schema.Set)),
		Items:       expandIgnoreErrors(d.Get("ignore_errors").(*schema.Set)),
	}
	log.Printf("[Dotcom-Monitor] Attempting to update filter ID: %v", fmt.Sprint(filter.ID))

	api := meta.(*client.APIClient)
	err := api.UpdateFilter(filter)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to update filter: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Filter ID: %v successfully updated", fmt.Sprint(filter.ID))

	mutex.Unlock()
	d.Partial(false)
	return resourceFilterRead(d, meta)
}

func resourceFilterDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull filter ID from state
	filterID, _ := strconv.Atoi(d.Id())

	filter := &client.Filter{
		ID: filterID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteFilter(filter)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete filter: %s", err)
	}

	d.SetId("")

	return nil
}

//////////////////////////////
// Filter helpers
//////////////////////////////

// expandFilterRules ... constructs a dotcommonitor.Rule struct based on the set of rules in the TF configuration
func expandFilterRules(rules *schema.Set) client.Rule {
	ruleList := make([]client.Rule, len(rules.List()))

	for i, item := range rules.List() { // should only be 1 item
		var schemaMap = item.(map[string]interface{})

		ruleList[i] = client.Rule{
			TimeMin:           schemaMap["num_minutes"].(int),
			NumberOfLocations: schemaMap["num_locations"].(int),
			NumberOfTasks:     schemaMap["num_tasks"].(int),
			OwnerDevice:       schemaMap["owner_device_down"].(bool),
		}
	}

	if len(ruleList) == 1 {
		return ruleList[0]
	} else {
		return client.Rule{}
	}
}

// flattenFilterRules ... flattens filter rule objects to generic interface for state
func flattenFilterRules(rule *client.Rule) []map[string]interface{} {
	l := make([]map[string]interface{}, 0)

	m := make(map[string]interface{})
	m["num_minutes"] = &rule.TimeMin
	m["num_locations"] = &rule.NumberOfLocations
	m["num_tasks"] = &rule.NumberOfTasks
	m["owner_device_down"] = &rule.OwnerDevice
	l = append(l, m)

	return l
}

// expandIgnoreErrors ... constructs dotcommonitor.Items structs based on the set of rules in the TF configuration
func expandIgnoreErrors(items *schema.Set) []client.Item {
	itemList := make([]client.Item, len(items.List()))

	for i, item := range items.List() {
		var schemaMap = item.(map[string]interface{})

		itemList[i] = client.Item{
			ErrorType:         strings.ToLower(schemaMap["type"].(string)),
			ErrorCodeToIgnore: expandErrorCodesString(schemaMap["codes"].(string)),
		}
	}

	return itemList
}

// flattenIgnoreErrors ... flattens ignore errors objects to generic interface for state
func flattenIgnoreErrors(ignoreErrors *[]client.Item) []map[string]interface{} {
	l := make([]map[string]interface{}, 0)

	for _, item := range *ignoreErrors {
		m := make(map[string]interface{})
		m["type"] = strings.ToLower(item.ErrorType)
		m["codes"] = flattenIgnoreErrorsCodes(item.ErrorCodeToIgnore)

		l = append(l, m)
	}

	return l
}

// expandErrorCodesString ... expands an error codes string to a list of interfaces
func expandErrorCodesString(codes string) []interface{} {
	var l []interface{}

	parts := strings.Split(codes, IgnoreErrorsCodesSeparator)

	for _, item := range parts {
		// first see if current item is a single error code
		if i, err := strconv.Atoi(item); err != nil {
			// not a single code, try a range
			r := strings.Split(item, IgnoreErrorsCodesRangeSeparator)
			if len(r) == 2 {
				from, _ := strconv.Atoi(r[0])
				to, _ := strconv.Atoi(r[1])
				var errorRange = client.ErrorCodeToIgnoreRange{
					From: from,
					To:   to,
				}
				l = append(l, errorRange)
			}
		} else {
			// current item must be a single error code
			l = append(l, i)
		}
	}

	return l
}

// flattenIgnoreErrorsCodes ... flattens a list of interfaces into an error codes string
func flattenIgnoreErrorsCodes(codes []interface{}) string {
	var sb strings.Builder

	for _, item := range codes {
		// first see if current item is a single error code
		if i, ok := item.(float64); !ok {
			// not a single code, try a range
			var schemaMap = item.(map[string]interface{})
			r := client.ErrorCodeToIgnoreRange{
				From: int(schemaMap["From"].(float64)),
				To:   int(schemaMap["To"].(float64)),
			}
			sb.WriteString(strconv.Itoa(r.From))
			sb.WriteString(IgnoreErrorsCodesRangeSeparator)
			sb.WriteString(strconv.Itoa(r.To))
			sb.WriteString(IgnoreErrorsCodesSeparator)

		} else {
			// current item must be a single error code
			sb.WriteString(strconv.Itoa(int(i)))
			sb.WriteString(IgnoreErrorsCodesSeparator)
		}
	}

	return strings.Trim(sb.String(), IgnoreErrorsCodesSeparator)
}
