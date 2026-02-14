package dto

// CreateContractorRequest is the request body for POST /contractors.
type CreateContractorRequest struct {
	Name              string   `json:"name" validate:"required"`
	Phone             string   `json:"phone" validate:"required"`
	Email             *string  `json:"email" validate:"omitempty,email"`
	Type              string   `json:"type" validate:"required,oneof=INDIVIDUAL IP LLC"`
	Regions           []string `json:"regions" validate:"required"`
	EquipmentTypes    []string `json:"equipmentTypes" validate:"required"`
	PriceExpectations *string  `json:"priceExpectations"`
	IsAvailable       *bool    `json:"isAvailable"`
	ConsentGiven      bool     `json:"consentGiven"`
}

// UpdateContractorRequest is the request body for PATCH /contractors/:id.
type UpdateContractorRequest struct {
	Name              *string  `json:"name"`
	Phone             *string  `json:"phone"`
	Email             *string  `json:"email" validate:"omitempty,email"`
	Type              *string  `json:"type" validate:"omitempty,oneof=INDIVIDUAL IP LLC"`
	Regions           []string `json:"regions"`
	EquipmentTypes    []string `json:"equipmentTypes"`
	PriceExpectations *string  `json:"priceExpectations"`
	IsAvailable       *bool    `json:"isAvailable"`
	ConsentGiven      *bool    `json:"consentGiven"`
}

// QueryContractorsRequest is the query params for GET /contractors.
type QueryContractorsRequest struct {
	Region      *string `json:"region"`
	Type        *string `json:"type" validate:"omitempty,oneof=INDIVIDUAL IP LLC"`
	IsAvailable *bool   `json:"isAvailable"`
	Search      *string `json:"search"`
	Page        int     `json:"page" validate:"min=1"`
	Limit       int     `json:"limit" validate:"min=1,max=100"`
}
