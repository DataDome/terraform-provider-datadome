package datadome

import (
	"encoding/json"
	"testing"
)

// TestOverriddenBotMarshalJSON verifies that OverriddenBot is serialized as a
// bare UUID string (matching the create/update request schema).
func TestOverriddenBotMarshalJSON(t *testing.T) {
	ob := OverriddenBot{UUID: "550e8400-e29b-41d4-a716-446655440000", Name: "My Test Bot"}
	data, err := json.Marshal(ob)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(data)
	want := `"550e8400-e29b-41d4-a716-446655440000"`
	if got != want {
		t.Errorf("MarshalJSON = %s, want %s", got, want)
	}
}

// TestCustomRuleMarshalJSON verifies that a CustomRule with a non-nil
// OverriddenBot embeds the UUID string (not an object) in the JSON payload
// sent to the create/update API.
func TestCustomRuleMarshalJSON(t *testing.T) {
	id := 42
	rule := CustomRule{
		ID:       &id,
		Name:     "My rule",
		Response: "allow",
		Query:    "ip:1.2.3.4",
		OverriddenBot: &OverriddenBot{
			UUID: "550e8400-e29b-41d4-a716-446655440000",
			Name: "My Test Bot",
		},
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Deserialize into a map to inspect the overridden_bot field type.
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	val, ok := m["overridden_bot"]
	if !ok {
		t.Fatal("overridden_bot field missing from marshalled JSON")
	}
	strVal, ok := val.(string)
	if !ok {
		t.Fatalf("overridden_bot should be a string, got %T: %v", val, val)
	}
	if strVal != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("overridden_bot = %q, want %q", strVal, "550e8400-e29b-41d4-a716-446655440000")
	}
}

// TestOverriddenBotUnmarshalJSON verifies that an API response carrying
// overridden_bot as an object {uuid, name} is decoded correctly into the struct.
func TestOverriddenBotUnmarshalJSON(t *testing.T) {
	payload := `{"overridden_bot":{"uuid":"550e8400-e29b-41d4-a716-446655440000","name":"My Test Bot"}}`

	var rule struct {
		OverriddenBot *OverriddenBot `json:"overridden_bot"`
	}
	if err := json.Unmarshal([]byte(payload), &rule); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.OverriddenBot == nil {
		t.Fatal("OverriddenBot should not be nil")
	}
	if rule.OverriddenBot.UUID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("UUID = %q, want %q", rule.OverriddenBot.UUID, "550e8400-e29b-41d4-a716-446655440000")
	}
	if rule.OverriddenBot.Name != "My Test Bot" {
		t.Errorf("Name = %q, want %q", rule.OverriddenBot.Name, "My Test Bot")
	}
}

// TestCustomRuleNilOverriddenBot verifies that when OverriddenBot is nil,
// the field is omitted from the serialized payload.
func TestCustomRuleNilOverriddenBot(t *testing.T) {
	id := 1
	rule := CustomRule{ID: &id, Name: "r", Response: "block", Query: "url:*"}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if _, ok := m["overridden_bot"]; ok {
		t.Error("overridden_bot should be absent when nil, but was present")
	}
}
