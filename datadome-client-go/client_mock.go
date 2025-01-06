package datadome

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

// MockClientCustomRule structure for test purposes on the custom rules
type MockClientCustomRule struct {
	CreateFunc func(ctx context.Context, params CustomRule) (*int, error)
	ReadFunc   func(ctx context.Context, id int) (*CustomRule, error)
	UpdateFunc func(ctx context.Context, params CustomRule) (*CustomRule, error)
	DeleteFunc func(ctx context.Context, id int) error

	resources map[int]*CustomRule
}

// NewMockClientCustomRule returns a new MockClient for custom rule management
func NewMockClientCustomRule() *MockClientCustomRule {
	return &MockClientCustomRule{
		resources: make(map[int]*CustomRule),
	}
}

// Create mock method
func (m *MockClientCustomRule) Create(ctx context.Context, params CustomRule) (*int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, params)
	}

	if params.ID == nil {
		ID := rand.Int()
		params.ID = &ID
	}

	newResource := &params
	m.resources[*newResource.ID] = newResource
	return newResource.ID, nil
}

// Read mock method
func (m *MockClientCustomRule) Read(ctx context.Context, id int) (*CustomRule, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(ctx, id)
	}

	var value *CustomRule
	for _, v := range m.resources {
		if *v.ID == id {
			value = v
		}
	}

	return value, nil
}

// Update mock method
func (m *MockClientCustomRule) Update(ctx context.Context, params CustomRule) (*CustomRule, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, params)
	}

	_, exists := m.resources[*params.ID]
	if !exists {
		return nil, fmt.Errorf("resource not found with ID %d", params.ID)
	}

	m.resources[*params.ID] = &params
	return &params, nil
}

// Delete mock method
func (m *MockClientCustomRule) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	_, exists := m.resources[id]
	if !exists {
		return fmt.Errorf("resource not found with ID %d", id)
	}

	delete(m.resources, id)
	return nil
}

// MockClientEndpoint structure for test purposes on the endpoints
type MockClientEndpoint struct {
	CreateFunc func(ctx context.Context, params Endpoint) (*string, error)
	ReadFunc   func(ctx context.Context, id string) (*Endpoint, error)
	UpdateFunc func(ctx context.Context, params Endpoint) (*Endpoint, error)
	DeleteFunc func(ctx context.Context, id string) error

	resources map[string]*Endpoint
}

// NewMockClientCustomRule returns a new MockClient for custom rule management
func NewMockClientEndpoint() *MockClientEndpoint {
	return &MockClientEndpoint{
		resources: make(map[string]*Endpoint),
	}
}

// Create mock method
func (m *MockClientEndpoint) Create(ctx context.Context, params Endpoint) (*string, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, params)
	}

	if params.ID == nil {
		newUUID := uuid.New()
		ID := newUUID.String()
		params.ID = &ID
	}

	newResource := &params
	m.resources[*newResource.ID] = newResource
	return newResource.ID, nil
}

// Read mock method
func (m *MockClientEndpoint) Read(ctx context.Context, id string) (*Endpoint, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(ctx, id)
	}

	var value *Endpoint
	for _, v := range m.resources {
		if *v.ID == id {
			value = v
		}
	}

	return value, nil
}

// Update mock method
func (m *MockClientEndpoint) Update(ctx context.Context, params Endpoint) (*Endpoint, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, params)
	}

	_, exists := m.resources[*params.ID]
	if !exists {
		return nil, fmt.Errorf("resource not found with ID %s", *params.ID)
	}

	m.resources[*params.ID] = &params
	return &params, nil
}

// Delete mock method
func (m *MockClientEndpoint) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	_, exists := m.resources[id]
	if !exists {
		return fmt.Errorf("resource not found with ID %s", id)
	}

	delete(m.resources, id)
	return nil
}
