package validation

import (
	"testing"

	"artha-kosha/apps/finance-api/internal/auth"
)

func TestRegistrationValidation_RequiredFields(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	tests := []struct {
		name        string
		request     auth.RegisterUserRequest
		expectError bool
	}{
		{
			name: "Missing full name",
			request: auth.RegisterUserRequest{
				FullName:        "",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			},
			expectError: true,
		},
		{
			name: "Missing date of birth",
			request: auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "",
				MobileNumber:    "+1234567890",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			},
			expectError: true,
		},
		{
			name: "Missing mobile number",
			request: auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			},
			expectError: true,
		},
		{
			name: "Missing email",
			request: auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           "",
				Username:        "testuser",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			},
			expectError: true,
		},
		{
			name: "Missing username",
			request: auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           "test@example.com",
				Username:        "",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			},
			expectError: true,
		},
		{
			name: "Missing password",
			request: auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "",
				ConfirmPassword: "Password123!",
			},
			expectError: true,
		},
		{
			name: "Missing confirm password",
			request: auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        "Password123!",
				ConfirmPassword: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Register(tt.request)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestRegistrationValidation_PasswordComplexity(t *testing.T) {
	provider := auth.NewLocalAuthProvider()

	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "Too short",
			password:    "Short1!",
			expectError: true,
		},
		{
			name:        "No uppercase",
			password:    "lowercase123!",
			expectError: true,
		},
		{
			name:        "No lowercase",
			password:    "UPPERCASE123!",
			expectError: true,
		},
		{
			name:        "No digit",
			password:    "NoDigitsHere!",
			expectError: true,
		},
		{
			name:        "No special character",
			password:    "NoSpecial123",
			expectError: true,
		},
		{
			name:        "Valid password",
			password:    "ValidPassword123!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           "test@example.com",
				Username:        "testuser",
				Password:        tt.password,
				ConfirmPassword: tt.password,
			}

			_, err := provider.Register(request)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestRegistrationValidation_UsernameFormat(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		expectError bool
	}{
		{
			name:        "Too short",
			username:    "abc",
			expectError: true,
		},
		{
			name:        "Too long",
			username:    "thisusernameistoolongandexceedsthemaximumlength",
			expectError: true,
		},
		{
			name:        "Invalid characters",
			username:    "user@name",
			expectError: true,
		},
		{
			name:        "Valid username",
			username:    "valid_user123",
			expectError: false,
		},
		{
			name:        "Valid with dots",
			username:    "valid.user.name",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := auth.NewLocalAuthProvider() // Use fresh provider for each test
			request := auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           tt.username + "@example.com", // Use unique email per test
				Username:        tt.username,
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			}

			_, err := provider.Register(request)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestRegistrationValidation_EmailFormat(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "Missing @",
			email:       "invalidemail.com",
			expectError: true,
		},
		{
			name:        "Missing domain",
			email:       "user@",
			expectError: true,
		},
		{
			name:        "Valid email",
			email:       "valid@example.com",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := auth.NewLocalAuthProvider() // Use fresh provider for each test
			request := auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     "1990-01-01",
				MobileNumber:    "+1234567890",
				Email:           tt.email,
				Username:        "testuser123", // Use valid username for email tests
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			}

			_, err := provider.Register(request)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestRegistrationValidation_DateOfBirth(t *testing.T) {
	tests := []struct {
		name        string
		dateOfBirth string
		expectError bool
	}{
		{
			name:        "Invalid format",
			dateOfBirth: "01-01-1990",
			expectError: true,
		},
		{
			name:        "Future date",
			dateOfBirth: "2050-01-01",
			expectError: true,
		},
		{
			name:        "Valid past date",
			dateOfBirth: "1990-01-01",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := auth.NewLocalAuthProvider() // Use fresh provider for each test
			request := auth.RegisterUserRequest{
				FullName:        "Test User",
				DateOfBirth:     tt.dateOfBirth,
				MobileNumber:    "+1234567890",
				Email:           "user_" + tt.name + "@example.com", // Use unique email per test
				Username:        "testuser123",                      // Use valid username for DOB tests
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
			}

			_, err := provider.Register(request)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
