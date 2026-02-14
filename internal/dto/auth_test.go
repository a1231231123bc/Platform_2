package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterRequest_Valid(t *testing.T) {
	req := RegisterRequest{
		OrganizationName: "ООО Ромашка",
		AdminName:        "Иван Иванов",
		AdminEmail:       "admin@romashka.ru",
		AdminPassword:    "SecurePass123!",
	}
	err := ValidateStruct(req)
	assert.NoError(t, err)
}

func TestRegisterRequest_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		req  RegisterRequest
		field string
	}{
		{
			name:  "missing organization name",
			req:   RegisterRequest{AdminName: "Admin", AdminEmail: "a@b.com", AdminPassword: "12345678"},
			field: "organizationName",
		},
		{
			name:  "missing admin email",
			req:   RegisterRequest{OrganizationName: "Org", AdminName: "Admin", AdminPassword: "12345678"},
			field: "adminEmail",
		},
		{
			name:  "missing admin password",
			req:   RegisterRequest{OrganizationName: "Org", AdminName: "Admin", AdminEmail: "a@b.com"},
			field: "adminPassword",
		},
		{
			name:  "missing admin name",
			req:   RegisterRequest{OrganizationName: "Org", AdminEmail: "a@b.com", AdminPassword: "12345678"},
			field: "adminName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.req)
			assert.Error(t, err)
			errs := ValidationErrors(err)
			assert.Contains(t, errs, tt.field)
		})
	}
}

func TestRegisterRequest_InvalidEmail(t *testing.T) {
	req := RegisterRequest{
		OrganizationName: "Org",
		AdminName:        "Admin",
		AdminEmail:       "not-an-email",
		AdminPassword:    "12345678",
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "adminEmail")
}

func TestRegisterRequest_ShortPassword(t *testing.T) {
	req := RegisterRequest{
		OrganizationName: "Org",
		AdminName:        "Admin",
		AdminEmail:       "a@b.com",
		AdminPassword:    "short",
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "adminPassword")
}

func TestLoginRequest_Valid(t *testing.T) {
	req := LoginRequest{
		Email:    "admin@romashka.ru",
		Password: "SecurePass123!",
	}
	err := ValidateStruct(req)
	assert.NoError(t, err)
}

func TestLoginRequest_Invalid(t *testing.T) {
	req := LoginRequest{
		Email:    "not-email",
		Password: "short",
	}
	err := ValidateStruct(req)
	assert.Error(t, err)
	errs := ValidationErrors(err)
	assert.Contains(t, errs, "email")
	assert.Contains(t, errs, "password")
}
