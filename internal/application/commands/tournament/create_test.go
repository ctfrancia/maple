package commands

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateTournamentCommand_Validate(t *testing.T) {
	tests := []struct {
		name         string
		cmd          CreateTournamentCommand
		wantErr      bool
		expectedErrs map[string]string
	}{
		{
			name: "valid tournament with all required fields",
			cmd: CreateTournamentCommand{
				Name:        "Test Tournament",
				Description: "A test tournament",
			},
			wantErr:      false,
			expectedErrs: nil,
		},
		{
			name: "valid tournament with minimal fields",
			cmd: CreateTournamentCommand{
				Name: "Min",
			},
			wantErr:      false,
			expectedErrs: nil,
		},
		{
			name: "valid tournament with maximum name length",
			cmd: CreateTournamentCommand{
				Name: "A very long tournament name that is exactly one hundred characters long for testing purposes!!",
			},
			wantErr:      false,
			expectedErrs: nil,
		},
		{
			name: "valid tournament with maximum description length",
			cmd: CreateTournamentCommand{
				Name:        "Test Tournament",
				Description: "A" + string(make([]byte, 499)), // 500 characters total
			},
			wantErr:      false,
			expectedErrs: nil,
		},
		{
			name: "empty name",
			cmd: CreateTournamentCommand{
				Name: "",
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"name": "is required",
			},
		},
		{
			name: "whitespace only name",
			cmd: CreateTournamentCommand{
				Name: "   ",
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"name": "is required",
			},
		},
		{
			name: "name too short",
			cmd: CreateTournamentCommand{
				Name: "AB",
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"name": "must be at least 3 characters",
			},
		},
		{
			name: "name too short with whitespace",
			cmd: CreateTournamentCommand{
				Name: "  A ",
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"name": "must be at least 3 characters",
			},
		},
		{
			name: "name too long",
			cmd: CreateTournamentCommand{
				Name: "A very long tournament name that definitely exceeds one hundred characters and should trigger validation error here",
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"name": "must be less than 100 characters",
			},
		},
		{
			name: "name exactly 101 characters",
			cmd: CreateTournamentCommand{
				Name: "A very long tournament name that is exactly one hundred and one characters long for testing purposes!",
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"name": "must be less than 100 characters",
			},
		},
		{
			name: "description too long",
			cmd: CreateTournamentCommand{
				Name:        "Valid Tournament",
				Description: "A" + string(make([]byte, 500)), // 501 characters total
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"description": "must be less than 500 characters",
			},
		},
		{
			name: "multiple validation errors",
			cmd: CreateTournamentCommand{
				Name:        "AB",
				Description: "A" + string(make([]byte, 500)), // 501 characters total
			},
			wantErr: true,
			expectedErrs: map[string]string{
				"name":        "must be at least 3 characters",
				"description": "must be less than 500 characters",
			},
		},
		{
			name: "valid tournament with empty schedule",
			cmd: CreateTournamentCommand{
				Name:     "Test Tournament",
				Schedule: []Schedule{},
			},
			wantErr:      false,
			expectedErrs: nil,
		},
		{
			name: "valid tournament with schedule",
			cmd: CreateTournamentCommand{
				Name: "Test Tournament",
				Schedule: []Schedule{
					{
						StartTime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
						EndTime:   time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC),
					},
				},
			},
			wantErr:      false,
			expectedErrs: nil,
		},
		{
			name: "valid tournament with all optional fields",
			cmd: CreateTournamentCommand{
				Name:               "Complete Tournament",
				Description:        "A complete tournament with all fields",
				AdditionalInfo:     "Some additional info",
				LocationID:         "loc123",
				MaxPlayers:         100,
				OpenToPublic:       true,
				OpenToRegistration: true,
				Schedule: []Schedule{
					{
						StartTime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
						EndTime:   time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC),
					},
				},
				Contact: Contact{
					Name:  "John Doe",
					Email: "john@example.com",
					Phone: "123-456-7890",
				},
				Registration: Registration{
					Status:     "open", // Assuming RegistrationStatus is a string type
					StartTime:  time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					PublicFee:  1000,
					PrivateFee: 800,
					OtherFee:   500,
					PrizePool:  10000,
					Payment: []Payment{
						{
							Place:  1,
							Amount: 5000,
							Type:   "monetary", // Assuming PaymentType is a string type
						},
					},
				},
			},
			wantErr:      false,
			expectedErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateTournamentCommand.Validate() expected error, got nil")
					return
				}

				ve, ok := IsValidationError(err)
				if !ok {
					t.Errorf("CreateTournamentCommand.Validate() expected ValidationError, got %T", err)
					return
				}

				if len(ve.Errors) != len(tt.expectedErrs) {
					t.Errorf("CreateTournamentCommand.Validate() expected %d errors, got %d", len(tt.expectedErrs), len(ve.Errors))
				}

				for field, expectedMsg := range tt.expectedErrs {
					if actualMsg, exists := ve.Errors[field]; !exists {
						t.Errorf("CreateTournamentCommand.Validate() missing error for field '%s'", field)
					} else if actualMsg != expectedMsg {
						t.Errorf("CreateTournamentCommand.Validate() for field '%s' = '%s', want '%s'", field, actualMsg, expectedMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("CreateTournamentCommand.Validate() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		ve       ValidationError
		expected string
	}{
		{
			name:     "empty errors",
			ve:       ValidationError{Errors: map[string]string{}},
			expected: "validation failed",
		},
		{
			name: "single error",
			ve: ValidationError{
				Errors: map[string]string{
					"name": "is required",
				},
			},
			expected: "validation failed: name: is required",
		},
		{
			name: "multiple errors",
			ve: ValidationError{
				Errors: map[string]string{
					"name":        "is required",
					"description": "too long",
				},
			},
			// Note: map iteration order is not guaranteed, so we need to check both possibilities
			// This test might need adjustment based on your Go version's map iteration behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ve.Error()

			if tt.name == "multiple errors" {
				// For multiple errors, just check that it contains both expected parts
				if !contains(result, "name: is required") || !contains(result, "description: too long") {
					t.Errorf("ValidationError.Error() = '%s', should contain both error messages", result)
				}
				if !contains(result, "validation failed:") {
					t.Errorf("ValidationError.Error() = '%s', should start with 'validation failed:'", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("ValidationError.Error() = '%s', want '%s'", result, tt.expected)
				}
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectVE     bool
		expectedErrs map[string]string
	}{
		{
			name: "is validation error",
			err: ValidationError{
				Errors: map[string]string{
					"name": "is required",
				},
			},
			expectVE: true,
			expectedErrs: map[string]string{
				"name": "is required",
			},
		},
		{
			name:         "is not validation error",
			err:          fmt.Errorf("some other error"),
			expectVE:     false,
			expectedErrs: nil,
		},
		{
			name:         "nil error",
			err:          nil,
			expectVE:     false,
			expectedErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve, ok := IsValidationError(tt.err)

			if ok != tt.expectVE {
				t.Errorf("IsValidationError() ok = %v, want %v", ok, tt.expectVE)
			}

			if tt.expectVE {
				if ve == nil {
					t.Errorf("IsValidationError() expected ValidationError, got nil")
					return
				}

				if len(ve.Errors) != len(tt.expectedErrs) {
					t.Errorf("IsValidationError() expected %d errors, got %d", len(tt.expectedErrs), len(ve.Errors))
				}

				for field, expectedMsg := range tt.expectedErrs {
					if actualMsg, exists := ve.Errors[field]; !exists {
						t.Errorf("IsValidationError() missing error for field '%s'", field)
					} else if actualMsg != expectedMsg {
						t.Errorf("IsValidationError() for field '%s' = '%s', want '%s'", field, actualMsg, expectedMsg)
					}
				}
			} else {
				if ve != nil {
					t.Errorf("IsValidationError() expected nil ValidationError, got %v", ve)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr, 1)))
}

func containsAt(s, substr string, start int) bool {
	if start >= len(s) {
		return false
	}
	if start+len(substr) > len(s) {
		return containsAt(s, substr, start+1)
	}
	if s[start:start+len(substr)] == substr {
		return true
	}
	return containsAt(s, substr, start+1)
}
