package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/godbus/dbus/v5"
)

const triggerPrefix = "qc "

type remoteMatch struct {
	Id         string
	Text       string
	IconName   string
	Type       int32
	Relevance  float64
	Properties map[string]dbus.Variant
}

type remoteAction struct {
	Id       string
	Text     string
	IconName string
}

// Runner implements the org.kde.krunner1 DBus interface.
type Runner struct{}

func (r *Runner) Match(query string) ([]remoteMatch, *dbus.Error) {
	if !strings.HasPrefix(query, triggerPrefix) {
		return nil, nil
	}
	expr := strings.TrimPrefix(query, triggerPrefix)
	expr = strings.TrimSpace(expr)
	if expr == "" || !strings.ContainsAny(expr, "0123456789+-*/^=().") {
		return nil, nil
	}

	approx, errA := Evaluate(expr, "approximate")
	exact, errE := Evaluate(expr, "exact")

	var matches []remoteMatch

	if errA == nil && approx != "" && approx != expr {
		if isFiatCurrency(approx) {
			approx = formatFiat(approx)
		}
		matches = append(matches, newMatch(approx, expr, 1.0))
	}
	prettyExpr := prettifyExpr(expr)
	if errE == nil && exact != "" && exact != approx && exact != expr && exact != prettyExpr &&
		!isFiatCurrency(approx) && !strings.ContainsRune(approx, 'e') &&
		!strings.Contains(expr, " to ") {
		matches = append(matches, newMatch(exact, expr, 0.9))
	}

	if len(matches) == 0 && errA != nil {
		log.Printf("qalc error for %q: %v", expr, errA)
	}

	return matches, nil
}

// isFiatCurrency checks if the result contains a fiat currency symbol.
// Crypto symbols like ₿ (Bitcoin) are excluded — exact results are useful there.
func isFiatCurrency(result string) bool {
	for _, r := range result {
		if strings.ContainsRune(fiatSymbols, r) {
			return true
		}
	}
	return false
}

// formatFiat truncates decimals to 2 and inserts a space after the currency symbol.
var decimalRe = regexp.MustCompile(`(\d)([.,])(\d{2})\d+`)

const fiatSymbols = "$€£¥₹₽₺₴₾₫₦₱₡₲₵₼₸₪₩₭₮₯₷₨"

func formatFiat(result string) string {
	result = decimalRe.ReplaceAllString(result, "${1}${2}${3}")
	for i, r := range result {
		if strings.ContainsRune(fiatSymbols, r) {
			pos := i + len(string(r))
			if pos < len(result) && result[pos] != ' ' {
				result = result[:pos] + " " + result[pos:]
			}
			break
		}
	}
	return result
}

// prettifyExpr replaces ASCII function names with Unicode equivalents for display.
var exprReplacements = strings.NewReplacer(
	"sqrt(", "√(",
	"cbrt(", "∛(",
	"pi", "π",
	"infinity", "∞",
	">=", "≥",
	"<=", "≤",
	"!=", "≠",
)

var superscriptDigits = strings.NewReplacer(
	"0", "⁰", "1", "¹", "2", "²", "3", "³", "4", "⁴",
	"5", "⁵", "6", "⁶", "7", "⁷", "8", "⁸", "9", "⁹",
)

// prettifyExponent converts "2^64" → "2⁶⁴".
var caretRe = regexp.MustCompile(`\^(\d+)`)

func prettifyExpr(expr string) string {
	expr = exprReplacements.Replace(expr)
	expr = caretRe.ReplaceAllStringFunc(expr, func(m string) string {
		return superscriptDigits.Replace(m[1:]) // skip ^
	})
	return expr
}

func newMatch(result, expr string, relevance float64) remoteMatch {
	return remoteMatch{
		Id:         result,
		Text:       result + "  ←  " + prettifyExpr(expr),
		IconName:   "accessories-calculator",
		Type:       100, // ExactMatch
		Relevance:  relevance,
		Properties: map[string]dbus.Variant{},
	}
}

func (r *Runner) Actions() ([]remoteAction, *dbus.Error) {
	return []remoteAction{
		{"copy", "Copy to clipboard", "edit-copy"},
	}, nil
}

func (r *Runner) Run(matchId string, actionId string) *dbus.Error {
	if err := CopyToClipboard(matchId); err != nil {
		log.Printf("clipboard error: %v", err)
	}
	return nil
}
