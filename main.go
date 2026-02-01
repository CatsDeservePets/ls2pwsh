package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

var progName = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")

type style string

func (s style) String() string {
	return "\"`e[" + string(s) + "m\""
}

type fileInfo struct {
	directory    style
	symbolicLink style
	executable   style
	extension    map[string]style
}

func main() {
	args, err := readInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage: %s LS_COLORS\n", progName)
		os.Exit(2)
	}
	fi := fromLSCOLORS(args)
	fmt.Println(fi.toPWSH())
}

func readInput() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	if len(os.Args) > 1 {
		return os.Args[1], nil
	}

	return "", fmt.Errorf("no input")
}

func fromLSCOLORS(s string) fileInfo {
	var fi fileInfo

	for ent := range strings.SplitSeq(s, ":") {
		k, v, found := strings.Cut(ent, "=")
		if !found {
			continue
		}
		switch k {
		case "di":
			fi.directory = style(v)
		case "ln":
			fi.symbolicLink = style(v)
		case "ex":
			fi.executable = style(v)
		default:
			if after, ok := strings.CutPrefix(k, "*."); ok {
				ext := "." + after
				if fi.extension == nil {
					fi.extension = make(map[string]style)
				}
				fi.extension[ext] = style(v)
			}
		}
	}

	return fi
}

func (fi fileInfo) toPWSH() string {
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
