package dto

// UpdateOrganizationRequest is the request body for PATCH /organizations/:id.
type UpdateOrganizationRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=2,max=255"`
	INN          *string `json:"inn"`
	ContactEmail *string `json:"contactEmail" validate:"omitempty,email"`
	ContactPhone *string `json:"contactPhone"`
}
