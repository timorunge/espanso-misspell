// Tests for line parsers.

package fetch

import "testing"

func TestPipeLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		line           string
		wantTypo       string
		wantCorrection string
		wantOK         bool
	}{
		{
			name:           "valid entry",
			line:           "Aaachen||Aachen",
			wantTypo:       "Aaachen",
			wantCorrection: "Aachen",
			wantOK:         true,
		},
		{
			name:           "valid with flags",
			line:           "Abbild||Bild",
			wantTypo:       "Abbild",
			wantCorrection: "Bild",
			wantOK:         true,
		},
		{
			name:   "empty line",
			line:   "",
			wantOK: false,
		},
		{
			name:   "comment line",
			line:   "# this is a comment",
			wantOK: false,
		},
		{
			name:   "heading line",
			line:   "== A ==",
			wantOK: false,
		},
		{
			name:   "template line",
			line:   "{{some template}}",
			wantOK: false,
		},
		{
			name:   "list marker",
			line:   "* list item",
			wantOK: false,
		},
		{
			name:   "too few pipes",
			line:   "word|correction",
			wantOK: false,
		},
		{
			name:   "wiki markup in typo",
			line:   "[[word]]||correction",
			wantOK: false,
		},
		{
			name:   "ambiguous correction with comma",
			line:   "word||correction1, correction2",
			wantOK: false,
		},
		{
			name:   "parentheses in correction",
			line:   "word||correction (variant)",
			wantOK: false,
		},
		{
			name:   "empty typo",
			line:   "||correction",
			wantOK: false,
		},
		{
			name:   "empty correction",
			line:   "word||",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			typo, correction, ok := PipeLine(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("PipeLine(%q) ok = %v, want %v", tt.line, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if typo != tt.wantTypo {
				t.Errorf("typo = %q, want %q", typo, tt.wantTypo)
			}
			if correction != tt.wantCorrection {
				t.Errorf("correction = %q, want %q", correction, tt.wantCorrection)
			}
		})
	}
}

func TestTemplateLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		template       string
		line           string
		wantTypo       string
		wantCorrection string
		wantOK         bool
	}{
		{
			name:           "valid Suggestion template",
			template:       "Suggestion",
			line:           "{{Suggestion|aigu|aigu}}",
			wantTypo:       "aigu",
			wantCorrection: "aigu",
			wantOK:         true,
		},
		{
			name:           "valid BR1 template",
			template:       "BR1",
			line:           "{{BR1|conbattere|combattere}}",
			wantTypo:       "conbattere",
			wantCorrection: "combattere",
			wantOK:         true,
		},
		{
			name:           "template with extra text around",
			template:       "Suggestion",
			line:           "some text {{Suggestion|foo|bar}} more text",
			wantTypo:       "foo",
			wantCorrection: "bar",
			wantOK:         true,
		},
		{
			name:     "skips named params",
			template: "Suggestion",
			line:     "{{Suggestion|key=value|foo|bar}}",
			wantTypo: "foo", wantCorrection: "bar", wantOK: true,
		},
		{
			name:     "wrong template name",
			template: "Suggestion",
			line:     "{{BR1|foo|bar}}",
			wantOK:   false,
		},
		{
			name:     "no closing braces",
			template: "Suggestion",
			line:     "{{Suggestion|foo|bar",
			wantOK:   false,
		},
		{
			name:     "too few args",
			template: "Suggestion",
			line:     "{{Suggestion|onlyone}}",
			wantOK:   false,
		},
		{
			name:     "regex in typo",
			template: "BR1",
			line:     "{{BR1|foo\\d+|bar}}",
			wantOK:   false,
		},
		{
			name:     "wildcard in typo",
			template: "Suggestion",
			line:     "{{Suggestion|foo*|bar}}",
			wantOK:   false,
		},
		{
			name:     "wiki markup in correction",
			template: "Suggestion",
			line:     "{{Suggestion|foo|[[bar]]}}",
			wantOK:   false,
		},
		{
			name:     "comma in correction",
			template: "BR1",
			line:     "{{BR1|foo|bar, baz}}",
			wantOK:   false,
		},
		{
			name:     "no template at all",
			template: "Suggestion",
			line:     "just plain text",
			wantOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			parse := TemplateLine(tt.template)
			typo, correction, ok := parse(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("TemplateLine(%q)(%q) ok = %v, want %v", tt.template, tt.line, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if typo != tt.wantTypo {
				t.Errorf("typo = %q, want %q", typo, tt.wantTypo)
			}
			if correction != tt.wantCorrection {
				t.Errorf("correction = %q, want %q", correction, tt.wantCorrection)
			}
		})
	}
}

func TestCodespellLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		line           string
		wantTypo       string
		wantCorrection string
		wantOK         bool
	}{
		{
			name:           "valid entry",
			line:           "acessibility->accessibility",
			wantTypo:       "acessibility",
			wantCorrection: "accessibility",
			wantOK:         true,
		},
		{
			name:           "with whitespace",
			line:           "  recieve -> receive  ",
			wantTypo:       "recieve",
			wantCorrection: "receive",
			wantOK:         true,
		},
		{
			name:   "empty line",
			line:   "",
			wantOK: false,
		},
		{
			name:   "no arrow separator",
			line:   "just a word",
			wantOK: false,
		},
		{
			name:   "ambiguous with comma",
			line:   "typo->fix1, fix2",
			wantOK: false,
		},
		{
			name:   "empty typo",
			line:   "->correction",
			wantOK: false,
		},
		{
			name:   "empty correction",
			line:   "typo->",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			typo, correction, ok := CodespellLine(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("CodespellLine(%q) ok = %v, want %v", tt.line, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if typo != tt.wantTypo {
				t.Errorf("typo = %q, want %q", typo, tt.wantTypo)
			}
			if correction != tt.wantCorrection {
				t.Errorf("correction = %q, want %q", correction, tt.wantCorrection)
			}
		})
	}
}
