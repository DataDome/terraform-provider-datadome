package datadome

type HttpResponse struct {
	Data    CustomRules `json:"data"`
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

type CustomRules struct {
	CustomRules []CustomRule `json:"custom_rules"`
}

type CustomRule struct {
	ID           int    `json:"id"`
	Name         string `json:"rule_name"`
	Response     string `json:"rule_response"`
	Query        string `json:"query"`
	IPStart      string `json:"ip_start"`
	IPEnd        string `json:"ip_end"`
	EndpointType string `json:"endpoint_type"`
	Priority     string `json:"rule_priority"`
	Hits         int    `json:"hits"`
}
