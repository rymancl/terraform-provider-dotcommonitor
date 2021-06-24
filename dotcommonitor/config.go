package dotcommonitor

import (
	"fmt"
	"log"

	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// Config struct for data required to log into API
type Config struct {
	UID string
}

// Client returns a new client.
func (c *Config) Client() (*client.APIClient, error) {
	client := client.NewAPIClient()

	// API Login
	err := client.Login(c.UID)

	if err != nil {
		return nil, fmt.Errorf("Error logging into API: %s", err)
	}

	log.Printf("[INFO] [Dotcom-Monitor] client configured for API key: %s", c.UID)
	//log.Print(resp)

	return client, nil
}
