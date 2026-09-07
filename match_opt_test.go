package main

import (
	"sync"
	"testing"
	"time"
)

// stubEvaluate replaces the qalc seam and returns a restore func.
func stubEvaluate(t *testing.T, fn func(expr, mode string) (string, error)) {
	t.Helper()
	orig := evaluate
	evaluate = fn
	t.Cleanup(func() { evaluate = orig })
	resetCache()
	t.Cleanup(resetCache)
}

func resetCache() {
	cacheMu.Lock()
	cache = map[string]cacheEntry{}
	cacheMu.Unlock()
}

// recorder counts evaluate calls per mode, safely under concurrency.
type recorder struct {
	mu    sync.Mutex
	modes []string
}

func (rec *recorder) record(mode string) {
	rec.mu.Lock()
	rec.modes = append(rec.modes, mode)
	rec.mu.Unlock()
}

func (rec *recorder) count(mode string) int {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	n := 0
	for _, m := range rec.modes {
		if m == mode {
			n++
		}
	}
	return n
}

func TestMatchSkipsExactForUnitConversion(t *testing.T) {
	rec := &recorder{}
	stubEvaluate(t, func(expr, mode string) (string, error) {
		rec.record(mode)
		return "62.137119 mi", nil
	})

	r := &Runner{}
	if _, err := r.Match("qc 100 km to miles"); err != nil {
		t.Fatalf("Match error: %v", err)
	}

	if got := rec.count("exact"); got != 0 {
		t.Errorf("exact evaluations = %d, want 0 for a ' to ' conversion", got)
	}
	if got := rec.count("approximate"); got != 1 {
		t.Errorf("approximate evaluations = %d, want 1", got)
	}
}

func TestMatchEvaluatesExactForPlainExpression(t *testing.T) {
	rec := &recorder{}
	stubEvaluate(t, func(expr, mode string) (string, error) {
		rec.record(mode)
		if mode == "exact" {
			return "1/3", nil
		}
		return "0.333333333", nil
	})

	r := &Runner{}
	matches, err := r.Match("qc 2/6")
	if err != nil {
		t.Fatalf("Match error: %v", err)
	}

	if got := rec.count("exact"); got != 1 {
		t.Errorf("exact evaluations = %d, want 1", got)
	}
	if len(matches) != 2 {
		t.Fatalf("matches = %d, want 2 (approximate + exact)", len(matches))
	}
	if matches[1].Id != "1/3" {
		t.Errorf("exact match Id = %q, want %q", matches[1].Id, "1/3")
	}
}

func TestMatchRunsBothModesConcurrently(t *testing.T) {
	const delay = 60 * time.Millisecond
	stubEvaluate(t, func(expr, mode string) (string, error) {
		time.Sleep(delay)
		if mode == "exact" {
			return "1/3", nil
		}
		return "0.333333333", nil
	})

	r := &Runner{}
	start := time.Now()
	if _, err := r.Match("qc 2/6"); err != nil {
		t.Fatalf("Match error: %v", err)
	}
	elapsed := time.Since(start)

	// Sequential would take 2*delay. Allow generous headroom for scheduling.
	if elapsed >= 2*delay {
		t.Errorf("Match took %v, want well under %v — evaluations are not concurrent", elapsed, 2*delay)
	}
}

func TestMatchCachesRepeatedQuery(t *testing.T) {
	rec := &recorder{}
	stubEvaluate(t, func(expr, mode string) (string, error) {
		rec.record(mode)
		return "4", nil
	})

	r := &Runner{}
	for i := 0; i < 3; i++ {
		if _, err := r.Match("qc 2+2"); err != nil {
			t.Fatalf("Match error: %v", err)
		}
	}

	if got := rec.count("approximate"); got != 1 {
		t.Errorf("approximate evaluations = %d, want 1 — repeated query was not cached", got)
	}
}

func TestMatchCacheDistinguishesExpressions(t *testing.T) {
	rec := &recorder{}
	stubEvaluate(t, func(expr, mode string) (string, error) {
		rec.record(mode)
		return expr, nil
	})

	r := &Runner{}
	if _, err := r.Match("qc 2+2"); err != nil {
		t.Fatalf("Match error: %v", err)
	}
	if _, err := r.Match("qc 3+3"); err != nil {
		t.Fatalf("Match error: %v", err)
	}

	if got := rec.count("approximate"); got != 2 {
		t.Errorf("approximate evaluations = %d, want 2 — distinct expressions must not share a cache entry", got)
	}
}

func TestCacheEvictsWhenFull(t *testing.T) {
	resetCache()
	t.Cleanup(resetCache)

	for i := 0; i < cacheMaxEntries+10; i++ {
		cacheStore(string(rune(i))+"k", "v", nil)
	}

	cacheMu.Lock()
	size := len(cache)
	cacheMu.Unlock()

	if size > cacheMaxEntries {
		t.Errorf("cache size = %d, want <= %d — cache grows without bound", size, cacheMaxEntries)
	}
}

func TestCacheExpiresStaleEntries(t *testing.T) {
	resetCache()
	t.Cleanup(resetCache)

	cacheMu.Lock()
	cache["stale"] = cacheEntry{result: "old", at: time.Now().Add(-cacheTTL - time.Second)}
	cacheMu.Unlock()

	if _, _, ok := cacheLookup("stale"); ok {
		t.Error("cacheLookup returned an entry older than cacheTTL, want miss")
	}
}
