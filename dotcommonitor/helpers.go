package dotcommonitor

import (
	"bytes"
	"log"
	"strconv"
	"strings"

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

// constructCustomDNSHostsString ... returns a string required for the syntax of "CustomDNSHosts"
//  Syntax:  <host>=<ip>;
func constructCustomDNSHostsString(hosts []interface{}) string {
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
	log.Printf("[Dotcom-Monitor] [constructCustomDNSHostsString] Result string: %s", resultString)

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

// convertLocationsToIntList ... type asserting a list of string to a list of int
func convertLocationsToIntList(interfaceList []interface{}) []int {
	//log.Printf("[Dotcom-Monitor] convertLocationsToIntList - Converting %v to int[]", interfaceList)
	intList := make([]int, len(interfaceList))

	for i := range interfaceList {
		intList[i] = interfaceList[i].(int)
	}

	return intList
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
