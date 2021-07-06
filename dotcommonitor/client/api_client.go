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

	var resp CreateTaskResponseBlock

	if err := c.Do("PUT", apiPath, task, &resp); err != nil {
		return fmt.Errorf("Failed to create task: %s", err)
	}

	task.ID = resp.CreateResponseBlock.Result

	return nil
}

// GetTask ... gets the task by ID & returns a ref to the task and any error
// https://wiki.dotcom-monitor.com/knowledge-base/get-task-info/
func (c *APIClient) GetTask(task *Task) error {
	apiPath := fmt.Sprintf("task/%s", fmt.Sprint(task.ID))
	task.ID = 0 // reset ID for provider checks

	if err := c.Do("GET", apiPath, nil, &task); err != nil {
		return fmt.Errorf("Failed to get task: %s", err)
	}

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

	var resp UpdateTaskResponseBlock

	if err := c.Do("POST", apiPath, task, &resp); err != nil {
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

	var resp CreateDeviceResponseBlock

	if err := c.Do("PUT", apiPath, device, &resp); err != nil {
		return fmt.Errorf("Failed to create device: %s", err)
	}

	device.ID = resp.CreateResponseBlock.Result

	return nil
}

// GetDevice ... gets the device by ID & returns a ref to the device and any error
// https://wiki.dotcom-monitor.com/knowledge-base/get-device-info/
func (c *APIClient) GetDevice(device *Device) error {
	apiPath := fmt.Sprintf("device/%s", fmt.Sprint(device.ID))
	device.ID = 0 // reset ID for provider checks

	if err := c.Do("GET", apiPath, nil, &device); err != nil {
		return fmt.Errorf("Failed to get device with ID %v: %s", device.ID, err)
	}

	return nil
}

// UpdateDevice ... updates the device by ID & returns a ref to the device and any error
// https://wiki.dotcom-monitor.com/knowledge-base/edit-device/
func (c *APIClient) UpdateDevice(device *Device) error {
	apiPath := fmt.Sprintf("device/%s", fmt.Sprint(device.ID))

	var resp UpdateDeviceResponseBlock

	if err := c.Do("POST", apiPath, device, &resp); err != nil {
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

	var platformDevicesResp []int

	if err := c.Do("GET", deviceIdsAPIPath, nil, &platformDevicesResp); err != nil {
		return fmt.Errorf("Failed to get device ID's by platform ID: %s", err)
	}

	for _, item := range platformDevicesResp {
		device := &Device{}
		device.ID = item
		if deviceErr := c.GetDevice(device); deviceErr != nil {
			return fmt.Errorf("GetDevicesByName failed: %v", deviceErr)
		}

		// check if the resulting device is the one we're looking for by name
		//   if it is, add it to our results list
		if device.Name == name {
			*devices = append(*devices, *device)
		}
	}

	return nil
}

//////////////////////////////
// Group functions
//////////////////////////////

// CreateGroup ... creates a new group & returns a ref to the group and any error
// https://wiki.dotcom-monitor.com/knowledge-base/create-new-notification-group/
func (c *APIClient) CreateGroup(group *Group) error {
	apiPath := "groups"

	var resp CreateGroupResponseBlock

	if err := c.Do("PUT", apiPath, group, &resp); err != nil {
		return fmt.Errorf("Failed to create group: %s", err)
	}

	group.ID = resp.CreateResponseBlock.Result

	return nil
}

// GetGroup ... gets the group by ID & returns a ref to the group and any error
// https://wiki.dotcom-monitor.com/knowledge-base/get-notification-group-info/
func (c *APIClient) GetGroup(group *Group) error {
	apiPath := fmt.Sprintf("group/%s", fmt.Sprint(group.ID))
	group.ID = 0 // reset ID for provider checks

	if err := c.Do("GET", apiPath, nil, &group); err != nil {
		return fmt.Errorf("Failed to get group: %s", err)
	}

	return nil
}

// UpdateGroup ... updates the group by ID & returns a ref to the group and any error
// https://wiki.dotcom-monitor.com/knowledge-base/edit-alert-group/
func (c *APIClient) UpdateGroup(group *Group) error {
	apiPath := fmt.Sprintf("group/%s", fmt.Sprint(group.ID))

	var resp UpdateGroupResponseBlock

	if err := c.Do("POST", apiPath, group, &resp); err != nil {
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

//////////////////////////////
// Location functions
//////////////////////////////

// GetLocations ... gets the list of all locations available in the account by platform ID
func (c *APIClient) GetLocations(platformID int, includeUnavailable bool, locations *[]Location) error {
	apiPath := fmt.Sprintf("locations/%s", fmt.Sprint(platformID))

	var resp []Location

	if err := c.Do("GET", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("GetLocations failed: %s", err)
	}

	for _, item := range resp {
		if !item.IsDeleted { // don't include deleted locations
			if (item.Available) || (!item.Available && includeUnavailable) {
				*locations = append(*locations, item)
			}
		}
	}

	return nil
}

// GetLocation ... gets the location by id and platform ID
func (c *APIClient) GetLocation(platformID int, locationID int, location *Location) error {
	var locationsList []Location

	if locationsErr := c.GetLocations(platformID, true, &locationsList); locationsErr != nil {
		return fmt.Errorf("GetLocation failed: %v", locationsErr)
	}

	for _, item := range locationsList {
		if item.ID == locationID {
			*location = item
			return nil
		}
	}

	return nil
}

// GetLocationsByName ... gets the locations by name and platform ID
func (c *APIClient) GetLocationsByName(platformID int, name string, includeUnavailable bool, locations *[]Location) error {
	var locationsList []Location

	if locationsErr := c.GetLocations(platformID, includeUnavailable, &locationsList); locationsErr != nil {
		return fmt.Errorf("GetLocationsByName failed: %v", locationsErr)
	}

	for _, item := range locationsList {
		if item.Name == name {
			*locations = append(*locations, item)
		}
	}

	return nil
}

// GetPublicLocations ... gets the public locations in the account for the platform ID
func (c *APIClient) GetPublicLocations(platformID int, includeUnavailable bool, locations *[]Location) error {
	var locationsList []Location

	if locationsErr := c.GetLocations(platformID, includeUnavailable, &locationsList); locationsErr != nil {
		return fmt.Errorf("GetPublicLocations failed: %v", locationsErr)
	}

	for _, item := range locationsList {
		if !item.IsPrivate {
			*locations = append(*locations, item)
		}
	}

	return nil
}

// GetPrivateLocations ... gets the private locations in the account for the platform ID
func (c *APIClient) GetPrivateLocations(platformID int, includeUnavailable bool, locations *[]Location) error {
	var locationsList []Location

	if locationsErr := c.GetLocations(platformID, includeUnavailable, &locationsList); locationsErr != nil {
		return fmt.Errorf("GetPrivateLocations failed: %v", locationsErr)
	}

	// API returned no locations
	if len(locationsList) < 1 {
		return fmt.Errorf("GetPrivateLocations: No private locations found on the platform %v.", platformID)
	}

	for _, item := range locationsList {
		if item.IsPrivate {
			*locations = append(*locations, item)
		}
	}

	return nil
}

//////////////////////////////
// Scheduler functions
//////////////////////////////

// CreateScheduler ... creates a new scheduler & returns a ref to the scheduler and any error
// https://www.dotcom-monitor.com/wiki/knowledge-base/create-new-scheduler/
func (c *APIClient) CreateScheduler(scheduler *Scheduler) error {
	apiPath := "schedulers"

	var resp CreateSchedulerResponseBlock

	if err := c.Do("PUT", apiPath, scheduler, &resp); err != nil {
		return fmt.Errorf("Failed to create scheduler: %s", err)
	}

	scheduler.ID = resp.CreateResponseBlock.Result

	return nil
}

// GetScheduler ... gets the scheduler by ID & returns a ref to the scheduler and any error
// https://www.dotcom-monitor.com/wiki/knowledge-base/get-specific-scheduler-info/
func (c *APIClient) GetScheduler(scheduler *Scheduler) error {
	apiPath := fmt.Sprintf("scheduler/%s", fmt.Sprint(scheduler.ID))
	scheduler.ID = 0 // reset ID for provider checks

	if err := c.Do("GET", apiPath, nil, &scheduler); err != nil {
		return fmt.Errorf("Failed to get scheduler: %s", err)
	}

	return nil
}

// GetSchedulers ... gets all scheduler IDs & returns a ref to the schedulers and any error
// https://www.dotcom-monitor.com/wiki/knowledge-base/get-list-of-schedulers/
func (c *APIClient) GetSchedulers(schedulerIds *[]int) error {
	apiPath := "schedulers"

	if err := c.Do("GET", apiPath, nil, &schedulerIds); err != nil {
		return fmt.Errorf("Failed to get schedulers: %s", err)
	}

	return nil
}

// UpdateScheduler ... updates the scheduler & returns a ref to the scheduler and any error
// https://www.dotcom-monitor.com/wiki/knowledge-base/edit-scheduler/
func (c *APIClient) UpdateScheduler(scheduler *Scheduler) error {
	apiPath := fmt.Sprintf("scheduler/%s", fmt.Sprint(scheduler.ID))

	var resp UpdateSchedulerResponseBlock

	if err := c.Do("POST", apiPath, scheduler, &resp); err != nil {
		return fmt.Errorf("Failed to update scheduler: %s", err)
	}

	return nil
}

// DeleteScheduler ... deletes the scheduler & returns a ref to the scheduler and any error
// https://www.dotcom-monitor.com/wiki/knowledge-base/edit-scheduler/
func (c *APIClient) DeleteScheduler(scheduler *Scheduler) error {
	apiPath := fmt.Sprintf("scheduler/%s", fmt.Sprint(scheduler.ID))

	var resp DeleteSchedulerResponseBlock

	if err := c.Do("DELETE", apiPath, nil, &resp); err != nil {
		return fmt.Errorf("Failed to delete scheduler: %s", err)
	}

	return nil
}

// GetSchedulersByName ... gets the schedulers by name
func (c *APIClient) GetSchedulersByName(name string, schedulers *[]Scheduler) error {
	var allSchedulerIds []int
	var schedulersList []Scheduler

	// first, get all scheduler IDs
	if schedulersErr := c.GetSchedulers(&allSchedulerIds); schedulersErr != nil {
		return fmt.Errorf("GetSchedulersByName failed: %v", schedulersErr)
	}

	for _, item := range allSchedulerIds {
		var scheduler Scheduler
		scheduler.ID = item
		// get full scheduler details for each scheduler ID
		if schedulerErr := c.GetScheduler(&scheduler); schedulerErr != nil {
			return fmt.Errorf("GetSchedulersByName failed: %v", schedulerErr)
		}
		schedulersList = append(schedulersList, scheduler)
	}

	for _, item := range schedulersList {
		if item.Name == name {
			*schedulers = append(*schedulers, item)
		}
	}

	return nil
}
