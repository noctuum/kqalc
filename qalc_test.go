package main

import (
	"os/exec"
	"strings"
	"testing"
)

func requireQalc(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("qalc"); err != nil {
		t.Skip("qalc not installed")
	}
}

func TestEvaluateApproximate(t *testing.T) {
	requireQalc(t)
	tests := []struct {
		expr, want string
	}{
		{"2+2", "4"},
		{"10/3", "3.33333333"},
		{"sqrt(2)", "1.41421356"},
	}
	for _, tt := range tests {
		got, err := Evaluate(tt.expr, "approximate")
		if err != nil {
			t.Errorf("Evaluate(%q, approximate) error: %v", tt.expr, err)
			continue
		}
		// Normalize locale decimal separator (comma → dot) for comparison
		norm := strings.ReplaceAll(got, ",", ".")
		if norm != tt.want {
			t.Errorf("Evaluate(%q, approximate) = %q, want %q", tt.expr, got, tt.want)
		}
	}
}

func TestEvaluateExact(t *testing.T) {
	requireQalc(t)
	tests := []struct {
		expr, want string
	}{
		{"1/3", "1/3"},
		{"sqrt(2)", "√(2)"},
	}
	for _, tt := range tests {
		got, err := Evaluate(tt.expr, "exact")
		if err != nil {
			t.Errorf("Evaluate(%q, exact) error: %v", tt.expr, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Evaluate(%q, exact) = %q, want %q", tt.expr, got, tt.want)
		}
	}
}

func TestEvaluateInvalidExpr(t *testing.T) {
	requireQalc(t)
	// qalc handles most inputs without error, but we verify it doesn't panic
	_, _ = Evaluate("", "approximate")
	_, _ = Evaluate("???", "exact")
}
