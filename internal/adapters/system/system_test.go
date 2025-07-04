package system

import (
	"testing"

	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSystemAdapter(t *testing.T) {
	t.Run("should create new system adapter instance", func(t *testing.T) {
		// Act
		adapter := NewSystemAdapter()

		// Assert
		assert.NotNil(t, adapter)
		assert.IsType(t, &SystemAdapter{}, adapter)
	})
}

func TestSystemAdapter_GetSystemInfo(t *testing.T) {
	t.Run("should return system info with correct version", func(t *testing.T) {
		// Arrange
		adapter := NewSystemAdapter()
		expectedVersion := "1.0.0"

		// Act
		result := adapter.GetSystemInfo()

		// Assert
		assert.Equal(t, expectedVersion, result.Version)
		assert.IsType(t, domain.System{}, result)
	})

	t.Run("should return consistent system info on multiple calls", func(t *testing.T) {
		// Arrange
		adapter := NewSystemAdapter()

		// Act
		result1 := adapter.GetSystemInfo()
		result2 := adapter.GetSystemInfo()

		// Assert
		assert.Equal(t, result1.Version, result2.Version)
		assert.Equal(t, result1, result2)
	})

	t.Run("should return non-empty version", func(t *testing.T) {
		// Arrange
		adapter := NewSystemAdapter()

		// Act
		result := adapter.GetSystemInfo()

		// Assert
		assert.NotEmpty(t, result.Version)
	})
}

func TestSystemAdapter_Interface(t *testing.T) {
	t.Run("should implement expected interface methods", func(t *testing.T) {
		// Arrange
		adapter := NewSystemAdapter()

		// Act & Assert - this will fail at compile time if interface is not implemented
		var _ interface {
			GetSystemInfo() domain.System
		} = adapter

		// Additional runtime check
		require.NotPanics(t, func() {
			adapter.GetSystemInfo()
		})
	})
}

// Table-driven test example (useful if you add more fields to System)
func TestSystemAdapter_GetSystemInfo_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		expected domain.System
	}{
		{
			name: "should return expected system information",
			expected: domain.System{
				Version: "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			adapter := NewSystemAdapter()

			// Act
			result := adapter.GetSystemInfo()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}
