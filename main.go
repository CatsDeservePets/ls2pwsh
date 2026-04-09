package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

var progName = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")

type fileInfo struct {
	directory    string
	symbolicLink string
	executable   string
	extension    map[string]string
}

type format byte

const (
	unknown format = iota
	pwsh
	gnu
	// TODO: bsd?
)

func (f format) String() string {
	switch f {
	case pwsh:
		return "pwsh"
	case gnu:
		return "gnu"
	default:
		return "unknown"
	}
}

const help = `usage: %s LS_COLORS | PSStyle.FileInfo

Convert color strings between LS_COLORS and PowerShell PSStyle.FileInfo format

If the input is a single dash ('-') or absent, %s reads from the standard input.
`

func usage() {
	fmt.Fprintf(os.Stderr, help, progName, progName)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args, err := readInput(flag.Args(), os.Stdin)
	if err != nil {
		flag.Usage()
	}
	switch detectFormat(args) {
	case gnu:
		fi := fromLSCOLORS(args)
		fmt.Println(fi.toPSStyle())
	case pwsh:
		fi := fromPSStyle(args)
		fmt.Println(fi.toLSCOLORS())
	default:
		fmt.Fprintf(os.Stderr, "%s: unrecognized input format\n", progName)
		os.Exit(1)
	}
}

func readInput(args []string, f *os.File) (string, error) {
	if len(args) > 0 && args[0] != "-" {
		return args[0], nil
	}

	if !term.IsTerminal(int(f.Fd())) {
		b, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return "", fmt.Errorf("no input")
}

func detectFormat(s string) format {
	switch {
	case strings.Contains(s, "di="), strings.Contains(s, "ln="), strings.Contains(s, "ex="):
		return gnu
	case strings.Contains(s, "Directory"), strings.Contains(s, "SymbolicLink"), strings.Contains(s, "Executable"), strings.Contains(s, "Extension"):
		return pwsh
	default:
		return unknown
	}
}

func fromPSStyle(s string) fileInfo {
	var fi fileInfo

	unescape := func(s string) string {
		return strings.TrimSuffix(strings.TrimPrefix(strings.Trim(s, "\""), "`e["), "m")
	}

	for ent := range strings.Lines(s) {
		ent = strings.TrimSpace(strings.ReplaceAll(ent, " ", ""))
		k, v, _ := strings.Cut(ent, ":")
		v = unescape(v)
		switch k {
		case "Directory":
			fi.directory = v
		case "SymbolicLink":
			fi.symbolicLink = v
		case "Executable":
			fi.executable = v
		case "Extension":
			// The first pair appears after 'Extension:' on the same line.
			k = v
			fallthrough
		default:
			ext, val, found := strings.Cut(k, "=")
			if !found {
				continue
			}
			if fi.extension == nil {
				fi.extension = make(map[string]string)
			}
			fi.extension[ext] = unescape(val)
		}
	}

	return fi
}

func fromLSCOLORS(s string) fileInfo {
	var fi fileInfo

	escape := func(s string) string {
		return "\"`e[" + s + "m\""
	}

	for ent := range strings.SplitSeq(s, ":") {
		k, v, found := strings.Cut(ent, "=")
		if !found {
			continue
		}
		v = escape(v)
		switch k {
		case "di":
			fi.directory = v
		case "ln":
			fi.symbolicLink = v
		case "ex":
			fi.executable = v
		default:
			if after, ok := strings.CutPrefix(k, "*."); ok {
				ext := "." + after
				if fi.extension == nil {
					fi.extension = make(map[string]string)
				}
				fi.extension[ext] = v
			}
		}
	}

	return fi
}

func (fi fileInfo) toPSStyle() string {
	var b strings.Builder

	if fi.directory != "" {
		fmt.Fprintf(&b, "$PSStyle.FileInfo.Directory = %s\n", fi.directory)
	}
	if fi.symbolicLink != "" {
		fmt.Fprintf(&b, "$PSStyle.FileInfo.SymbolicLink = %s\n", fi.symbolicLink)
	}
	if fi.executable != "" {
		fmt.Fprintf(&b, "$PSStyle.FileInfo.Executable = %s\n", fi.executable)
	}
	if len(fi.extension) != 0 {
		for k, v := range fi.extension {
			fmt.Fprintf(&b, "$PSStyle.FileInfo.Extension['%s'] = %s\n", k, v)
		}
	}

	return b.String()
}

func (fi fileInfo) toLSCOLORS() string {
	var b strings.Builder

	if fi.directory != "" {
		fmt.Fprintf(&b, "di=%s:", fi.directory)
	}
	if fi.symbolicLink != "" {
		fmt.Fprintf(&b, "ln=%s:", fi.symbolicLink)
	}
	if fi.executable != "" {
		fmt.Fprintf(&b, "ex=%s:", fi.executable)
	}
	if len(fi.extension) != 0 {
		for k, v := range fi.extension {
			fmt.Fprintf(&b, "*%s=%s:", k, v)
		}
	}

	return b.String()
}
