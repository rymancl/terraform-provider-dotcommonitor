// Copyright (C) 2019 IHS Markit.
// All Rights Reserved
//
// NOTICE: All information contained herein is, and remains
// the property of IHS Markit and its suppliers,
// if any. The intellectual and technical concepts contained
// herein are proprietary to IHS Markit and its suppliers
// and may be covered by U.S. and Foreign Patents, patents in
// process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from IHS Markit.

package dotcommonitor

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor/client"
)

func resourceTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTaskCreate,
		Read:   resourceTaskRead,
		Update: resourceTaskUpdate,
		Delete: resourceTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"request_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "GET",
				ValidateFunc: validation.StringInSlice([]string{"GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS", "TRACE", "PATCH"}, true),
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"device_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true, // API disallows moving a task to a different device - Example error: "Task 407648 does not belong to site 195811"
				ValidateFunc: validation.IntAtLeast(0),
			},
			"keyword1": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"keyword2": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"keyword3": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"userpass": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_page_download": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_html": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_frames": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_style_sheets": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_scripts": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_images": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_objects": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_applets": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"download_additional": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ssl_check_certificate_authority": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ssl_check_certificate_cn": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ssl_check_certificate_date": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ssl_check_certificate_revocation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ssl_check_certificate_usage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ssl_expiration_reminder_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"ssl_client_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"get_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
				ConfigMode:    schema.SchemaConfigModeAttr,
				ConflictsWith: []string{"post_params"},
			},
			"post_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
				ConfigMode:    schema.SchemaConfigModeAttr,
				ConflictsWith: []string{"get_params"},
			},
			"header_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
				ConfigMode: schema.SchemaConfigModeAttr,
			},
			"prepare_script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_resolve_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Device Cached", "Non Cached", "TTL Cached", "External DNS Server"}, true),
				// External DNS Server requires dns_server_ip
			},
			"dns_server_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"custom_dns_hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"host": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
				ConfigMode: schema.SchemaConfigModeAttr,
			},
			"task_type_id": { // https://wiki.dotcom-monitor.com/knowledge-base/serverview/
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2, // HTTPS
				ValidateFunc: validation.IntBetween(1, 20),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      120, // API default is 120 seconds
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceTaskCreate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	api := meta.(*client.APIClient)

	task := &client.Task{
		RequestType:                   d.Get("request_type").(string),
		URL:                           d.Get("url").(string),
		Name:                          d.Get("name").(string),
		DeviceID:                      d.Get("device_id").(int),
		Keyword1:                      d.Get("keyword1").(string),
		Keyword2:                      d.Get("keyword2").(string),
		Keyword3:                      d.Get("keyword3").(string),
		UserName:                      d.Get("username").(string),
		UserPass:                      d.Get("userpass").(string),
		SSLCheckCertificateAuthority:  d.Get("ssl_check_certificate_authority").(bool),
		SSLCheckCertificateCN:         d.Get("ssl_check_certificate_cn").(bool),
		SSLCheckCertificateDate:       d.Get("ssl_check_certificate_date").(bool),
		SSLCheckCertificateRevocation: d.Get("ssl_check_certificate_revocation").(bool),
		SSLCheckCertificateUsage:      d.Get("ssl_check_certificate_usage").(bool),
		SSLExpirationReminderInDays:   d.Get("ssl_expiration_reminder_in_days").(int),
		SSLClientCertificate:          d.Get("ssl_client_certificate").(string),
		FullPageDownload:              d.Get("full_page_download").(bool),
		DownloadHTML:                  d.Get("download_html").(bool),
		DownloadFrames:                d.Get("download_frames").(bool),
		DownloadStyleSheets:           d.Get("download_style_sheets").(bool),
		DownloadScripts:               d.Get("download_scripts").(bool),
		DownloadImages:                d.Get("download_images").(bool),
		DownloadObjects:               d.Get("download_objects").(bool),
		DownloadApplets:               d.Get("download_applets").(bool),
		DownloadAdditional:            d.Get("download_additional").(bool),
		GetParams:                     expandInterfaceListToTaskParamList(d.Get("get_params").([]interface{})),
		PostParams:                    expandInterfaceListToTaskParamList(d.Get("post_params").([]interface{})),
		HeaderParams:                  expandInterfaceListToTaskParamList(d.Get("header_params").([]interface{})),
		PrepareScript:                 d.Get("prepare_script").(string),
		DNSResolveMode:                d.Get("dns_resolve_mode").(string),
		DNSserverIP:                   d.Get("dns_server_ip").(string),
		CustomDNSHosts:                constructCustomDNSHostsString(d.Get("custom_dns_hosts").([]interface{})),
		TaskTypeID:                    d.Get("task_type_id").(int),
		Timeout:                       d.Get("timeout").(int),
	}
	log.Printf("[Dotcom-Monitor] task create configuration: %v", task)

	// HACK: Dotcom-Monitor states timeout is in seconds, but it is actually stored in milliseconds
	//  We store it in seconds in state, so here we convert state data into milliseconds
	if task.Timeout != 0 {
		task.Timeout = task.Timeout * 1000
	}

	// create the task
	err := api.CreateTask(task)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to create task: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Task successfully created - ID: %v", fmt.Sprint(task.ID))

	// Set ID
	strID := fmt.Sprint(task.ID)
	d.SetId(strID)

	mutex.Unlock()
	return resourceTaskRead(d, meta)
}

func resourceTaskRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull task ID from state
	taskID, _ := strconv.Atoi(d.Id())

	// Check if task exists before trying to read it
	if !doesTaskExist(taskID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Task does not exist, removing ID %v from state", taskID)
		d.SetId("")
		return nil
	}

	task := &client.Task{}

	api := meta.(*client.APIClient)
	err := api.GetTask(task)

	if task == nil {
		return fmt.Errorf("[Dotcom-Monitor] Task %v does not exist - removing ID from state", task.ID)
	}

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get task: %s", err)
	}

	// HACK: Dotcom-Monitor states timeout is in seconds, but it is actually stored in milliseconds
	//  We store it in seconds in state, so here we convert API data back to seconds
	if task.Timeout != 0 {
		task.Timeout = task.Timeout / 1000
	}

	return nil
}

func resourceTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()

	// Pull task ID from state
	taskID, _ := strconv.Atoi(d.Id())

	// Check if task exists before trying to update it
	if !doesTaskExist(taskID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Task does not exist, removing ID %v from state", taskID)
		d.SetId("")
		return nil
	}

	task := &client.Task{
		ID:                            taskID,
		RequestType:                   d.Get("request_type").(string),
		URL:                           d.Get("url").(string),
		Name:                          d.Get("name").(string),
		DeviceID:                      d.Get("device_id").(int),
		Keyword1:                      d.Get("keyword1").(string),
		Keyword2:                      d.Get("keyword2").(string),
		Keyword3:                      d.Get("keyword3").(string),
		UserName:                      d.Get("username").(string),
		UserPass:                      d.Get("userpass").(string),
		SSLCheckCertificateAuthority:  d.Get("ssl_check_certificate_authority").(bool),
		SSLCheckCertificateCN:         d.Get("ssl_check_certificate_cn").(bool),
		SSLCheckCertificateDate:       d.Get("ssl_check_certificate_date").(bool),
		SSLCheckCertificateRevocation: d.Get("ssl_check_certificate_revocation").(bool),
		SSLCheckCertificateUsage:      d.Get("ssl_check_certificate_usage").(bool),
		SSLExpirationReminderInDays:   d.Get("ssl_expiration_reminder_in_days").(int),
		SSLClientCertificate:          d.Get("ssl_client_certificate").(string),
		FullPageDownload:              d.Get("full_page_download").(bool),
		DownloadHTML:                  d.Get("download_html").(bool),
		DownloadFrames:                d.Get("download_frames").(bool),
		DownloadStyleSheets:           d.Get("download_style_sheets").(bool),
		DownloadScripts:               d.Get("download_scripts").(bool),
		DownloadImages:                d.Get("download_images").(bool),
		DownloadObjects:               d.Get("download_objects").(bool),
		DownloadApplets:               d.Get("download_applets").(bool),
		DownloadAdditional:            d.Get("download_additional").(bool),
		GetParams:                     expandInterfaceListToTaskParamList(d.Get("get_params").([]interface{})),
		PostParams:                    expandInterfaceListToTaskParamList(d.Get("post_params").([]interface{})),
		HeaderParams:                  expandInterfaceListToTaskParamList(d.Get("header_params").([]interface{})),
		PrepareScript:                 d.Get("prepare_script").(string),
		DNSResolveMode:                d.Get("dns_resolve_mode").(string),
		DNSserverIP:                   d.Get("dns_server_ip").(string),
		CustomDNSHosts:                constructCustomDNSHostsString(d.Get("custom_dns_hosts").([]interface{})),
		TaskTypeID:                    d.Get("task_type_id").(int),
		Timeout:                       d.Get("timeout").(int),
	}
	log.Printf("[Dotcom-Monitor] task update configuration: %v", task)

	// HACK: Dotcom-Monitor states timeout is in seconds, but it is actually stored in milliseconds
	//  We store it in seconds in state, so here we convert state data into milliseconds
	if task.Timeout != 0 {
		task.Timeout = task.Timeout * 1000
	}

	api := meta.(*client.APIClient)
	err := api.UpdateTask(task)

	if err != nil {
		mutex.Unlock()
		return fmt.Errorf("[Dotcom-Monitor] Failed to update task: %s", err)
	}

	log.Printf("[Dotcom-Monitor] Task ID: %v successfully updated", fmt.Sprint(task.ID))

	mutex.Unlock()
	return resourceTaskRead(d, meta)
}

func resourceTaskDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull task ID from state
	taskID, _ := strconv.Atoi(d.Id())

	// Check if task exists before trying to remove it
	if !doesTaskExist(taskID, meta) {
		log.Printf("[Dotcom-Monitor] [WARNING] Task does not exist, removing ID %v from state", taskID)
		d.SetId("")
		return nil
	}

	task := &client.Task{
		ID: taskID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteTask(task)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete task: %s", err)
	}

	return nil
}

// doesTaskExist ... determintes if a task with the given taskID exists
func doesTaskExist(taskID int, meta interface{}) bool {
	log.Printf("[Dotcom-Monitor] [DEBUG] Checking if task exists with ID: %v", taskID)
	task := &client.Task{
		ID: taskID,
	}

	// Since an empty HTTP response is a valid 200 from the API, we will determine if
	//  the task exists by comparing the hash of the struct before and after the HTTP call.
	//  If the has does not change, it means nothing else was added, therefore it does not exist.
	//  If the hash changes, the API found the task and added the rest of the fields.
	h := sha256.New()
	t := fmt.Sprintf("%v", task)
	sum := h.Sum([]byte(t))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash before: %x", sum)

	// Try to get task from API
	api := meta.(*client.APIClient)
	err := api.GetTask(task)

	t2 := fmt.Sprintf("%v", task)
	sum2 := h.Sum([]byte(t2))
	log.Printf("[Dotcom-Monitor] [DEBUG] Hash after: %x", sum2)

	// Compare the hashes, and if there was an error from the API we will assume the task exists
	//  to be safe that we do not improperly remove an existing task from state
	if bytes.Equal(sum, sum2) && err == nil {
		log.Println("[Dotcom-Monitor] [DEBUG] No new fields added to the task, therefore the task did not exist")
		return false
	}

	// If we get here, we can assume the task does exist
	return true
}
