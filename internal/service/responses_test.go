package service

import (
	"testing"

	"platform/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestIsOfferRespondable(t *testing.T) {
	t.Parallel()

	assert.True(t, isOfferRespondable(models.OfferStatusSent))
	assert.False(t, isOfferRespondable(models.OfferStatusAccepted))
	assert.False(t, isOfferRespondable(models.OfferStatusDeclined))
	assert.False(t, isOfferRespondable(models.OfferStatusExpired))
}
