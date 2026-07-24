package integration

import (
    "strings"
    "testing"

    "artha-kosha/apps/finance-api/internal/auth"
)

func TestPasswordHashing_Argon2idFormatAndVerification(t *testing.T) {
    provider := auth.NewLocalAuthProvider()

    req := auth.RegisterUserRequest{
        FullName:        "Security Test",
        DateOfBirth:     "1990-01-01",
        MobileNumber:    "+10000000000",
        Email:           "sec@example.com",
        Username:        "securityuser",
        Password:        "S3cureP@ss!",
        ConfirmPassword: "S3cureP@ss!",
    }

    resp, err := provider.Register(req)
    if err != nil {
        t.Fatalf("register failed: %v", err)
    }

    if !strings.Contains(resp.PasswordHash, "argon2id") {
        t.Fatalf("expected argon2id hash format, got: %s", resp.PasswordHash)
    }

    // Verify login works with plain password
    loginReq := auth.LoginRequest{Username: req.Username, Password: req.Password}
    _, err = provider.Login(loginReq)
    if err != nil {
        t.Fatalf("login failed: %v", err)
    }
}
