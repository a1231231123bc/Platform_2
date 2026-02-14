package service

import (
	"testing"

	"platform/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestParseCommunicationDirection(t *testing.T) {
	outgoing, err := parseCommunicationDirection("OUTGOING")
	assert.NoError(t, err)
	assert.Equal(t, models.CommunicationDirectionOutgoing, outgoing)

	incoming, err := parseCommunicationDirection("INCOMING")
	assert.NoError(t, err)
	assert.Equal(t, models.CommunicationDirectionIncoming, incoming)

	_, err = parseCommunicationDirection("INVALID")
	assert.Error(t, err)
}
