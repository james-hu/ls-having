package lsh

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gobwas/glob"
)

type Options struct {
	Depth        int
	Exclude      glob.Glob
	FlagFile     glob.Glob
	CheckFile    string // "" means it is not used
	CheckRegexp  regexp.Regexp
	CheckInverse bool
}

func LsHaving(options *Options, dir string) []string {
	var found []string = make([]string, 0, 100)
	doLsHaving(options, &found, dir, 0, nil) // root dir has depth 0
	return found
}

type dirEntryEx struct {
	Path  string
	Depth int
	Entry fs.DirEntry
}

// Read all entries under the specified directory.
// The "depth" parameter is the depth of the directory specified by "dir" parameter.
func readEntries(dir string, depth int) *[]dirEntryEx {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	entriesEx := make([]dirEntryEx, 0, len(entries))
	for _, entry := range entries {
		entriesEx = append(entriesEx, dirEntryEx{filepath.Join(dir, entry.Name()), depth + 1, entry})
	}
	return &entriesEx
}

func doLsHaving(options *Options, found *[]string, dir string, depth int, entriesInDir *[]dirEntryEx) {
	if entriesInDir == nil {
		entriesInDir = readEntries(dir, depth)
	}
	for _, entry := range *entriesInDir {
		if shouldCheck(options, &entry) {
			entriesInSubDir := readEntries(entry.Path, entry.Depth)
			if match(options, &entry, entriesInSubDir) {
				*found = append(*found, entry.Path)
			}
			doLsHaving(options, found, entry.Path, entry.Depth, entriesInSubDir)
		}
	}
}

func shouldCheck(options *Options, dir *dirEntryEx) bool {
	if !dir.Entry.IsDir() {
		return false
	}
	if options.Depth >= 0 && options.Depth < dir.Depth {
		return false
	}
	if options.Exclude != nil && options.Exclude.Match(dir.Path) {
		return false
	}
	return true
}

func match(options *Options, dir *dirEntryEx, entries *[]dirEntryEx) bool {
	foundFlagFile := false
	for _, entry := range *entries {
		if options.FlagFile.Match(entry.Entry.Name()) {
			foundFlagFile = true
			break
		}
	}
	var checkFileMismatch bool
	if options.CheckFile == "" {
		checkFileMismatch = false
	} else {
		// read the file
		// do regexp match
	}
	return foundFlagFile && !checkFileMismatch
}
