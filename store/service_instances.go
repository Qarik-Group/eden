package store

import (
	"time"
)

// ServiceInstances is the set of all service instances created with this CLI
type ServiceInstances struct {
	Instances []ServiceInstance `json:"service_instances"`
}

// ServiceInstance represents a service instance that was created with this CLI
type ServiceInstance struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	BrokerURL string           `json:"broker_url"`
	Bindings  []ServiceBinding `json:"bindings"`
	CreatedAt *time.Location   `json:"created_at"`
}

// ServiceBinding represents a binding with credentials
type ServiceBinding struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Credentials interface{}    `json:"credentials"`
	CreatedAt   *time.Location `json:"created_at"`
}
