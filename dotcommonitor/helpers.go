package dotcommonitor

import (
	"bytes"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

//////////////////////////////
// Common helpers
//////////////////////////////

// convertStringListToIntList ... type asserting a list of string to a list of int
func convertStringListToIntList(csvString string) []int {
	//log.Printf("[Dotcom-Monitor] convertStringListToIntList - Converting %v to int[]", csvString)
	sp := strings.Split(csvString, ",")
	intList := make([]int, len(sp))

	for i, item := range sp {
		intList[i], _ = strconv.Atoi(item)
	}

	return intList
}

// convertInterfaceListToIntList ... type asserting a list of interfaces to a list of int
func convertInterfaceListToIntList(interfaceList []interface{}) []int {
	//log.Printf("[Dotcom-Monitor] convertInterfaceListToIntList - Converting %v to int[]", interfaceList)
	intList := make([]int, len(interfaceList))

	for i := range interfaceList {
		intList[i] = interfaceList[i].(int)
	}

	return intList
}

// convertInterfaceListToStringList ... type asserting a list of interfaces to a list of string
func convertInterfaceListToStringList(interfaceList []interface{}) []string {
	//log.Printf("[Dotcom-Monitor] convertInterfaceListToStringList - Converting %v to int[]", interfaceList)
	stringList := make([]string, len(interfaceList))

	for i := range interfaceList {
		stringList[i] = interfaceList[i].(string)
	}

	return stringList
}

// expandIntSet ... type asserting a set to a list of int
func expandIntSet(set *schema.Set) []int {
	return convertInterfaceListToIntList(set.List())
}

// intInList .. checks if the int is in the list of ints
func intInList(intList []int, num int) bool {
	sort.Ints(intList)
	index := sort.Search(len(intList), func(i int) bool { return intList[i] >= num })
	result := (index < len(intList)) && (intList[index] == num)
	return result
}

// stringInList .. checks if the string is in the list of strings
func stringInList(stringList []string, s string) bool {
	sort.Strings(stringList)
	index := sort.Search(len(stringList), func(i int) bool { return stringList[i] >= s })
	result := (index < len(stringList)) && (stringList[index] == s)
	return result
}

//////////////////////////////
// Task helpers
//////////////////////////////

// expandInterfaceListToTaskParamList ... type asserting a list of interfaces to a list of TaskParam
func expandInterfaceListToTaskParamList(schemaInterfaceList []interface{}) []client.TaskParam {
	taskParamList := make([]client.TaskParam, len(schemaInterfaceList))

	for i, item := range schemaInterfaceList {
		var schemaMap = item.(map[string]interface{})
		taskParamList[i] = client.TaskParam{
			Name:  schemaMap["name"].(string),
			Value: schemaMap["value"].(string),
		}
		//log.Printf("[Dotcom-Monitor] [expandInterfaceListToTaskParamList] Added TaskParam to list - Name: %v  Value: %v", taskParamList[i].Name, taskParamList[i].Value)
	}
	return taskParamList
}

// flattenCustomDnsHostsToString ... returns a string required for the syntax of "CustomDNSHosts"
//  Syntax:  <host>=<ip>;
func flattenCustomDnsHostsToString(hosts []interface{}) string {
	buf := bytes.Buffer{}

	for _, item := range hosts {
		var schemaMap = item.(map[string]interface{})
		var ip = schemaMap["ip_address"].(string)
		var host = schemaMap["host"].(string)

		buf.WriteString(host)
		buf.WriteString("=")
		buf.WriteString(ip)
		buf.WriteString(";")

	}
	resultString := buf.String()
	log.Printf("[Dotcom-Monitor] [flattenCustomDnsHostsToString] Result string: %s", resultString)

	return resultString
}

//////////////////////////////
// Device helpers
//////////////////////////////

// constructNotificationsNotificationGroupList ... constructs a list of dotcommonitor.NotificationsNotificationGroups structs based on the list of notifications_group in the TF configuration
func constructNotificationsNotificationGroupList(notificationGroups []interface{}) []client.NotificationsNotificationGroups {
	//log.Printf("[Dotcom-Monitor] Converting notifications_group list to dotcommonitor.NotificationsNotificationGroups list")

	nnGroupList := make([]client.NotificationsNotificationGroups, len(notificationGroups))

	for i, item := range notificationGroups {
		var schemaMap = item.(map[string]interface{})

		nnGroupList[i] = client.NotificationsNotificationGroups{
			ID:           schemaMap["id"].(int),
			TimeShiftMin: schemaMap["time_shift_min"].(int),
		}
		//log.Printf("[Dotcom-Monitor] [constructNotificationsNotificationGroupList] Added NotificationGroup to list - ID: %v  TimeShiftMin: %v", nnGroupList[i].ID, nnGroupList[i].TimeShiftMin)
	}

	return nnGroupList
}



//////////////////////////////
// Group helpers
//////////////////////////////

// constructGroupAddresses .. constructs a list of client.Addresses structs based on the addresses schema in the TF configuration
func constructGroupAddresses(schemaAddresses []interface{}) []client.Addresses {
	//log.Printf("[Dotcom-Monitor] Converting addresses interface list to dotcommonitor.Addresses list")

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

		//log.Printf("[Dotcom-Monitor] [constructNotificationsNotificationGroupList] Added an address to list with Type: %v", addressList[i].Type)
	}

	return addressList
}

//////////////////////////////
// Location helpers
//////////////////////////////

// locationListContainsLocationID .. checks if provided location ID is valid in the list of locations
func locationListContainsLocationID(locations []client.Location, id int) bool {
    for _, item := range locations {
        if item.ID == id {
            return true
        }
    }
    return false
}

// locationListContainsLocationName .. checks if provided location name is valid in the list of locations
func locationListContainsLocationName(locations []client.Location, name string) bool {
    for _, item := range locations {
        if item.Name == name {
            return true
        }
    }
    return false
}

// removeRestrictiveLocations .. removes any locations that may be considered restrictive by
//  country-wide firewalls, government regulations, restrictions, etc.
//  This list can be updated as appropriate. 
func removeRestrictiveLocations(locations []client.Location) []client.Location {
	var restrictiveLocationIds = []int {11, 72, 184, 445, 446, 447, 448}
	// 11  = Hong Kong
	// 72  = Shanghai
	// 184 = Beijing
	// 445 = Chengdu
	// 446 = Guangzhou
	// 447 = Qingdao
	// 448 = Shenzhen

	var trimmedLocationList []client.Location
	for _, item := range locations { // iterate selected locations
		if !intInList(restrictiveLocationIds, item.ID) && !locationListContainsLocationID(trimmedLocationList, item.ID) {
			trimmedLocationList = append(trimmedLocationList, item)
		}
	}
	return trimmedLocationList
}

//////////////////////////////
// Scheduler helpers
//////////////////////////////

// constructSchedulerWeeklyIntervalsList ... constructs a list of dotcommonitor.WeeklyInterval structs based on the list of weekly_intervals in the TF configuration
func constructSchedulerWeeklyIntervalsList(weeklyIntervals []interface{}) []client.WeeklyInterval {
	//log.Printf("[Dotcom-Monitor] Converting weekly_intervals list to dotcommonitor.WeeklyIntervals list")

	wiList := make([]client.WeeklyInterval, len(weeklyIntervals))

	for i, item := range weeklyIntervals {
		var schemaMap = item.(map[string]interface{})

		wiList[i] = client.WeeklyInterval{
			Days:       convertInterfaceListToStringList(schemaMap["days"].([]interface{})),
			FromMinute: schemaMap["from_minute"].(int),
			ToMinute:   schemaMap["to_minute"].(int),
			Enabled:    schemaMap["enabled"].(bool),
		}
	}

	return wiList
}

// constructSchedulerExcludedTimeIntervalsList ... constructs a list of dotcommonitor.DateTimeInterval structs based on the list of excluded_time_intervals in the TF configuration
func constructSchedulerExcludedTimeIntervalsList(excludedTimeIntervals []interface{}) []client.DateTimeInterval {
	//log.Printf("[Dotcom-Monitor] Converting excluded_time_intervals list to dotcommonitor.DateTimeInterval list")

	etList := make([]client.DateTimeInterval, len(excludedTimeIntervals))

	for i, item := range excludedTimeIntervals {
		var schemaMap = item.(map[string]interface{})

		etList[i] = client.DateTimeInterval{
			From: schemaMap["from_unix"].(int),
			To:   schemaMap["to_unix"].(int),
		}
	}

	return etList
}
