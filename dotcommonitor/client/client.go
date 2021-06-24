package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// DotcomMonitorAPIBaseURL ... Base API URL
	// https://wiki.dotcom-monitor.com/knowledge-base/api/
	DotcomMonitorAPIBaseURL = "https://api.dotcom-monitor.com/config_api_v1"

	// AuthCookieName ... Authentication cookie name
	AuthCookieName = ".ASPXFORMSAUTH"
)

// Client ... A client for use with Dotcom-Monitor's REST API.
type Client struct {
	UID        string
	LoggedIn   bool
	AuthCookie string
	Transport  http.RoundTripper
	verbose    bool
}

// NewClient ... Creates a new Httpclient.
func NewClient() *Client {
	return &Client{
		Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
	}
}

// Verbose ... Enable, or disable verbose output from the client.
//
// This will enable (or disable) logging messages that explain what the client
// is about to do, like the endpoint it is about to make a request to. If the
// request fails with an unexpected HTTP response code, then the response body
// will be logged out, as well.
func (c *Client) Verbose(p bool) {
	c.verbose = p
}

// Login ... Establishes a new session with the Dotcom-Monitor API.
func (c *Client) Login(uid string) error {
	var req = LoginBlock{
		UID: uid,
	}

	var resp LoginResponse

	err := c.Do("POST", "login", req, &resp)
	if err != nil {
		return err
	}

	c.LoggedIn = resp.ResponseBlock.Success
	return nil
}

// Logout ... clears cookie
func (c *Client) Logout() {
	c.LoggedIn = false
	c.AuthCookie = ""
}

// IsLoggedIn ... Determines if user is logged in
func (c *Client) IsLoggedIn() bool {
	return c.LoggedIn
}

// newRequest creates a new *http.Request, and sets the following headers:
// <ul>
// <li>Content-Type</li>
// <li>Set-Cookie</li>
// </ul>
func (c *Client) newRequest(method, urlStr string, data []byte) (*http.Request, error) {
	var r *http.Request
	var err error

	if data != nil {
		r, err = http.NewRequest(method, urlStr, bytes.NewReader(data))
	} else {
		r, err = http.NewRequest(method, urlStr, nil)
	}

	r.AddCookie(&http.Cookie{Name: AuthCookieName, Value: c.AuthCookie})
	r.Header.Set("Content-Type", "application/json")

	return r, err
}

// Do ... master function for performing all HTTP calls
func (c *Client) Do(method, endpoint string, requestData, responseData interface{}) error {
	// Throw an error if the user tries to make a request if the client is
	// logged out/unauthenticated, but make an exemption for when the
	// caller is trying to log in.
	if !c.IsLoggedIn() && method != "POST" && endpoint != "login" {
		return errors.New("Will not perform request; client is closed")
	}

	var err error

	// Marshal the request data into a byte slice.
	if c.verbose {
		log.Println("[Dotcom-Monitor] marshaling request data")
	}
	var js []byte
	if requestData != nil {
		js, err = json.Marshal(requestData)
	} else {
		js = []byte("")
	}
	if err != nil {
		return err
	}

	urlStr := fmt.Sprintf("%s/%s", DotcomMonitorAPIBaseURL, endpoint)

	// Create a new http.Request.
	req, err := c.newRequest(method, urlStr, js)
	if err != nil {
		return err
	}

	if c.verbose {
		log.Printf("Making %s request to %q", method, urlStr)
	}

	var resp *http.Response
	resp, err = c.Transport.RoundTrip(req)

	if err != nil {
		if c.verbose {
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			log.Printf("%s", string(respBody))
		}
		return err
	}

	// Get cookies (session token)
	for _, cookie := range resp.Cookies() {
		if cookie.Name == AuthCookieName {
			c.AuthCookie = cookie.Value
			break
		}
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		if resp.ContentLength == 0 {
			// Zero-length content body?
			log.Println("[Dotcom-Monitor] [WARNING] zero-length response body; skipping decoding of response")
			return nil
		}

		//dec := json.NewDecoder(resp.Body)
		text, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Could not read response body")
		}
		if err := json.Unmarshal(text, &responseData); err != nil {
			return fmt.Errorf("error unmarshalling response: %v", err)
		}

		return nil

	case 401:
		// https://wiki.dotcom-monitor.com/knowledge-base/authentication/
		log.Println("[Dotcom-Monitor]: 401 - Unauthorized")
		//c.Login(c.UID)
	}

	// If we got here, this means that the client does not know how to
	// interpret the response, and it should just error out.
	c.Logout()
	reason, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read in response body")
	}
	return fmt.Errorf("server responded with %v: %v",
		resp.Status,
		string(reason))
}
