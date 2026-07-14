package auth

import "testing"

func TestLocalAuthProviderExposesAuthAbstraction(t *testing.T) {
	provider := NewLocalAuthProvider()
	if provider == nil {
		t.Fatal("expected a local auth provider implementation")
	}
}
