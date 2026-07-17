package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		p := Path{Name: name}
		p.Full, p.Err = exec.LookPath(p.Name)
		if eval && p.Err == nil {
			if r, e := scoopFind(p.Name); e {
				p.Real, p.Err = r, nil
			} else if r, e := miseFind(p.Name); e {
				p.Real, p.Err = r, nil
			} else {
				p.Real, p.Err = filepath.EvalSymlinks(p.Full)
			}
		}
		if p.Real == p.Full {
			p.Real = ""
		}
		if p.Err != nil {
			if errors.Is(p.Err, exec.ErrNotFound) {
				p.ErrMsg = "executable file not found"
			} else {
				p.ErrMsg = p.Err.Error()
			}
		}
		result.Length.Name = max(result.Length.Name, len(p.Name))
		result.Length.Full = max(result.Length.Full, len(p.Full))
		result.Length.Real = max(result.Length.Real, len(p.Real))
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

func linePrint(result Result, i int) {
	p := result.Path[i]
	var w io.Writer
	if p.Err != nil {
		w = os.Stderr
	} else {
		w = os.Stdout
	}

	if p.Err == nil && !result.Verbose {
		if p.Real != "" {
			fmt.Fprintf(w, "%s\n", p.Real)
		} else {
			fmt.Fprintf(w, "%s\n", p.Full)
		}
		return
	}

	if p.Err != nil {
		if !result.NoErr {
			fmt.Fprintf(w, "%s: %s\n", p.Name, p.ErrMsg)
		}
		return
	}

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
	}
	Path []Path
}

type Path struct {
	Name   string `json:"name"`
	Full   string `json:"full,omitempty"`
	Real   string `json:"real,omitempty"`
	Err    error  `json:"-"`
	ErrMsg string `json:"err,omitempty"`
}

func miseFind(name string) (string, bool) {
	n, err := exec.Command("mise", "which", name).Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(n)), true
}

func scoopFind(name string) (string, bool) {
	n, err := exec.Command("scoop", "which", name).Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(n)), true
}
