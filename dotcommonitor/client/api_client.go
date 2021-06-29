package client

import (
	"fmt"
	"net/http"
)

// APIClient A client with extra helper methods for common actions
type APIClient struct {
	Client
}

// NewAPIClient Creates a new APIClient
func NewAPIClient() *APIClient {
	return &APIClient{
		Client{
			Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
			//verbose:   true,
		}}
}

//////////////////////////////
// Task functions
//////////////////////////////

// CreateTask ... creates a new task & returns a ref to the task and any error
// https://wiki.dotcom-monitor.com/knowledge-base/create-new-task/
func (c *APIClient) CreateTask(task *Task) error {
	apiPath := "tasks"

	data := &Task{
		RequestType:                   task.RequestType,
		URL:                           task.URL,
		DeviceID:                      task.DeviceID,
		TaskTypeID:                    task.TaskTypeID,
		Name:                          task.Name,
		Keyword1:                      task.Keyword1,
		Keyword2:                      task.Keyword2,
		Keyword3:                      task.Keyword3,
		UserName:                      task.UserName,
		UserPass:                      task.UserPass,
		SSLCheckCertificateAuthority:  task.SSLCheckCertificateAuthority,
		SSLCheckCertificateCN:         task.SSLCheckCertificateCN,
		SSLCheckCertificateDate:       task.SSLCheckCertificateDate,
		SSLCheckCertificateRevocation: task.SSLCheckCertificateRevocation,
		SSLCheckCertificateUsage:      task.SSLCheckCertificateUsage,
		SSLExpirationReminderInDays:   task.SSLExpirationReminderInDays,
		SSLClientCertificate:          task.SSLClientCertificate,
		FullPageDownload:              task.FullPageDownload,
		DownloadHTML:                  task.DownloadHTML,
		DownloadFrames:                task.DownloadFrames,
		DownloadStyleSheets:           task.DownloadStyleSheets,
		DownloadScripts:               task.DownloadScripts,
		DownloadImages:                task.DownloadImages,
		DownloadObjects:               task.DownloadObjects,
		DownloadApplets:               task.DownloadApplets,
		DownloadAdditional:            task.DownloadAdditional,
		GetParams:                     task.GetParams,
		PostParams:                    task.PostParams,
		HeaderParams:                  task.HeaderParams,
		Timeout:                       task.Timeout,
		PrepareScript:                 task.PrepareScript,
		DNSResolveMode:                task.DNSResolveMode,
		DNSserverIP:                   task.DNSserverIP,
		CustomDNSHosts:                task.CustomDNSHosts,
	}

	var resp CreateTaskResponseBlock

	if err := c.Do("PUT", apiPath, data, &resp); err != nil {
		return fmt.Errorf("Failed to create task: %s", err)
	}

	task.ID = resp.CreateResponseBlock.Result

	return nil
}

// GetTask ... gets the task by ID & returns a ref to the task and any error
// https://wiki.dotcom-monitor.com/knowledge-base/get-task-info/
func (c *APIClient) GetTask(task *Task) error {
	apiPath := fmt.Sprintf("task/%s", fmt.Sprint(task.ID))

	var resp Task

	if err := c.Do("GET", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to get task: %s", err)
	}

	task.RequestType = resp.RequestType
	task.URL = resp.URL
	task.Keyword1 = resp.Keyword1
	task.Keyword2 = resp.Keyword2
	task.Keyword3 = resp.Keyword3
	task.UserName = resp.UserName
	task.UserPass = resp.UserPass
	task.SSLCheckCertificateAuthority = resp.SSLCheckCertificateAuthority
	task.SSLCheckCertificateCN = resp.SSLCheckCertificateCN
	task.SSLCheckCertificateDate = resp.SSLCheckCertificateDate
	task.SSLCheckCertificateRevocation = resp.SSLCheckCertificateRevocation
	task.SSLCheckCertificateUsage = resp.SSLCheckCertificateUsage
	task.SSLExpirationReminderInDays = resp.SSLExpirationReminderInDays
	task.SSLClientCertificate = resp.SSLClientCertificate
	task.FullPageDownload = resp.FullPageDownload
	task.DownloadHTML = resp.DownloadHTML
	task.DownloadFrames = resp.DownloadFrames
	task.DownloadStyleSheets = resp.DownloadStyleSheets
	task.DownloadScripts = resp.DownloadScripts
	task.DownloadImages = resp.DownloadImages
	task.DownloadObjects = resp.DownloadObjects
	task.DownloadApplets = resp.DownloadApplets
	task.DownloadAdditional = resp.DownloadAdditional
	task.GetParams = resp.GetParams
	task.PostParams = resp.PostParams
	task.HeaderParams = resp.HeaderParams
	task.PrepareScript = resp.PrepareScript
	task.DNSResolveMode = resp.DNSResolveMode
	task.DNSserverIP = resp.DNSserverIP
	task.CustomDNSHosts = resp.CustomDNSHosts
	task.DeviceID = resp.DeviceID
	task.TaskTypeID = resp.TaskTypeID
	task.Name = resp.Name
	task.Timeout = resp.Timeout

	return nil
}

// GetTaskListByDevice ... gets a list of tasks for the device & returns a ref to the tasks and any error
// https://www.dotcom-monitor.com/wiki/knowledge-base/get-task-list-by-device/
func (c *APIClient) GetTaskListByDevice(device *Device, tasks *[]Task) error {
	apiPath := fmt.Sprintf("device/%s/tasks", fmt.Sprint(device.ID))

	var resp []int

	if err := c.Do("GET", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to get task list by device: %s", err)
	}

	for _, item := range resp {
		task := &Task{}
		task.ID = item
		if taskErr := c.GetTask(task); taskErr != nil {
			return fmt.Errorf("GetTaskListByDevice failed: %v", taskErr)
		}
		*tasks = append(*tasks, *task)
	}

	return nil
}

// GetDeviceTasksByName ... gets a task list for the device by name & returns a ref to the tasks and any error
func (c *APIClient) GetDeviceTasksByName(deviceID int, name string, tasks *[]Task) error {
	device := &Device{}
	device.ID = deviceID
	var taskList []Task

	if deviceErr := c.GetTaskListByDevice(device, &taskList); deviceErr != nil {
		return fmt.Errorf("GetDeviceTasksByName failed: %v", deviceErr)
	}

	for _, item := range taskList {
		if item.Name == name {
			*tasks = append(*tasks, item)
		}
	}

	return nil
}

// UpdateTask ... updates the task by ID & returns a ref to the task and any error
// https://wiki.dotcom-monitor.com/knowledge-base/edit-task/
func (c *APIClient) UpdateTask(task *Task) error {
	apiPath := fmt.Sprintf("task/%s", fmt.Sprint(task.ID))

	data := &Task{
		RequestType:                   task.RequestType,
		URL:                           task.URL,
		DeviceID:                      task.DeviceID,
		TaskTypeID:                    task.TaskTypeID,
		Name:                          task.Name,
		Keyword1:                      task.Keyword1,
		Keyword2:                      task.Keyword2,
		Keyword3:                      task.Keyword3,
		UserName:                      task.UserName,
		UserPass:                      task.UserPass,
		SSLCheckCertificateAuthority:  task.SSLCheckCertificateAuthority,
		SSLCheckCertificateCN:         task.SSLCheckCertificateCN,
		SSLCheckCertificateDate:       task.SSLCheckCertificateDate,
		SSLCheckCertificateRevocation: task.SSLCheckCertificateRevocation,
		SSLCheckCertificateUsage:      task.SSLCheckCertificateUsage,
		SSLExpirationReminderInDays:   task.SSLExpirationReminderInDays,
		SSLClientCertificate:          task.SSLClientCertificate,
		FullPageDownload:              task.FullPageDownload,
		DownloadHTML:                  task.DownloadHTML,
		DownloadFrames:                task.DownloadFrames,
		DownloadStyleSheets:           task.DownloadStyleSheets,
		DownloadScripts:               task.DownloadScripts,
		DownloadImages:                task.DownloadImages,
		DownloadObjects:               task.DownloadObjects,
		DownloadApplets:               task.DownloadApplets,
		DownloadAdditional:            task.DownloadAdditional,
		GetParams:                     task.GetParams,
		PostParams:                    task.PostParams,
		HeaderParams:                  task.HeaderParams,
		Timeout:                       task.Timeout,
		PrepareScript:                 task.PrepareScript,
		DNSResolveMode:                task.DNSResolveMode,
		DNSserverIP:                   task.DNSserverIP,
		CustomDNSHosts:                task.CustomDNSHosts,
	}

	var resp UpdateTaskResponseBlock

	if err := c.Do("POST", apiPath, data, &resp); err != nil {
		return fmt.Errorf("Failed to update task: %s", err)
	}

	return nil
}

// DeleteTask ... deletes the task by ID & returns a ref to the task and any error
// https://wiki.dotcom-monitor.com/knowledge-base/delete-task/
func (c *APIClient) DeleteTask(task *Task) error {
	apiPath := fmt.Sprintf("task/%s", fmt.Sprint(task.ID))

	var resp DeleteTaskResponseBlock

	if err := c.Do("DELETE", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to delete task: %s", err)
	}

	return nil
}

//////////////////////////////
// Device functions
//////////////////////////////

// CreateDevice ... creates a new device & returns a ref to the device and any error
// https://wiki.dotcom-monitor.com/knowledge-base/create-new-device/
func (c *APIClient) CreateDevice(device *Device) error {
	apiPath := "devices"

	data := &Device{
		Name:                    device.Name,
		PlatformID:              device.PlatformID,
		Frequency:               device.Frequency,
		Locations:               device.Locations,
		AvoidSimultaneousChecks: device.AvoidSimultaneousChecks,
		AlertSilenceMin:         device.AlertSilenceMin,
		FalsePositiveCheck:      device.FalsePositiveCheck,
		SendUptimeAlert:         device.SendUptimeAlert,
		Postpone:                device.Postpone,
		OwnerDeviceID:           device.OwnerDeviceID,
		FilterID:                device.FilterID,
		SchedulerID:             device.SchedulerID,
		Notifications:           device.Notifications,
	}

	var resp CreateDeviceResponseBlock

	if err := c.Do("PUT", apiPath, data, &resp); err != nil {
		return fmt.Errorf("Failed to create device: %s", err)
	}

	device.ID = resp.CreateResponseBlock.Result

	return nil
}

// GetDevice ... gets the device by ID & returns a ref to the device and any error
// https://wiki.dotcom-monitor.com/knowledge-base/get-device-info/
func (c *APIClient) GetDevice(device *Device) error {
	apiPath := fmt.Sprintf("device/%s", fmt.Sprint(device.ID))

	var resp Device

	if err := c.Do("GET", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to get device with ID %v: %s", device.ID, err)
	}

	device.Name = resp.Name
	device.PlatformID = resp.PlatformID
	device.Frequency = resp.Frequency
	device.Locations = resp.Locations
	device.AvoidSimultaneousChecks = resp.AvoidSimultaneousChecks
	device.AlertSilenceMin = resp.AlertSilenceMin
	device.FalsePositiveCheck = resp.FalsePositiveCheck
	device.SendUptimeAlert = resp.SendUptimeAlert
	device.StatusDescription = resp.StatusDescription
	device.Postpone = resp.Postpone
	device.OwnerDeviceID = resp.OwnerDeviceID
	device.FilterID = resp.FilterID
	device.SchedulerID = resp.SchedulerID
	device.NumberOfTasks = resp.NumberOfTasks
	device.PackageID = resp.PackageID
	device.Notifications = resp.Notifications

	return nil
}

// UpdateDevice ... updates the device by ID & returns a ref to the device and any error
// https://wiki.dotcom-monitor.com/knowledge-base/edit-device/
func (c *APIClient) UpdateDevice(device *Device) error {
	apiPath := fmt.Sprintf("device/%s", fmt.Sprint(device.ID))

	data := &Device{
		Name:                    device.Name,
		PlatformID:              device.PlatformID,
		Frequency:               device.Frequency,
		Locations:               device.Locations,
		AvoidSimultaneousChecks: device.AvoidSimultaneousChecks,
		AlertSilenceMin:         device.AlertSilenceMin,
		FalsePositiveCheck:      device.FalsePositiveCheck,
		SendUptimeAlert:         device.SendUptimeAlert,
		Postpone:                device.Postpone,
		OwnerDeviceID:           device.OwnerDeviceID,
		FilterID:                device.FilterID,
		SchedulerID:             device.SchedulerID,
		Notifications:           device.Notifications,
	}

	var resp UpdateDeviceResponseBlock

	if err := c.Do("POST", apiPath, data, &resp); err != nil {
		return fmt.Errorf("Failed to update device: %s", err)
	}

	return nil
}

// DeleteDevice ... deletes the device by ID & returns a ref to the device and any error
// https://wiki.dotcom-monitor.com/knowledge-base/delete-device/
func (c *APIClient) DeleteDevice(device *Device) error {
	apiPath := fmt.Sprintf("device/%s", fmt.Sprint(device.ID))

	var resp DeleteDeviceResponseBlock

	if err := c.Do("DELETE", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to delete device: %s", err)
	}

	return nil
}

// GetDevicesByName ... gets a list of devices on the given platform based on the given name
func (c *APIClient) GetDevicesByName(platformID int, name string, devices *[]Device) error {
	deviceIdsAPIPath := fmt.Sprintf("devices/%s", fmt.Sprint(platformID))

	//var resp PlatformDevices
	var platformDevicesResp []int

	if err := c.Do("GET", deviceIdsAPIPath, nil, &platformDevicesResp); err != nil {
		return fmt.Errorf("Failed to get device ID's by platform ID: %s", err)
	}

	//deviceList := make([]*Device, len(platformDevicesResp))
	//client := NewAPIClient()

	for _, item := range platformDevicesResp {
		device := &Device{}
		device.ID = item
		if deviceErr := c.GetDevice(device); deviceErr != nil {
			return fmt.Errorf("GetDevicesByName failed: %v", deviceErr)
		}

		// check if the resulting device is the one we're looking for by name
		//   if it is, add it to our results list
		if device.Name == name {
			//devices[i] = device
			*devices = append(*devices, *device)
		}
		//device = nil
	}

	//return fmt.Errorf("***** devices length: %v", len(devices))

	return nil
}

//////////////////////////////
// Group functions
//////////////////////////////

// CreateGroup ... creates a new group & returns a ref to the group and any error
// https://wiki.dotcom-monitor.com/knowledge-base/create-new-notification-group/
func (c *APIClient) CreateGroup(group *Group) error {
	apiPath := "groups"

	data := &Group{
		Name:        group.Name,
		SchedulerID: group.SchedulerID,
		Addresses:   group.Addresses,
		AssignedTo:  group.AssignedTo,
	}

	var resp CreateGroupResponseBlock

	if err := c.Do("PUT", apiPath, data, &resp); err != nil {
		return fmt.Errorf("Failed to create group: %s", err)
	}

	group.ID = resp.CreateResponseBlock.Result

	return nil
}

// GetGroup ... gets the group by ID & returns a ref to the group and any error
// https://wiki.dotcom-monitor.com/knowledge-base/get-notification-group-info/
func (c *APIClient) GetGroup(group *Group) error {
	apiPath := fmt.Sprintf("group/%s", fmt.Sprint(group.ID))

	var resp Group

	if err := c.Do("GET", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to get group: %s", err)
	}

	group.Name = resp.Name
	group.SchedulerID = resp.SchedulerID
	group.Addresses = resp.Addresses
	group.AssignedTo = resp.AssignedTo

	return nil
}

// UpdateGroup ... updates the group by ID & returns a ref to the group and any error
// https://wiki.dotcom-monitor.com/knowledge-base/edit-alert-group/
func (c *APIClient) UpdateGroup(group *Group) error {
	apiPath := fmt.Sprintf("group/%s", fmt.Sprint(group.ID))

	data := &Group{
		Name:        group.Name,
		SchedulerID: group.SchedulerID,
		Addresses:   group.Addresses,
		AssignedTo:  group.AssignedTo,
	}

	var resp UpdateGroupResponseBlock

	if err := c.Do("POST", apiPath, data, &resp); err != nil {
		return fmt.Errorf("Failed to update group: %s", err)
	}

	return nil
}

// DeleteGroup ... deletes the group by ID & returns a ref to the group and any error
// https://wiki.dotcom-monitor.com/knowledge-base/delete-alert-group/
func (c *APIClient) DeleteGroup(group *Group) error {
	apiPath := fmt.Sprintf("group/%s", fmt.Sprint(group.ID))

	var resp DeleteGroupResponseBlock

	if err := c.Do("DELETE", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to delete group: %s", err)
	}

	return nil
}

// GetGroupsByName ... gets a list of groups based on the given name
func (c *APIClient) GetGroupsByName(name string, groups *[]Group) error {
	groupIdsAPIPath := "groups"

	var groupsResp []int

	if err := c.Do("GET", groupIdsAPIPath, nil, &groupsResp); err != nil {
		return fmt.Errorf("Failed to get group ID's by name: %s", err)
	}

	for _, item := range groupsResp {
		group := &Group{}
		group.ID = item
		if groupErr := c.GetGroup(group); groupErr != nil {
			return fmt.Errorf("GetGroupsByName failed: %v", groupErr)
		}

		// check if the resulting group is the one we're looking for by name
		//   if it is, add it to our results list
		if group.Name == name {
			*groups = append(*groups, *group)
		}
	}

	return nil
}
