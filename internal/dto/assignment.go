package dto

// CreateAssignmentRequest is the request body for POST /assignments.
type CreateAssignmentRequest struct {
	JobID        string `json:"jobId" validate:"required,uuid"`
	ContractorID string `json:"contractorId" validate:"required,uuid"`
}

// UpdateAssignmentStatusRequest is the request body for PATCH /assignments/:id/status.
type UpdateAssignmentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=ASSIGNED CONFIRMED IN_WORK COMPLETED CANCELLED"`
}
