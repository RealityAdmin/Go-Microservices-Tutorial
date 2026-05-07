package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Generic interface for service registry
type Registry interface {
	// Create a server instance record in registry
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error

	// Remove a service instance record from registry
	Deregister(ctx context.Context, instanceID string, serviceName string) error

	// Return list of addresses of active instances
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)

	// Push mechanism to report healthy state to registry
	ReportHealthyState(instanceID string, serviceName string) error
}

var ErrNotFound = errors.New("No service addresses found")

// Generate pseudo-random service instance ID.: SERVICE_NAME-RANDOM NO
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
