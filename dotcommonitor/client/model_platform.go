package client

// Platform ... simple struct to hold platform details
type Platform struct {
	ID        int       `json:"Id,omitempty"`
	Name      string    `json:"Name"`
	Packages  []Package `json:"Packages"`
	Available bool      `json:"Available"`
}

// Package ... struct for platform packages
type Package struct {
	PackageID   int    `json:"Package_Id"`
	PackageName string `json:"Package_Name"`
	PlatformID  int    `json:"Platform_Id"`
}
