package contract

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestRegisterContract(t *testing.T) {
	provider := auth.NewLocalAuthProvider()
	//router := http.NewRouter()

	// Register handler would be attached here
	// For now, we'll test the contract directly against the provider

	tests := []struct {
		name          string
		request       auth.RegisterUserRequest
		shouldSucceed bool
	}{
		{
			name: "Valid registration request",
			request: auth.RegisterUserRequest{
				FullName:        "John Doe",
				DateOfBirth:     "1990-01-15",
				MobileNumber:    "+1234567890",
				Email:           "john.doe@example.com",
				Username:        "johndoe",
				Password:        "SecurePass123!",
				ConfirmPassword: "SecurePass123!",
			},
			shouldSucceed: true,
		},
		{
			name: "Missing required fields",
			request: auth.RegisterUserRequest{
				FullName:        "",
				DateOfBirth:     "1990-01-15",
				MobileNumber:    "+1234567890",
				Email:           "john.doe@example.com",
				Username:        "johndoe",
				Password:        "SecurePass123!",
				ConfirmPassword: "SecurePass123!",
			},
			shouldSucceed: false,
		},
		{
			name: "Password mismatch",
			request: auth.RegisterUserRequest{
				FullName:        "John Doe",
				DateOfBirth:     "1990-01-15",
				MobileNumber:    "+1234567890",
				Email:           "john.doe@example.com",
				Username:        "johndoe",
				Password:        "SecurePass123!",
				ConfirmPassword: "DifferentPass123!",
			},
			shouldSucceed: false,
		},
		{
			name: "Invalid email format",
			request: auth.RegisterUserRequest{
				FullName:        "John Doe",
				DateOfBirth:     "1990-01-15",
				MobileNumber:    "+1234567890",
				Email:           "invalid-email",
				Username:        "johndoe",
				Password:        "SecurePass123!",
				ConfirmPassword: "SecurePass123!",
			},
			shouldSucceed: false,
		},
		{
			name: "Invalid username format",
			request: auth.RegisterUserRequest{
				FullName:        "John Doe",
				DateOfBirth:     "1990-01-15",
				MobileNumber:    "+1234567890",
				Email:           "john.doe@example.com",
				Username:        "jo", // Too short
				Password:        "SecurePass123!",
				ConfirmPassword: "SecurePass123!",
			},
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the provider directly for contract validation
			_, err := provider.Register(tt.request)

			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error but got success")
				}
			}
		})
	}
}

func TestRegisterResponseContract(t *testing.T) {
	provider := auth.NewLocalAuthProvider()
	request := auth.RegisterUserRequest{
		FullName:        "Jane Doe",
		DateOfBirth:     "1985-05-20",
		MobileNumber:    "+1987654321",
		Email:           "jane.doe@example.com",
		Username:        "janedoe",
		Password:        "SecurePass456!",
		ConfirmPassword: "SecurePass456!",
	}

	response, err := provider.Register(request)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	// Verify response structure
	if response.UserID == "" {
		t.Error("UserID should not be empty")
	}
	if response.Username == "" {
		t.Error("Username should not be empty")
	}
	if response.FirstName == "" {
		t.Error("FirstName should not be empty")
	}
	if response.PasswordHash == "" {
		t.Error("PasswordHash should not be empty")
	}
}
