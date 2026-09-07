package main

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const evalTimeout = 3 * time.Second

// evaluate is the qalc seam; tests replace it to avoid spawning processes.
var evaluate = Evaluate

const (
	cacheTTL        = 60 * time.Second
	cacheMaxEntries = 256
)

type cacheEntry struct {
	result string
	err    error
	at     time.Time
}

var (
	cacheMu sync.Mutex
	cache   = map[string]cacheEntry{}
)

func cacheLookup(key string) (string, error, bool) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	e, ok := cache[key]
	if !ok || time.Since(e.at) > cacheTTL {
		return "", nil, false
	}
	return e.result, e.err, true
}

// cacheStore drops the whole map once full: entries live for cacheTTL at most,
// so LRU bookkeeping would cost more than the occasional re-evaluation.
func cacheStore(key, result string, err error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if len(cache) >= cacheMaxEntries {
		cache = make(map[string]cacheEntry, cacheMaxEntries)
	}
	cache[key] = cacheEntry{result: result, err: err, at: time.Now()}
}

// evaluateCached memoises qalc results — KRunner re-queries on every keystroke,
// so typing "2+2" would otherwise evaluate "2", "2+" and "2+2".
func evaluateCached(expr, mode string) (string, error) {
	key := mode + "\x00" + expr
	if result, err, ok := cacheLookup(key); ok {
		return result, err
	}
	result, err := evaluate(expr, mode)
	cacheStore(key, result, err)
	return result, err
}

// Evaluate runs qalc with the given expression and returns the result.
// Mode must be "approximate" or "exact".
func Evaluate(expression, mode string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), evalTimeout)
	defer cancel()

	args := []string{"-t", "+u8",
		"-set", "assumptions number",
		"-set", "lowercase e on",
		"-set", "unicode on",
	}

	switch mode {
	case "approximate":
		args = append(args,
			"-set", "precision 12",
			"-set", "max decimals 8",
			"-set", "exact off",
			"-set", "approximation approximate",
			"-set", "fractions off",
		)
	case "exact":
		args = append(args,
			"-set", "precision 100",
			"-set", "max decimals 26",
			"-set", "exact on",
			"-set", "approximation exact",
			"-set", "fractions on",
		)
	}

	args = append(args, expression)

	out, err := exec.CommandContext(ctx, "qalc", args...).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
