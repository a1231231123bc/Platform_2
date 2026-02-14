package dto

// UpdateUserRequest is the request body for PATCH /users/:id.
type UpdateUserRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2,max=255"`
	Role *string `json:"role" validate:"omitempty,oneof=ADMIN MANAGER OBSERVER"`
}
