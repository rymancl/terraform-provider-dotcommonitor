package client

// Filter ... simple struct to hold filter details
type Filter struct {
	ID          int    `json:"Id,omitempty"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	Rules       Rule   `json:"Rules,omitempty"`
	Items       []Item `json:"Items,omitempty"`
	AssignedTo  []int  `json:"Assigned_To,omitempty"`
}

// Rule ... struct for filter rule
type Rule struct {
	TimeMin           int  `json:"Time_Min"`
	NumberOfLocations int  `json:"Number_Of_Locations"`
	NumberOfTasks     int  `json:"Number_Of_Tasks"`
	OwnerDevice       bool `json:"Owner_Device"`
}

// Item ... struct for filter errors to ignore
type Item struct {
	ErrorType         string `json:"Error_Type"`
	ErrorCodeToIgnore []int  `json:"Error_Code_To_Ignore"`
}

// CreateFilterResponseBlock ... struct for create filter response
type CreateFilterResponseBlock struct {
	CreateResponseBlock
}

// UpdateFilterResponseBlock ... struct for update filter response
type UpdateFilterResponseBlock struct {
	ResponseBlock
}

// DeleteFilterResponseBlock ... struct for delete filter response
type DeleteFilterResponseBlock struct {
	ResponseBlock
}
