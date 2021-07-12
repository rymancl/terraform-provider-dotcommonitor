package dotcommonitor

import (
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
				ForceNew:     true, // API gets confused on certain attributes when changing the request type
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
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"full_page_download": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_html": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_frames": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_style_sheets": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_scripts": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_images": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_objects": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_applets": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"download_additional": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_check_certificate_authority": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_check_certificate_cn": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_check_certificate_date": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_check_certificate_revocation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_check_certificate_usage": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_expiration_reminder_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
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
			},
			"prepare_script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_resolve_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Device Cached",
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
		SSLExpirationReminderInDays:   strconv.Itoa(d.Get("ssl_expiration_reminder_in_days").(int)), // HACK: stored as string in API
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
		CustomDNSHosts:                flattenCustomDnsHostsToString(d.Get("custom_dns_hosts").([]interface{})),
		TaskTypeID:                    d.Get("task_type_id").(int),
		Timeout:                       d.Get("timeout").(int),
	}
	log.Printf("[Dotcom-Monitor] task create configuration: %v", task)

	// HACK: Dotcom-Monitor states timeout is in seconds, but it is actually stored in milliseconds
	//  We store it in seconds in state, so here we convert state data into milliseconds
	task.Timeout = task.Timeout * 1000

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

	task := &client.Task{}
	task.ID = taskID

	api := meta.(*client.APIClient)
	err := api.GetTask(task)

	if task == nil {
		return fmt.Errorf("[Dotcom-Monitor] Task %v does not exist - removing ID from state", task.ID)
	}

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to get task: %s", err)
	}

	// Check if task exists before trying to read it
	if !(task.ID > 0) {
		log.Printf("[Dotcom-Monitor] [WARNING] Task does not exist, removing ID %v from state", task.ID)
		d.SetId("")
		return nil
	}

	// HACK: Dotcom-Monitor states timeout is in seconds, but it is actually stored in milliseconds
	//  We store it in seconds in state, so here we convert API data back to seconds
	if task.Timeout != 0 {
		task.Timeout = task.Timeout / 1000
	}

	d.Set("request_type", task.RequestType)
	d.Set("url", task.URL)
	d.Set("name", task.Name)
	d.Set("device_id", task.DeviceID)
	d.Set("keyword1", task.Keyword1)
	d.Set("keyword2", task.Keyword2)
	d.Set("keyword3", task.Keyword3)
	d.Set("username", task.UserName)
	d.Set("userpass", task.UserPass)
	d.Set("full_page_download", task.FullPageDownload)
	d.Set("download_html", task.DownloadHTML)
	d.Set("download_frames", task.DownloadFrames)
	d.Set("download_style_sheets", task.DownloadStyleSheets)
	d.Set("download_scripts", task.DownloadScripts)
	d.Set("download_images", task.DownloadImages)
	d.Set("download_objects", task.DownloadObjects)
	d.Set("download_applets", task.DownloadApplets)
	d.Set("download_additional", task.DownloadAdditional)
	d.Set("ssl_check_certificate_authority", task.SSLCheckCertificateAuthority)
	d.Set("ssl_check_certificate_cn", task.SSLCheckCertificateCN)
	d.Set("ssl_check_certificate_date", task.SSLCheckCertificateDate)
	d.Set("ssl_check_certificate_revocation", task.SSLCheckCertificateRevocation)
	d.Set("ssl_check_certificate_usage", task.SSLCheckCertificateUsage)
	d.Set("ssl_expiration_reminder_in_days", task.SSLExpirationReminderInDays)
	d.Set("ssl_client_certificate", task.SSLClientCertificate)
	if task.GetParams != nil {
		d.Set("get_params", task.GetParams)
	}
	if task.PostParams != nil {
		d.Set("post_params", task.PostParams)
	}
	if task.HeaderParams != nil {
		d.Set("header_params", task.HeaderParams)
	}
	d.Set("prepare_script", task.PrepareScript)
	d.Set("dns_resolve_mode", task.DNSResolveMode)
	d.Set("dns_server_ip", task.DNSserverIP)
	d.Set("custom_dns_hosts", task.CustomDNSHosts) // not really sure why this works, but it does
	d.Set("task_type_id", task.TaskTypeID)
	d.Set("timeout", task.Timeout)

	return nil
}

func resourceTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	d.Partial(true)

	// Pull task ID from state
	taskID, _ := strconv.Atoi(d.Id())

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
		SSLExpirationReminderInDays:   strconv.Itoa(d.Get("ssl_expiration_reminder_in_days").(int)), // HACK: stored as string in API
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
		CustomDNSHosts:                flattenCustomDnsHostsToString(d.Get("custom_dns_hosts").([]interface{})),
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
	d.Partial(false)
	return resourceTaskRead(d, meta)
}

func resourceTaskDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Pull task ID from state
	taskID, _ := strconv.Atoi(d.Id())

	task := &client.Task{
		ID: taskID,
	}

	api := meta.(*client.APIClient)
	err := api.DeleteTask(task)

	if err != nil {
		return fmt.Errorf("[Dotcom-Monitor] Failed to delete task: %s", err)
	}

	d.SetId("")

	return nil
}
