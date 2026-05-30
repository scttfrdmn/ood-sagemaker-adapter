// Package ood defines the OOD interactive app spec types for the SageMaker adapter.
package ood

// AppSpec is the OOD interactive app launch payload.
type AppSpec struct {
	AppName  string            `json:"app_name"`
	AppType  string            `json:"app_type,omitempty"` // JupyterServer, KernelGateway, etc.
	UserName string            `json:"user_name"`
	DomainID string            `json:"domain_id,omitempty"`
	Env      map[string]string `json:"env,omitempty"`
}

// AppStatus maps SageMaker app states to OOD status strings.
type AppStatus struct {
	ID         string `json:"id"`
	Status     string `json:"status"` // queued, running, completed, failed
	PresignURL string `json:"presign_url,omitempty"`
	Message    string `json:"message,omitempty"`
}

const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusUnknown   = "undetermined"
)
