package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

var nameLength = 12
var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890.-")
var verbose bool

func main() {
	var n int
	var x int
	var outDir string
	var prefix string

	flag.IntVar(&n, "n", 1, "set the number of files to make")
	flag.IntVar(&x, "x", 36000, "set the size of each file")
	flag.StringVar(&outDir, "o", "output", "set the output directory")
	flag.StringVar(&prefix, "p", "", "file name prefix")
	flag.BoolVar(&verbose, "v", false, "when set, get more verbose output")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	vlog("going to generate %d files that are %d bytes in size in %s with prefix %s", n, x, outDir, prefix)

	vlog("pre-generating names...")
	names := genNames(n, prefix)

	vlog("pre-generating content...")
	content := genContent(x)

	vlog("ensuring output directory %q exists", outDir)

	err := os.Mkdir(outDir, 0755)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		// this could just error because the dir exists
		vlog("unable to create directory, ", err)
	}

	vlog("beginning to write files to %s...", outDir)
	var errs bool
	for _, name := range names {
		err := createFile(outDir, name, content)
		if err != nil {
			errs = true
			log.Println("Unable to create file", err)
		}
	}
	if errs {
		vlog("exiting unsuccessfully")
		os.Exit(1)
	}
	vlog("complete")
}

func genNames(n int, prefix string) []string {
	names := make([]string, n)
	p := []rune(prefix)
	for i := 0; i < n; i++ {
		name := make([]rune, len(p)+nameLength)
		for j := 0; j < len(p)+nameLength; j++ {
			if j < len(p) {
				name[j] = []rune(p)[j]
			} else {
				name[j] = runes[rand.Intn(len(runes))]
			}
		}
		names[i] = string(name)
	}
	return names
}

func genContent(x int) []byte {
	content := make([]rune, x)
	for i := range content {
		content[i] = runes[rand.Intn(len(runes))]
		// give a small chance for new lines
		if rand.Intn(100) <= 2 {
			content[i] = '\n'
		}
	}
	// there has to be a better way
	return []byte(string(content))
}

func createFile(outDir string, name string, content []byte) error {
	fullPath := path.Join(outDir, name)
	vlog("writing file %s", fullPath)
	return ioutil.WriteFile(fullPath, content, 0755)
}

func vlog(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}
