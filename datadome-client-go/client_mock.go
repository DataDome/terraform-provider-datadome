package datadome

import (
	"context"
	"fmt"
)

// MockClient structure for test purposes
type MockClient struct {
	CreateFunc func(ctx context.Context, params CustomRule) (*int, error)
	ReadFunc   func(ctx context.Context) ([]CustomRule, error)
	UpdateFunc func(ctx context.Context, params CustomRule) (*CustomRule, error)
	DeleteFunc func(ctx context.Context, id int) error

	resources map[int]*CustomRule
}

// NewMockClient returns a new MockClient
func NewMockClient() *MockClient {
	return &MockClient{
		resources: make(map[int]*CustomRule),
	}
}

// Create mock method
func (m *MockClient) Create(ctx context.Context, params CustomRule) (*int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, params)
	}

	newResource := &params
	m.resources[newResource.ID] = newResource
	return &newResource.ID, nil
}

// Read mock method
func (m *MockClient) Read(ctx context.Context) ([]CustomRule, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(ctx)
	}

	values := []CustomRule{}
	for _, value := range m.resources {
		values = append(values, *value)
	}

	return values, nil
}

// Update mock method
func (m *MockClient) Update(ctx context.Context, params CustomRule) (*CustomRule, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, params)
	}

	_, exists := m.resources[params.ID]
	if !exists {
		return nil, fmt.Errorf("resource not found with ID %d", params.ID)
	}

	m.resources[params.ID] = &params
	return &params, nil
}

// Delete mock method
func (m *MockClient) Delete(ctx context.Context, id int) error {
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
