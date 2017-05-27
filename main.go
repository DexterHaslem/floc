// floc is a quick and dirty to list *PHYSICAL*
// lines of code in each directory
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const FilterNone = "*"

var dir = flag.String("dir", ".", "Starting directory")
var filter = flag.String("filter", FilterNone, "extension filter, a comma separated list of extensions to include, or * for all")
var validExts []string

type dirinfo struct {
	name    string
	lines   int64
	parent  *dirinfo
	subdirs []*dirinfo
}

func (d *dirinfo) sublines() int64 {
	total := int64(0)
	for _, sd := range d.subdirs {
		total += sd.lines + sd.sublines()
	}
	return total
}

func (d *dirinfo) pretty(tabs int) {
	for i := 0; i < tabs; i++ {
		// doing \t was huge
		fmt.Print("  ")
	}

	fmt.Printf("%s: contains=%d (%d in dir)\n", d.name, d.sublines()+d.lines, d.lines)
	for _, sd := range d.subdirs {
		sd.pretty(tabs + 1)
	}
}

func parseExts() {
	validExts = []string{}
	f := *filter
	if f == FilterNone {
		return
	}

	chunks := strings.Split(f, ",")
	for _, c := range chunks {
		// filepath.Ext returns dot so add it here to simplify
		validExts = append(validExts, "."+c)
	}
}

func passes(fi os.FileInfo) bool {
	// TODO: try to skip binary
	if *filter == FilterNone {
		return true
	}
	// ext returns the dot so massasge things a bit

	ext := filepath.Ext(fi.Name())
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

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return -1, err
	}

	str := string(b)
	chunks := strings.Split(str, "\n") // TODO: line ending modes :-(
	return len(chunks), nil
}

func walk(root string, parent *dirinfo) (*dirinfo, error) {
	fs, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	ret := dirinfo{
		name:    filepath.Base(root),
		subdirs: []*dirinfo{},
		parent:  parent,
	}

	for _, fi := range fs {
		fn := filepath.Join(root, fi.Name())
		if fi.IsDir() {
			subdir, err := walk(fn, &ret)
			if err == nil {
				ret.subdirs = append(ret.subdirs, subdir)
			}
		} else if passes(fi) {
			ln, err := lines(fn)
			if err == nil {
				ret.lines += int64(ln)
			}
		}
	}
	return &ret, nil
}

func main() {
	flag.Parse()
	parseExts()
	fmt.Printf("floc: dir='%s' filter='%s'\n", *dir, *filter)
	wd, err := walk(*dir, nil)
	if err == nil {
		wd.pretty(0)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}
