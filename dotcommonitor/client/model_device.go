package client

// Device ... simple struct to hold device details
type Device struct {
	ID                      int                       `json:"Id,omitempty"`
	Name                    string                    `json:"Name"`
	PlatformID              int                       `json:"Platform_Id"`
	Frequency               int                       `json:"Frequency"`
	Locations               []int                     `json:"Locations"`
	AvoidSimultaneousChecks bool                      `json:"Avoid_Simultaneous_Checks"`
	AlertSilenceMin         int                       `json:"Alert_Silence_Min"`
	FalsePositiveCheck      bool                      `json:"False_Positive_Check"`
	SendUptimeAlert         bool                      `json:"Send_Uptime_Alert"`
	StatusDescription       string                    `json:"Status_Description"`
	Postpone                bool                      `json:"Postpone"`
	OwnerDeviceID           int                       `json:"Owner_Device_Id"`
	FilterID                int                       `json:"Filter_Id"`
	SchedulerID             int                       `json:"Scheduler_Id"`
	NumberOfTasks           int                       `json:"Number_Of_Tasks"`
	PackageID               int                       `json:"Package_Id"`
	Notifications           *DeviceNotificationsBlock `json:"Notifications"`
}

// DeviceNotificationsBlock ... struct for device notifications
type DeviceNotificationsBlock struct {
	EMailFlag               bool                              `json:"E_Mail_Flag"`
	EMailAddress            string                            `json:"E_Mail_Address,omitempty"`
	EMailTimeIntervalMin    int                               `json:"E_Mail_TimeInterval_Min,omitempty"`
	WLDeviceFlag            bool                              `json:"WL_Device_Flag"`
	WLDeviceEmailAddress    string                            `json:"WL_Device_Email_Address,omitempty"`
	WLDeviceTimeIntervalMin int                               `json:"WL_Device_TimeInterval_Min,omitempty"`
	PagerFlag               bool                              `json:"Pager_Flag"`
	PagerAreaCode           string                            `json:"Pager_Area_Code,omitempty"`
	PagerPhone              string                            `json:"Pager_Phone,omitempty"`
	PagerNumCode            string                            `json:"Pager_Num_Code,omitempty"`
	PagerTimeIntervalMin    int                               `json:"Pager_TimeInterval_Min,omitempty"`
	PhoneFlag               bool                              `json:"Phone_Flag"`
	PhoneAreaCode           string                            `json:"Phone_Area_Code,omitempty"`
	PhonePhone              string                            `json:"Phone_Phone,omitempty"`
	PhoneTimeIntervalMin    int                               `json:"Phone_TimeInterval_Min,omitempty"`
	SMSFlag                 bool                              `json:"SMS_Flag"`
	SMSPhone                string                            `json:"SMS_Phone,omitempty"`
	SMSTimeIntervalMin      int                               `json:"SMS_TimeInterval_Min,omitempty"`
	ScriptFlag              bool                              `json:"Script_Flag"`
	ScriptBatchFileName     string                            `json:"Script_Batch_File_Name,omitempty"`
	ScriptTimeIntervalMin   int                               `json:"Script_TimeInterval_Min,omitempty"`
	SNMPTimeIntervalMin     int                               `json:"SNMP_TimeInterval_Min,omitempty"`
	NotificationGroups      []NotificationsNotificationGroups `json:"Notification_Groups"`
}

// NotificationsNotificationGroups ... struct for device notifications notification groups
type NotificationsNotificationGroups struct {
	ID           int `json:"Id"`
	TimeShiftMin int `json:"Time_Shift_Min"`
}

// CreateDeviceResponseBlock ... struct for create device response
type CreateDeviceResponseBlock struct {
	CreateResponseBlock
}

// UpdateDeviceResponseBlock ... struct for update device response
type UpdateDeviceResponseBlock struct {
	ResponseBlock
}

// DeleteDeviceResponseBlock ... struct for delete device response
type DeleteDeviceResponseBlock struct {
	ResponseBlock
}

// PlatformDevices .. struct for devices on a platform
type PlatformDevices struct {
	DeviceIDList []int
}
