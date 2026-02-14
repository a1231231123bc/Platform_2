package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAssignmentRequestValidation(t *testing.T) {
	req := CreateAssignmentRequest{
		JobID:        "3f026267-f3ce-4b70-9ee5-6d82aa4f4f78",
		ContractorID: "de1178db-8326-4f3d-884d-4ea5e18e2a26",
	}
	assert.NoError(t, ValidateStruct(req))
}

func TestUpdateAssignmentStatusRequestValidation(t *testing.T) {
	okReq := UpdateAssignmentStatusRequest{Status: "CONFIRMED"}
	assert.NoError(t, ValidateStruct(okReq))

	badReq := UpdateAssignmentStatusRequest{Status: "INVALID"}
	err := ValidateStruct(badReq)
	assert.Error(t, err)
	assert.Contains(t, ValidationErrors(err), "status")
}
