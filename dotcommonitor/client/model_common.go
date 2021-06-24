package client

//////////////////////////////
// Common structs
//////////////////////////////

// ResponseBlock ... struct for generic API response
type ResponseBlock struct {
	Success          bool     `json:"Success"`
	Result           string   `json:"Result,omitempty"`
	ErrorDescription []string `json:"ErrorDescription,omitempty"`
}

// CreateResponseBlock ... struct for generic API response from creation endpoints
type CreateResponseBlock struct {
	Success          bool     `json:"Success"`
	Result           int      `json:"Result,omitempty"`
	ErrorDescription []string `json:"ErrorDescription,omitempty"`
}
