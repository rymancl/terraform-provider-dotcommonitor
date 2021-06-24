package client

// Task ... simple struct to hold task details
type Task struct {
	ID                            int         `json:"Id,omitempty"`
	RequestType                   string      `json:"RequestType"`
	URL                           string      `json:"Url"`
	Keyword1                      string      `json:"Keyword1,omitempty"`
	Keyword2                      string      `json:"Keyword2,omitempty"`
	Keyword3                      string      `json:"Keyword3,omitempty"`
	UserName                      string      `json:"UserName,omitempty"`
	UserPass                      string      `json:"UserPass,omitempty"`
	SSLCheckCertificateAuthority  bool        `json:"CheckCertificateAuthority"`
	SSLCheckCertificateCN         bool        `json:"CheckCertificateCN"`
	SSLCheckCertificateDate       bool        `json:"CheckCertificateDate"`
	SSLCheckCertificateRevocation bool        `json:"CheckCertificateRevocation"`
	SSLCheckCertificateUsage      bool        `json:"CheckCertificateUsage"`
	SSLExpirationReminderInDays   int         `json:"ExpirationReminderInDays,omitempty"`
	SSLClientCertificate          string      `json:"ClientCertificate,omitempty"`
	FullPageDownload              bool        `json:"FullPageDownload"`
	DownloadHTML                  bool        `json:"Download_Html"`
	DownloadFrames                bool        `json:"Download_Frames"`
	DownloadStyleSheets           bool        `json:"Download_StyleSheets"`
	DownloadScripts               bool        `json:"Download_Scripts"`
	DownloadImages                bool        `json:"Download_Images"`
	DownloadObjects               bool        `json:"Download_Objects"`
	DownloadApplets               bool        `json:"Download_Applets"`
	DownloadAdditional            bool        `json:"Download_Additional"`
	GetParams                     []TaskParam `json:"GetParams,omitempty"`
	PostParams                    []TaskParam `json:"PostParams,omitempty"`
	HeaderParams                  []TaskParam `json:"HeaderParams,omitempty"`
	PrepareScript                 string      `json:"PrepareScript,omitempty"`
	DNSResolveMode                string      `json:"DNSResolveMode,omitempty"`
	DNSserverIP                   string      `json:"DNSserverIP,omitempty"`
	CustomDNSHosts                string      `json:"CustomDNSHosts,omitempty"`
	DeviceID                      int         `json:"Device_Id"`
	TaskTypeID                    int         `json:"Task_Type_Id"`
	Name                          string      `json:"Name"`
	Timeout                       int         `json:"Timeout,omitempty"`
}

// TaskParam ... simple struct for defining a param
type TaskParam struct {
	Name  string
	Value string
}

// CreateTaskResponseBlock ... struct for create task response
type CreateTaskResponseBlock struct {
	CreateResponseBlock
}

// UpdateTaskResponseBlock ... struct for update task response
type UpdateTaskResponseBlock struct {
	ResponseBlock
}

// DeleteTaskResponseBlock ... struct for delete task response
type DeleteTaskResponseBlock struct {
	ResponseBlock
}
