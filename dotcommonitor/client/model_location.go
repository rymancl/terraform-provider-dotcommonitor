package client

// Location ... struct for Location
type Location struct {
	ID          int         `json:"Id,omitempty"`
	Name        string      `json:"Name"`
	Available	bool        `json:"Available"`
	IsDeleted	bool        `json:"IsDeleted"`
	IsPrivate	bool        `json:"IsPrivate"`
}
