package datadome

type HttpResponse struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Errors  []Error     `json:"errors"`
	Message string      `json:"message"`
}

type HttpRequest struct {
	Data CustomRule `json:"data"`
}

type Error struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

type ID struct {
	ID int `json:"id"`
}

type CustomRules struct {
	CustomRules []CustomRule `json:"custom_rules"`
}

type CustomRule struct {
	ID           int    `json:"id"`
	Name         string `json:"rule_name"`
	Response     string `json:"rule_response"`
	Query        string `json:"query"`
	EndpointType string `json:"endpoint_type"`
	Priority     string `json:"rule_priority"`
	Enabled      bool   `json:"rule_enabled"`
}
