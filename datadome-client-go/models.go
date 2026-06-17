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
	ID            *int           `json:"id"`
	Name          string         `json:"rule_name"`
	Response      string         `json:"rule_response"`
	Query         string         `json:"query"`
	EndpointType  string         `json:"endpoint_type"`
	Priority      string         `json:"rule_priority"`
	Enabled       *bool          `json:"rule_enabled,omitempty"`
	ActivatedAt   *string        `json:"activated_at,omitempty"`
	ExpiredAt     *string        `json:"expired_at,omitempty"`
	OverriddenBot *OverriddenBot `json:"overridden_bot,omitempty"`
	PolicyOptions *PolicyOptions `json:"policy_options,omitempty"`
}

// OverriddenBot identifies a Verified Bot or AI Agent that the custom rule applies to.
type OverriddenBot struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// PolicyOptions holds an optional rate-limit or time-box policy for a custom rule.
// At most one of TimeBox or RateLimit may be set.
// Only valid when the rule response is "allow" or "intent_based".
type PolicyOptions struct {
	TimeBox   *TimeBoxOptions   `json:"time_box,omitempty"`
	RateLimit *RateLimitOptions `json:"rate_limit,omitempty"`
}

// TimeBoxOptions restricts a rule to specific hours of the week, applying an
// alternative response outside the authorized window.
type TimeBoxOptions struct {
	AuthorizedHoursOfTheWeek []int  `json:"authorized_hours_of_the_week"`
	ResponseOutsideTimeBox   string `json:"response_outside_time_box"`
}

// RateLimitOptions triggers an alternative response once a request threshold is
// exceeded within a time window.
type datadome/resource_custom_rule.go struct {
	AppliesTo              string `json:"applies_to"`
	Threshold              int    `json:"threshold"`
	TimeFrame              string `json:"time_frame"`
	ResponseAfterThreshold string `json:"response_after_threshold"`
}

// Endpoint structure containing the information of an endpoint
type Endpoint struct {
	ID                 *string `json:"id"`
	Name               string  `json:"name"`
	Description        *string `json:"description"`
	PositionBefore     *string `json:"positionBefore"`
	TrafficUsage       string  `json:"trafficUsage"`
	Source             string  `json:"source"`
	CookieSameSite     string  `json:"cookieSameSite"`
	Domain             *string `json:"domain"`
	PathInclusion      *string `json:"pathInclusion"`
	PathExclusion      *string `json:"pathExclusion"`
	UserAgentInclusion *string `json:"userAgentInclusion"`
	Query              *string `json:"query"`
	ResponseFormat     string  `json:"responseFormat"`
	DetectionEnabled   bool    `json:"detectionEnabled"`
	ProtectionEnabled  bool    `json:"protectionEnabled"`
}
