package repository

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardMetricsJSONContract(t *testing.T) {
	t.Parallel()

	metrics := DashboardMetrics{
		TotalJobs:        10,
		ActiveJobs:       3,
		TotalContractors: 7,
		TotalAssignments: 15,
	}

	data, err := json.Marshal(metrics)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Contains(t, m, "totalJobs")
	assert.Contains(t, m, "activeJobs")
	assert.Contains(t, m, "totalContractors")
	assert.Contains(t, m, "totalAssignments")
}
