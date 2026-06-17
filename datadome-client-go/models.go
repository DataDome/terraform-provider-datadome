package datadome

import "encoding/json"

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
// On writes the API expects a bare UUID string; on reads it returns an object {uuid, name}.
// The custom MarshalJSON / UnmarshalJSON bridge that asymmetry transparently.
type OverriddenBot struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// MarshalJSON emits a bare UUID string to satisfy the write schema.
func (o OverriddenBot) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.UUID)
}

// UnmarshalJSON accepts either a UUID string (write echo) or an object {uuid, name} (read schema).
func (o *OverriddenBot) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		o.UUID = s
		return nil
	}
	type alias OverriddenBot
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*o = OverriddenBot(a)
	return nil
}

type PolicyOptions struct {
	TimeBox   *TimeBoxOptions   `json:"time_box,omitempty"`
	RateLimit *RateLimitOptions `json:"rate_limit,omitempty"`
}

type TimeBoxOptions struct {
	AuthorizedHoursOfTheWeek []int  `json:"authorized_hours_of_the_week"`
	ResponseOutsideTimeBox   string `json:"response_outside_time_box"`
}

type RateLimitOptions struct {
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
