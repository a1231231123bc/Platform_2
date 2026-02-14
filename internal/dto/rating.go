package dto

// CreateRatingRequest is the request body for POST /ratings.
type CreateRatingRequest struct {
	JobID        string  `json:"jobId" validate:"required,uuid"`
	ContractorID string  `json:"contractorId" validate:"required,uuid"`
	Score        int     `json:"score" validate:"required,min=1,max=5"`
	Comment      *string `json:"comment"`
}
