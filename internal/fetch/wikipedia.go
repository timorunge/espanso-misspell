// Wikipedia spell-check word list fetcher and parser.

package fetch

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/timorunge/espanso"
)

// Wikipedia downloads a Wikipedia spell-check word list and parses it into
// espanso matches. The expected format is pipe-delimited, one entry per line:
//
//	word|flags|correction
//
// Lines starting with #, =, {, *, or blank lines are skipped. Entries where
// the typo or correction contains wiki markup or where the correction is
// ambiguous (commas or parentheses) are dropped.
func Wikipedia(url, lang string) (espanso.Matches, error) {
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
		line := strings.TrimSpace(scanner.Text())
		if line == "" || isSkipPrefix(line[0]) {
			continue
		}

		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue
		}

		typo := strings.TrimSpace(parts[0])
		correction := strings.TrimSpace(parts[2])
		if typo == "" || correction == "" {
			continue
		}
		if strings.ContainsAny(typo, "<>[]{}*") {
			continue
		}
		if strings.ContainsAny(correction, ",()<>[]{}") {
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

func isSkipPrefix(b byte) bool {
	return b == '#' || b == '=' || b == '{' || b == '*'
}
