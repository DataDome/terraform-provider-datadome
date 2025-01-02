package datadome

// HttpResponse from the DataDome's API
type HttpResponse struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Errors  []Error     `json:"errors"`
	Message string      `json:"message"`
}

// HttpRequest with the CustomRule inside
type HttpRequest struct {
	Data interface{} `json:"data"`
}

// Error structure returned in case of HTTP error
type Error struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

// ID format
type ID struct {
	ID int `json:"id"`
}

// CustomRules structure containing a slice of CustomRule
type CustomRules struct {
	CustomRules []CustomRule `json:"custom_rules"`
}

// CustomRule structure containing the information of a custom rule
type CustomRule struct {
	ID           int    `json:"id"`
	Name         string `json:"rule_name"`
	Response     string `json:"rule_response"`
	Query        string `json:"query"`
	EndpointType string `json:"endpoint_type"`
	Priority     string `json:"rule_priority"`
	Enabled      bool   `json:"rule_enabled"`
}

// Endpoint structure containing the information of an endpoint
type Endpoint struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Description        *string `json:"description,omitempty"`
	PositionBefore     *string `json:"positionBefore,omitempty"`
	TrafficUsage       string  `json:"trafficUsage"`
	Source             string  `json:"source"`
	CookieSameSite     string  `json:"cookieSameSite"`
	Domain             *string `json:"domain,omitempty"`
	PathInclusion      *string `json:"pathInclusion,omitempty"`
	PathExclusion      *string `json:"pathExclusion,omitempty"`
	UserAgentInclusion *string `json:"userAgentInclusion,omitempty"`
	ResponseFormat     string  `json:"responseFormat"`
	DetectionEnabled   bool    `json:"detectionEnabled"`
	ProtectionEnabled  bool    `json:"protectionEnabled"`
}
