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
	Number         string `json:"Number,omitempty"`         // for Phone, Sms
	Code           string `json:"Code,omitempty"`           // for Phone
	IntegrationKey string `json:"IntegrationKey,omitempty"` // for PagerDuty
	IntegrationURL string `json:"IntegrationURL,omitempty"` // for AlertOps
	WebHook        string `json:"WebHook,omitempty"`        // for Slack, Teams
	Community      string `json:"Community,omitempty"`      // for SNMP
	Host           string `json:"Host,omitempty"`           // for SNMP
	UserID         int    `json:"UserId,omitempty"`         // for SNMP
	Version        string `json:"Version,omitempty"`        // for SNMP
	// Message        string `json:"Message,omitempty"`        // for Script?
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
