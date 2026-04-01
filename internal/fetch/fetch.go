// Package fetch downloads spell-check word lists and parses them into
// espanso matches.
package fetch

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/timorunge/espanso"
)

// LineParser extracts a (typo, correction) pair from a single line of a
// word list. It returns ok=false when the line should be skipped.
type LineParser func(line string) (typo, correction string, ok bool)

// Fetch downloads a word list from url and parses it into espanso matches
// using the given line parser. The lang parameter is used for User-Agent
// and error messages.
func Fetch(ctx context.Context, url, lang string, parse LineParser) (espanso.Matches, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	ua := fmt.Sprintf("espanso-misspell/%s (https://github.com/timorunge/espanso-misspell)", lang)
	req.Header.Set("User-Agent", ua)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s word list: %w", lang, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("fetch %s word list: status %d", lang, resp.StatusCode)
	}

	return ParseLines(resp.Body, parse)
}

// ParseLines reads lines from r and parses each with the given line parser.
func ParseLines(r io.Reader, parse LineParser) (espanso.Matches, error) {
	var matches espanso.Matches
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		typo, correction, ok := parse(scanner.Text())
		if !ok {
			continue
		}
		matches = append(matches, espanso.Match{
			Triggers: []string{typo},
			Replace:  correction,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read word list: %w", err)
	}

	return matches, nil
}

// PipeLine parses pipe-delimited lines of the form word|flags|correction.
// Lines starting with #, =, {, *, or blank lines are skipped. Entries
// containing wiki markup or ambiguous corrections are dropped.
func PipeLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || isSkipPrefix(line[0]) {
		return "", "", false
	}

	typoRaw, rest, ok := strings.Cut(line, "|")
	if !ok {
		return "", "", false
	}
	_, corrRaw, ok := strings.Cut(rest, "|")
	if !ok {
		return "", "", false
	}

	typo := strings.TrimSpace(typoRaw)
	correction := strings.TrimSpace(corrRaw)
	if typo == "" || correction == "" {
		return "", "", false
	}
	if strings.ContainsAny(typo, "<>[]{}*") {
		return "", "", false
	}
	if strings.ContainsAny(correction, ",()<>[]{}") {
		return "", "", false
	}

	return typo, correction, true
}

// TemplateLine returns a line parser that extracts pairs from wiki templates
// like {{Suggestion|typo|correction}} or {{BR1|typo|correction}}. Only simple
// literal entries are kept; regex patterns, wildcards, and wiki markup are
// skipped.
func TemplateLine(tmpl string) LineParser {
	prefix := "{{" + tmpl + "|"

	return func(line string) (string, string, bool) {
		_, after, found := strings.Cut(line, prefix)
		if !found {
			return "", "", false
		}

		inner, _, found := strings.Cut(after, "}}")
		if !found {
			return "", "", false
		}

		// Extract positional arguments, skipping named params (key=value).
		var args []string
		for _, part := range strings.Split(inner, "|") {
			part = strings.TrimSpace(part)
			if part == "" || strings.Contains(part, "=") {
				continue
			}
			args = append(args, part)
		}
		if len(args) < 2 {
			return "", "", false
		}

		typo := args[0]
		correction := args[1]

		if strings.ContainsAny(typo, "<>[]{}*$()\\\"") {
			return "", "", false
		}
		if strings.ContainsAny(correction, "<>[]{}*$()\\,\"") {
			return "", "", false
		}

		return typo, correction, true
	}
}

// CodespellLine parses codespell dictionary lines of the form
// misspelling->correction. Ambiguous entries with multiple comma-separated
// corrections are skipped.
func CodespellLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", false
	}

	typo, correction, ok := strings.Cut(line, "->")
	if !ok {
		return "", "", false
	}

	typo = strings.TrimSpace(typo)
	correction = strings.TrimSpace(correction)
	if typo == "" || correction == "" {
		return "", "", false
	}

	// Skip ambiguous entries (multiple corrections).
	if strings.Contains(correction, ",") {
		return "", "", false
	}

	return typo, correction, true
}

func isSkipPrefix(b byte) bool {
	return b == '#' || b == '=' || b == '{' || b == '*'
}
