package models

import (
	"testing"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error al hashear password: %v", err)
	}
	if !CheckPassword(password, hash) {
		t.Error("La verificación de password falló")
	}
	if CheckPassword("wrongpassword", hash) {
		t.Error("La verificación debería fallar con password incorrecto")
	}
}
