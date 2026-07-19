package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

func main() {
	var eval, verbose, fmtjson, noe bool
	pflag.Usage = func() { fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0])); pflag.PrintDefaults() }
	pflag.BoolVarP(&eval, "eval", "e", eval, "eval link if symlink")
	pflag.BoolVarP(&fmtjson, "json", "j", fmtjson, "json output")
	pflag.BoolVar(&noe, "noe", noe, "no err output")
	pflag.BoolVarP(&verbose, "verbose", "v", verbose, "verbose")
	pflag.Parse()

	names := pflag.Args()
	if len(names) == 0 {
		pflag.Usage()
		os.Exit(1)
	}

	var result Result
	result.Verbose = verbose
	result.NoErr = noe
	result.Eval = eval

	for _, name := range names {
		p := Find(name, eval)
		result.Length.Name = max(result.Length.Name, len(p.Name))
		result.Length.Full = max(result.Length.Full, len(p.Full))
		result.Length.Real = max(result.Length.Real, len(p.Real))
		result.Length.Shim = max(result.Length.Shim, len(p.Shim))
		result.Path = append(result.Path, p)
	}

	if fmtjson {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result.Path)
		return
	}

	result.Verbose = result.Verbose || len(result.Path) > 1
	for i := range result.Path {
		linePrint(result, i)
	}
}

func Find(name string, eval bool) (p Path) {
	p.Name = name

	if p.Name != "" {
		p.Full, _ = exec.LookPath(p.Name)
	}

	if eval {
		if p.Full != "" {
			p.Real, _ = filepath.EvalSymlinks(p.Full)
		}

		if p.Real != "" {
			p.Shim = shim(p.Real)
		}

		if p.Shim == p.Real {
			p.Shim = ""
		}

		if p.Real == p.Full {
			p.Real = ""
		}
	}

	return
}

func linePrint(result Result, i int) {
	p := result.Path[i]

	if p.Full == "" {
		if !result.NoErr {
			fmt.Fprintf(os.Stderr, "%s: not found\n", p.Name)
		}
		return
	}

	if !result.Verbose {
		fmt.Fprintf(os.Stdout, "%s\n", cmp.Or(p.Shim, p.Real, p.Full))
		return
	}

	w := os.Stdout

	// name
	fmt.Fprintf(w, "%*s", -result.Length.Name, p.Name)

	if result.Length.Full > 0 {
		if p.Full != "" {
			fmt.Fprintf(w, " => %*s", -result.Length.Full, p.Full)
		} else {
			fmt.Fprintf(w, "    %*s", -result.Length.Full, "")
		}
	}

	if result.Length.Real > 0 {
		if p.Real != "" {
			fmt.Fprintf(w, " => %*s", -result.Length.Real, p.Real)
		} else {
			fmt.Fprintf(w, "    %*s", -result.Length.Real, "")
		}
	}

	if result.Length.Shim > 0 {
		if p.Shim != "" {
			fmt.Fprintf(w, " => %*s", -result.Length.Shim, p.Shim)
		} else {
			fmt.Fprintf(w, "    %*s", -result.Length.Shim, "")
		}
	}

	fmt.Fprintln(w)
}

type Result struct {
	Verbose bool
	NoErr   bool
	Eval    bool
	Length  struct {
		Name int
		Full int
		Real int
		Shim int
	}
	Path []Path
}

type Path struct {
	Name string `json:"name"`
	Full string `json:"full,omitempty"`
	Real string `json:"real,omitempty"`
	Shim string `json:"shim,omitempty"`

	// Err    error  `json:"-"`
	// ErrMsg string `json:"err,omitempty"`
}

func shim(name string) (r string) {
	if name == "" {
		return
	}
	if strings.Contains(name, "shims") {
		if r = which("mise", name); r == "" || strings.Contains(r, "shims") {
			r = which("scoop", name)
		}
	}
	return
}

func which(command, name string) string {
	n, err := exec.Command(command, "which", filepath.Base(name)).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(n))
}
