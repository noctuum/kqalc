package main

import (
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

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

// idleDelay is how long kqalc stays up after KRunner ends a match session.
// Past cacheTTL every cached result has expired, so exiting loses nothing.
var idleDelay = cacheTTL

// Runner implements the org.kde.krunner1 DBus interface.
type Runner struct {
	// canExit reports whether the bus can start kqalc again on the next query.
	canExit func() bool
	// done is closed when kqalc has been idle for idleDelay and may exit.
	done chan struct{}
	once sync.Once

	mu sync.Mutex
	// activity counts Match and Teardown calls. A timer armed by Teardown
	// acts only if no call came in after it.
	activity int
}

func newRunner(canExit func() bool) *Runner {
	return &Runner{canExit: canExit, done: make(chan struct{})}
}

// Teardown is called by KRunner when a match session is over.
func (r *Runner) Teardown() *dbus.Error {
	seen := r.touch()
	time.AfterFunc(idleDelay, func() { r.expire(seen) })
	return nil
}

func (r *Runner) touch() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.activity++
	return r.activity
}

func (r *Runner) expire(seen int) {
	r.mu.Lock()
	stale := seen != r.activity
	r.mu.Unlock()
	if stale || !r.canExit() {
		return
	}
	r.once.Do(func() { close(r.done) })
}

func (r *Runner) Match(query string) ([]remoteMatch, *dbus.Error) {
	r.touch()
	if !strings.HasPrefix(query, triggerPrefix) {
		return nil, nil
	}
	expr := strings.TrimPrefix(query, triggerPrefix)
	expr = strings.TrimSpace(expr)
	if expr == "" || !strings.ContainsAny(expr, "0123456789+-*/^=().") {
		return nil, nil
	}

	// A conversion never yields a usable exact form, so it gets one qalc run.
	// Anything else runs both modes at once: each run costs about 58 ms, almost
	// all of it qalc startup.
	var (
		approx, exact string
		errA, errE    error
	)
	if strings.Contains(expr, " to ") {
		approx, errA = evaluateCached(expr, "approximate")
	} else {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); approx, errA = evaluateCached(expr, "approximate") }()
		go func() { defer wg.Done(); exact, errE = evaluateCached(expr, "exact") }()
		wg.Wait()
	}

	var matches []remoteMatch

	if errA == nil && approx != "" && approx != expr {
		if isFiatCurrency(approx) {
			approx = formatFiat(approx)
		}
		matches = append(matches, newMatch(approx, expr, 1.0))
	}
	prettyExpr := prettifyExpr(expr)
	if errE == nil && exact != "" && exact != approx && exact != expr && exact != prettyExpr &&
		!isFiatCurrency(approx) && !strings.ContainsRune(approx, 'e') {
		matches = append(matches, newMatch(exact, expr, 0.9))
	}

	if len(matches) == 0 && errA != nil {
		log.Printf("qalc error for %q: %v", expr, errA)
	}

	return matches, nil
}

// isFiatCurrency reports whether the result holds a fiat currency symbol.
// fiatSymbols leaves out ₿, so a Bitcoin result keeps its exact form.
func isFiatCurrency(result string) bool {
	for _, r := range result {
		if strings.ContainsRune(fiatSymbols, r) {
			return true
		}
	}
	return false
}

var decimalRe = regexp.MustCompile(`(\d)([.,])(\d{2})\d+`)

const fiatSymbols = "$€£¥₹₽₺₴₾₫₦₱₡₲₵₼₸₪₩₭₮₯₷₨"

// formatFiat truncates the decimals to two and puts a space after the
// currency symbol.
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

var caretRe = regexp.MustCompile(`\^(\d+)`)

// prettifyExpr rewrites an expression for display: sqrt(2) → √(2),
// 2^64 → 2⁶⁴.
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
