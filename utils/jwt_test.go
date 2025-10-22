package utils

import (
	"os"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	token, err := GenerateToken("test@example.com", "user")
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}
	if token == "" {
		t.Error("Token generado vacío")
	}
}

func TestValidateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	token, err := GenerateToken("test@example.com", "admin")
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}
	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("Error validando token: %v", err)
	}
	if claims.Email != "test@example.com" || claims.Role != "admin" {
		t.Errorf("Claims incorrectos: %+v", claims)
	}
}
