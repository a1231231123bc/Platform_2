package dto

// CreateJobRequest is the request body for POST /jobs.
type CreateJobRequest struct {
	Title         string   `json:"title" validate:"required,min=3,max=255"`
	Description   *string  `json:"description"`
	Region        string   `json:"region" validate:"required"`
	EquipmentType *string  `json:"equipmentType"`
	Volume        *string  `json:"volume"`
	Deadline      *string  `json:"deadline"`
	Price         *float64 `json:"price"`
	Conditions    *string  `json:"conditions"`
}

// UpdateJobRequest is the request body for PATCH /jobs/:id.
type UpdateJobRequest struct {
	Title         *string  `json:"title" validate:"omitempty,min=3,max=255"`
	Description   *string  `json:"description"`
	Region        *string  `json:"region"`
	EquipmentType *string  `json:"equipmentType"`
	Volume        *string  `json:"volume"`
	Deadline      *string  `json:"deadline"`
	Price         *float64 `json:"price"`
	Conditions    *string  `json:"conditions"`
}

// QueryJobsRequest is the query params for GET /jobs.
type QueryJobsRequest struct {
	Status *string `json:"status" validate:"omitempty,oneof=DRAFT ACTIVE IN_PROGRESS COMPLETED CANCELLED"`
	Region *string `json:"region"`
	Page   int     `json:"page" validate:"min=1"`
	Limit  int     `json:"limit" validate:"min=1,max=100"`
}
