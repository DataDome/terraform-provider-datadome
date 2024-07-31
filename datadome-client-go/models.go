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
	Data CustomRule `json:"data"`
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
