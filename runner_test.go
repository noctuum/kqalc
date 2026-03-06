package main

import (
	"os/exec"
	"testing"
)

func TestPrettifyExpr(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"sqrt(2)", "√(2)"},
		{"cbrt(27)", "∛(27)"},
		{"2*pi", "2*π"},
		{"infinity", "∞"},
		{"x>=1", "x≥1"},
		{"x<=1", "x≤1"},
		{"x!=1", "x≠1"},
		{"2^64", "2⁶⁴"},
		{"10^12", "10¹²"},
		{"2^10+3^5", "2¹⁰+3⁵"},
		{"sqrt(pi)", "√(π)"},
		{"42", "42"},
	}
	for _, tt := range tests {
		if got := prettifyExpr(tt.in); got != tt.want {
			t.Errorf("prettifyExpr(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIsFiatCurrency(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"$1848.52", true},
		{"€1234.56", true},
		{"£99.99", true},
		{"¥1000", true},
		{"₹500", true},
		{"₿0.5", false},  // crypto excluded
		{"42", false},    // no currency
		{"hello", false}, // plain text
		{"$ 100", true},  // with space
	}
	for _, tt := range tests {
		if got := isFiatCurrency(tt.in); got != tt.want {
			t.Errorf("isFiatCurrency(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestFormatFiat(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		// truncate decimals to 2
		{"$1848.5234", "$ 1848.52"},
		{"€1234.5678", "€ 1234.56"},
		// already 2 decimals — no truncation, add space
		{"£99.99", "£ 99.99"},
		// comma as decimal separator
		{"€1234,5678", "€ 1234,56"},
		// already has space — don't double
		{"$ 100.00", "$ 100.00"},
		// no decimals
		{"¥1000", "¥ 1000"},
	}
	for _, tt := range tests {
		if got := formatFiat(tt.in); got != tt.want {
			t.Errorf("formatFiat(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestNewMatch(t *testing.T) {
	m := newMatch("42", "6*7", 1.0)
	if m.Id != "42" {
		t.Errorf("Id = %q, want %q", m.Id, "42")
	}
	if m.Text != "42  ←  6*7" {
		t.Errorf("Text = %q, want %q", m.Text, "42  ←  6*7")
	}
	if m.IconName != "accessories-calculator" {
		t.Errorf("IconName = %q, want %q", m.IconName, "accessories-calculator")
	}
	if m.Type != 100 {
		t.Errorf("Type = %d, want %d", m.Type, 100)
	}
	if m.Relevance != 1.0 {
		t.Errorf("Relevance = %f, want %f", m.Relevance, 1.0)
	}
}

func TestNewMatchPrettifiesExpr(t *testing.T) {
	m := newMatch("1.41421356", "sqrt(2)", 1.0)
	want := "1.41421356  ←  √(2)"
	if m.Text != want {
		t.Errorf("Text = %q, want %q", m.Text, want)
	}
}

func requireQalcRunner(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("qalc"); err != nil {
		t.Skip("qalc not installed")
	}
}

func TestMatchBasicMath(t *testing.T) {
	requireQalcRunner(t)
	r := &Runner{}
	matches, dbusErr := r.Match("qc 2+2")
	if dbusErr != nil {
		t.Fatalf("Match returned dbus error: %v", dbusErr)
	}
	if len(matches) == 0 {
		t.Fatal("Match returned no results for '2+2'")
	}
	if matches[0].Id != "4" {
		t.Errorf("Id = %q, want %q", matches[0].Id, "4")
	}
}

func TestMatchNoPrefix(t *testing.T) {
	r := &Runner{}
	matches, _ := r.Match("2+2")
	if len(matches) != 0 {
		t.Errorf("Match without prefix returned %d matches, want 0", len(matches))
	}
}

func TestMatchEmptyExpr(t *testing.T) {
	r := &Runner{}
	matches, _ := r.Match("qc ")
	if len(matches) != 0 {
		t.Errorf("Match with empty expr returned %d matches, want 0", len(matches))
	}
}

func TestMatchNoDigits(t *testing.T) {
	r := &Runner{}
	matches, _ := r.Match("qc hello world")
	if len(matches) != 0 {
		t.Errorf("Match with no digits returned %d matches, want 0", len(matches))
	}
}

func TestMatchExactAndApprox(t *testing.T) {
	requireQalcRunner(t)
	r := &Runner{}
	matches, _ := r.Match("qc 1/3")
	if len(matches) < 2 {
		t.Skipf("expected 2 matches for '1/3', got %d (may vary by qalc version)", len(matches))
	}
	if matches[1].Relevance != 0.9 {
		t.Errorf("exact match Relevance = %f, want 0.9", matches[1].Relevance)
	}
}

func TestActions(t *testing.T) {
	r := &Runner{}
	actions, dbusErr := r.Actions()
	if dbusErr != nil {
		t.Fatalf("Actions returned dbus error: %v", dbusErr)
	}
	if len(actions) != 1 {
		t.Fatalf("Actions returned %d actions, want 1", len(actions))
	}
	if actions[0].Id != "copy" {
		t.Errorf("action Id = %q, want %q", actions[0].Id, "copy")
	}
}

func TestRunReturnsNil(t *testing.T) {
	r := &Runner{}
	// CopyToClipboard will fail (no wl-copy/xclip in test) but Run handles the error
	dbusErr := r.Run("test-value", "copy")
	if dbusErr != nil {
		t.Errorf("Run returned dbus error: %v", dbusErr)
	}
}
