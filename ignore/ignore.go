package ignore

import (
	"bufio"
	"errors"
	"io"
	"log"
	"path/filepath"
	"strings"

	bogieio "github.com/sethpollack/bogie/io"
)

// BogieIgnore default name of an ignorefile.
const BogieIgnore = ".bogieignore"

// Rules is a collection of path matching rules.
//
// Parse() and ParseFile() will construct and populate new Rules.
// Empty() will create an immutable empty ruleset.
type Rules struct {
	patterns []*pattern
}

func (r *Rules) Clone() *Rules {
	return &Rules{patterns: r.patterns}
}

// Empty builds an inital ruleset.
func Init() *Rules {
	r := Rules{patterns: []*pattern{}}
	r.parseRule(".bogieignore")
	return &r
}

// ParseFile parses a BogieIgnore file.
func (r *Rules) ParseFile(filepath string) error {
	rules, err := bogieio.ReadInput(filepath)
	if err != nil {
		return err
	}
	return r.Parse(strings.NewReader(rules))
}

// Parse parses a rules file
func (r *Rules) Parse(file io.Reader) error {
	s := bufio.NewScanner(file)
	for s.Scan() {
		if err := r.parseRule(s.Text()); err != nil {
			return err
		}
	}
	return s.Err()
}

// Len returns the number of patterns in this rule set.
func (r *Rules) Len() int {
	return len(r.patterns)
}

// Ignore evalutes the file at the given path, and returns true if it should be ignored.
//
// Ignore evaluates path against the rules in order. Evaluation stops when a match
// is found. Matching a negative rule will stop evaluation.
func (r *Rules) Ignore(path string, isDir bool) bool {
	// Disallow ignoring the current working directory.
	// See issue:
	// 1776 (New York City) Hamilton: "Pardon me, are you Aaron Burr, sir?"
	if path == "." || path == "./" {
		return false
	}
	for _, p := range r.patterns {
		if p.match == nil {
			log.Printf("ignore: no matcher supplied for %q", p.raw)
			return false
		}

		// For negative rules, we need to capture and return non-matches,
		// and continue for matches.
		if p.negate {
			if p.mustDir && !isDir {
				return true
			}
			if !p.match(path) {
				return true
			}
			continue
		}

		// If the rule is looking for directories, and this is not a directory,
		// skip it.
		if p.mustDir && !isDir {
			continue
		}
		if p.match(path) {
			return true
		}
	}
	return false
}

// parseRule parses a rule string and creates a pattern, which is then stored in the Rules object.
func (r *Rules) parseRule(rule string) error {
	rule = strings.TrimSpace(rule)

	// Ignore blank lines
	if rule == "" {
		return nil
	}
	// Comment
	if strings.HasPrefix(rule, "#") {
		return nil
	}

	// Fail any rules that contain **
	if strings.Contains(rule, "**") {
		return errors.New("double-star (**) syntax is not supported")
	}

	// Fail any patterns that can't compile. A non-empty string must be
	// given to Match() to avoid optimization that skips rule evaluation.
	if _, err := filepath.Match(rule, "abc"); err != nil {
		return err
	}

	p := &pattern{raw: rule}

	// Negation is handled at a higher level, so strip the leading ! from the
	// string.
	if strings.HasPrefix(rule, "!") {
		p.negate = true
		rule = rule[1:]
	}

	// Directory verification is handled by a higher level, so the trailing /
	// is removed from the rule. That way, a directory named "foo" matches,
	// even if the supplied string does not contain a literal slash character.
	if strings.HasSuffix(rule, "/") {
		p.mustDir = true
		rule = strings.TrimSuffix(rule, "/")
	}

	if strings.HasPrefix(rule, "/") {
		// Require path matches the root path.
		p.match = func(n string) bool {
			rule = strings.TrimPrefix(rule, "/")
			ok, err := filepath.Match(rule, n)
			if err != nil {
				log.Printf("Failed to compile %q: %s", rule, err)
				return false
			}
			return ok
		}
	} else if strings.Contains(rule, "/") {
		// require structural match.
		p.match = func(n string) bool {
			ok, err := filepath.Match(rule, n)
			if err != nil {
				log.Printf("Failed to compile %q: %s", rule, err)
				return false
			}
			return ok
		}
	} else {
		p.match = func(n string) bool {
			// When there is no slash in the pattern, we evaluate ONLY the
			// filename.
			n = filepath.Base(n)
			ok, err := filepath.Match(rule, n)
			if err != nil {
				log.Printf("Failed to compile %q: %s", rule, err)
				return false
			}
			return ok
		}
	}

	r.patterns = append(r.patterns, p)
	return nil
}

// matcher is a function capable of computing a match.
//
// It returns true if the rule matches.
type matcher func(name string) bool

// pattern describes a pattern to be matched in a rule set.
type pattern struct {
	// raw is the unparsed string, with nothing stripped.
	raw string
	// match is the matcher function.
	match matcher
	// negate indicates that the rule's outcome should be negated.
	negate bool
	// mustDir indicates that the matched file must be a directory.
	mustDir bool
}
