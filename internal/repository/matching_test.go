package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildMatchingFiltersWithEquipment(t *testing.T) {
	t.Parallel()

	equipment := "Excavator"
	where, args := buildMatchingFilters("Moscow", &equipment)

	assert.Contains(t, where, "c.is_available = true")
	assert.Contains(t, where, "$1 = ANY(c.regions)")
	assert.Contains(t, where, "c.equipment_types @> ARRAY[$2]::text[]")
	require.Len(t, args, 2)
	assert.Equal(t, "Moscow", args[0])
	assert.Equal(t, "Excavator", args[1])
}

func TestBuildMatchingFiltersWithoutEquipment(t *testing.T) {
	t.Parallel()

	where, args := buildMatchingFilters("Moscow", nil)

	assert.Contains(t, where, "c.is_available = true")
	assert.Contains(t, where, "$1 = ANY(c.regions)")
	assert.NotContains(t, where, "equipment_types")
	require.Len(t, args, 1)
	assert.Equal(t, "Moscow", args[0])
}
