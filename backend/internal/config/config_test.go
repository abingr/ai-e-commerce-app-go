package config

import "testing"

func TestGetCSVEnvTrimsAndSkipsEmptyValues(t *testing.T) {
	t.Setenv("ECOMMERCE_TEST_CSV", " http://localhost:5173, https://app.example.com, ")

	values := getCSVEnv("ECOMMERCE_TEST_CSV", "")

	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}

	if values[0] != "http://localhost:5173" {
		t.Fatalf("expected first origin to be trimmed, got %q", values[0])
	}

	if values[1] != "https://app.example.com" {
		t.Fatalf("expected second origin, got %q", values[1])
	}
}
