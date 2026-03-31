// Wikipedia spell-check word list fetcher and parser.

package fetch

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/timorunge/espanso"
)

// LineParser extracts a (typo, correction) pair from a single line of raw
// wiki text. It returns ok=false when the line should be skipped.
type LineParser func(line string) (typo, correction string, ok bool)

// Wikipedia downloads a Wikipedia spell-check word list and parses it into
// espanso matches using the given line parser.
func Wikipedia(url, lang string, parse LineParser) (espanso.Matches, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	ua := fmt.Sprintf("espanso-misspell/%s (https://github.com/timorunge/espanso-misspell)", lang)
	req.Header.Set("User-Agent", ua)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s word list: %w", lang, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s word list: status %d", lang, resp.StatusCode)
	}

	var matches espanso.Matches
	scanner := bufio.NewScanner(resp.Body)
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
		return nil, fmt.Errorf("read %s word list: %w", lang, err)
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

	parts := strings.SplitN(line, "|", 3)
	if len(parts) != 3 {
		return "", "", false
	}

	typo := strings.TrimSpace(parts[0])
	correction := strings.TrimSpace(parts[2])
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
func TemplateLine(template string) LineParser {
	prefix := "{{" + template + "|"

	return func(line string) (string, string, bool) {
		start := strings.Index(line, prefix)
		if start < 0 {
			return "", "", false
		}
		rest := line[start+len(prefix):]

		end := strings.Index(rest, "}}")
		if end < 0 {
			return "", "", false
		}

		// Extract positional arguments, skipping named params (key=value).
		var args []string
		for _, part := range strings.Split(rest[:end], "|") {
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

func isSkipPrefix(b byte) bool {
	return b == '#' || b == '=' || b == '{' || b == '*'
}
