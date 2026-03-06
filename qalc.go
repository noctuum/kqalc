package main

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

const evalTimeout = 3 * time.Second

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
