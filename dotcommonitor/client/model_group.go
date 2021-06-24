package client

// Group ... struct for Group
type Group struct {
	ID          int         `json:"Id,omitempty"`
	Name        string      `json:"Name"`
	SchedulerID int         `json:"Scheduler_Id,omitempty"`
	Addresses   []Addresses `json:"Addresses,omitempty"`
	AssignedTo  []int       `json:"Assigned_To,omitempty"`
}

// Addresses ... struct for delivery addresses
type Addresses struct {
	Type           string `json:"Type"`
	TemplateID     int    `json:"Template_Id,omitempty"`
	Address        string `json:"Address,omitempty"`        // for Email
	Number         string `json:"Number,omitempty"`         // for Phone, Pager, Sms
	Code           string `json:"Code,omitempty"`           // for Phone, Pager
	Message        string `json:"Message,omitempty"`        // for Pager
	IntegrationKey string `json:"IntegrationKey,omitempty"` // for PagerDuty
}

// CreateGroupResponseBlock ... struct for create group response
type CreateGroupResponseBlock struct {
	CreateResponseBlock
}

// UpdateGroupResponseBlock ... struct for update group response
type UpdateGroupResponseBlock struct {
	ResponseBlock
}

// DeleteGroupResponseBlock ... struct for delete group response
type DeleteGroupResponseBlock struct {
	ResponseBlock
}
