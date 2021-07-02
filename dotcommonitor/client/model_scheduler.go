package client

// Scheduler ... struct for Scheduler
type Scheduler struct {
	ID                    int                `json:"Id,omitempty"`
	Name                  string             `json:"Name"`
	Description           string             `json:"Description,omitempty"`
	WeeklyIntervals       []WeeklyInterval   `json:"Weekly_Intervals,omitempty"`
	ExcludedTimeIntervals []DateTimeInterval `json:"Date_Time_Intervals,omitempty"` // API calls this Date_Time_Intervals but it represents excluded monitoring times, i.e. maintenance windows
	AssignedTo            AssignedTo         `json:"Assigned_To,omitempty"`
}

// WeeklyInterval ... struct for repeating weekly interval
type WeeklyInterval struct {
	Days       []string `json:"Days,omitempty"`
	FromMinute int      `json:"From_Min,omitempty"`
	ToMinute   int      `json:"To_Min,omitempty"`
	Enabled    bool     `json:"Included,omitempty"`
}

// DateTimeInterval ... struct for non-repeat date time interval
type DateTimeInterval struct {
	From int `json:"From,omitempty"` // UNIX time format
	To   int `json:"To,omitempty"`   // UNIX time format
}

// AssignedTo ... struct for non-repeat date time interval
type AssignedTo struct {
	Devices []int `json:"Devices,omitempty"`
	Groups  []int `json:"Notification_Groups,omitempty"`
}

// CreateSchedulerResponseBlock ... struct for create scheduler response
type CreateSchedulerResponseBlock struct {
	CreateResponseBlock
}

// UpdateSchedulerResponseBlock ... struct for update scheduler response
type UpdateSchedulerResponseBlock struct {
	ResponseBlock
}

// DeleteSchedulerResponseBlock ... struct for delete scheduler response
type DeleteSchedulerResponseBlock struct {
	ResponseBlock
}
