package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	db := Connect()
	if db == nil {
		t.Fatal("No se pudo conectar a la base de datos")
	}
}
