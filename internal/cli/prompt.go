package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/initrepo"
)

// stdinReader and interactive are variables so tests can script the terminal.
var (
	stdinReader io.Reader = os.Stdin
	// interactive reports whether both standard input and standard output are
	// terminals: the only case in which the CLI asks anything.
	interactive = func(stdout io.Writer) bool { return isCharDevice(os.Stdin) && isCharDevice(stdout) }
)

func isCharDevice(v any) bool {
	f, ok := v.(*os.File)
	if !ok {
		return false
	}
	st, err := f.Stat()
	return err == nil && st.Mode()&os.ModeCharDevice != 0
}

// prompter asks plain questions on out and reads the answers from one buffered
// reader, so a script of answers is read in order. End of input takes the default.
type prompter struct {
	in  *bufio.Reader
	out io.Writer
}

func newPrompter(r io.Reader, w io.Writer) *prompter { return &prompter{bufio.NewReader(r), w} }

// line prints q and returns the trimmed answer; ok is false at end of input.
func (p *prompter) line(q string) (string, bool) {
	fmt.Fprint(p.out, q)
	s, err := p.in.ReadString('\n')
	if err != nil && s == "" {
		fmt.Fprintln(p.out)
		return "", false
	}
	return strings.TrimSpace(s), true
}

// yesNo asks "q [y/N] " (the capital letter is the default).
func (p *prompter) yesNo(q string, def bool) bool {
	hint := " [y/N] "
	if def {
		hint = " [Y/n] "
	}
	for {
		a, ok := p.line(q + hint)
		if !ok {
			return def
		}
		switch strings.ToLower(a) {
		case "":
			return def
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
		fmt.Fprintln(p.out, "Please answer y or n.")
	}
}

// docsDir asks where the records go until the answer is a valid folder name;
// an empty answer or end of input takes docs.
func (p *prompter) docsDir() string {
	for {
		a, ok := p.line("Where should the records go? (" + initrepo.DefaultDocsDir + "): ")
		if !ok || a == "" {
			return initrepo.DefaultDocsDir
		}
		err := initrepo.ValidateDocsDir(a)
		if err == nil {
			return a
		}
		fmt.Fprintf(p.out, "%q %v.\n", a, err)
	}
}

// choice asks "q [a/b] (def): " until the answer is one of valid or empty.
func (p *prompter) choice(q string, valid []string, def string) string {
	for {
		a, ok := p.line(fmt.Sprintf("%s [%s] (%s): ", q, strings.Join(valid, "/"), def))
		if !ok || a == "" {
			return def
		}
		for _, v := range valid {
			if strings.EqualFold(a, v) {
				return v
			}
		}
		fmt.Fprintf(p.out, "Please answer one of: %s.\n", strings.Join(valid, ", "))
	}
}
