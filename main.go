// floc is a tool to list *PHYSICAL* lines of code in each directory
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FilterNone = "*"

var (
	dir = flag.String("dir", ".",
		"Starting directory")

	filter = flag.String("filter", FilterNone,
		"extension filter, a comma separated list of extensions to include (with no dot), or * for all")
	validExts []string
)

func parseExts() {
	validExts = []string{}
	f := *filter
	if f == FilterNone {
		return
	}

	chunks := strings.Split(f, ",")
	for _, c := range chunks {
		// if user adds dot, clean it so its consistent below
		validExts = append(validExts, strings.Trim(c, "."))
	}
}

func passes(ext string) bool {
	if *filter == FilterNone {
		return true
	}

	for _, ve := range validExts {
		if ve == ext {
			return true
		}
	}
	return false
}

func lines(fn string) (int, error) {
	f, err := os.Open(fn)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var c int
	for c = 0; s.Scan(); c++ {
		// dont bother, use file filters
		//if notText(l)
		// whee
	}

	return c, nil
}

func main() {
	if len(os.Args) < 3 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()
	parseExts()
	fmt.Printf("floc: parsing dir='%s' for file filter='%s'\n", *dir, *filter)

	type ret struct {
		// ext:[path: loc]
		ByDir map[string]map[string]int `json:"byDir"`
		// ext: total
		Total map[string]int `json:"total"`
	}
	r := &ret{
		ByDir: map[string]map[string]int{},
	}

	filepath.Walk(*dir, func(p string, i os.FileInfo, err error) error {
		if i.IsDir() {
			// this causes empty dirs to be reported
			//lines[p] = 0
			return nil
		}

		// strip dot for printing
		ext := strings.Trim(filepath.Ext(i.Name()), ".")

		if passes(ext) {
			d := filepath.Dir(p)
			// chop off root so its relative to our start
			d = strings.Replace(d, *dir, "", -1)
			d = strings.Trim(d, "\\/")
			d = filepath.Clean(d)
			lc, err := lines(p)
			if err == nil {
				if r.ByDir[d] == nil {
					r.ByDir[d] = map[string]int{}
				}
				r.ByDir[d][ext] += lc
			} else {
				fmt.Fprintf(os.Stderr, "failed to read file %s: %s", i.Name(), err)
			}
		}
		return nil
	})

	r.Total = map[string]int{}
	for _, m := range r.ByDir {
		for e, v := range m {
			r.Total[e] += v
		}
	}

	j, err := json.MarshalIndent(r, "", "  ")
	if err == nil {
		fmt.Printf("%s\n", j)
	}
}
