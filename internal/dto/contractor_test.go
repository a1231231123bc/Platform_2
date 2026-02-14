package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContractorRequest_Valid(t *testing.T) {
	req := CreateContractorRequest{
		Name:           "Иван Иванов",
		Phone:          "+79991234567",
		Type:           "INDIVIDUAL",
		Regions:        []string{"Moscow"},
		EquipmentTypes: []string{"Самосвал"},
		ConsentGiven:   true,
	}
	err := ValidateStruct(req)
	assert.NoError(t, err)
}

func TestCreateContractorRequest_MissingRequired(t *testing.T) {
	req := CreateContractorRequest{}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "name")
	assert.Contains(t, errs, "phone")
	assert.Contains(t, errs, "type")
	assert.Contains(t, errs, "regions")
	assert.Contains(t, errs, "equipmentTypes")
}

func TestCreateContractorRequest_InvalidType(t *testing.T) {
	req := CreateContractorRequest{
		Name:           "Test",
		Phone:          "+79991234567",
		Type:           "INVALID",
		Regions:        []string{"Moscow"},
		EquipmentTypes: []string{"Самосвал"},
		ConsentGiven:   true,
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "type")
}

func TestCreateContractorRequest_InvalidEmail(t *testing.T) {
	email := "not-email"
	req := CreateContractorRequest{
		Name:           "Test",
		Phone:          "+79991234567",
		Email:          &email,
		Type:           "INDIVIDUAL",
		Regions:        []string{"Moscow"},
		EquipmentTypes: []string{"Самосвал"},
		ConsentGiven:   true,
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "email")
}

func TestQueryContractorsRequest_ValidType(t *testing.T) {
	tp := "IP"
	req := QueryContractorsRequest{
		Type:  &tp,
		Page:  1,
		Limit: 20,
	}
	err := ValidateStruct(req)
	assert.NoError(t, err)
}

func TestQueryContractorsRequest_InvalidType(t *testing.T) {
	tp := "INVALID"
	req := QueryContractorsRequest{
		Type:  &tp,
		Page:  1,
		Limit: 20,
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "type")
}
