package auth

import (
	"testing"
)

func TestPasswordService_ChangePassword(t *testing.T) {
	svc := NewLocalAuthProvider()
	
	err := svc.ChangePassword(ChangePasswordRequest{})
	if err == nil {
		t.Error("expected error")
	}
}
