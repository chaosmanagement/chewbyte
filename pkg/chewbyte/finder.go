package chewbyte

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TraverseDirection int

const (
	None TraverseDirection = iota
	TowardsRoot
	TowardsChildren
)

var chrewbyteExtensions = []string{"json", "yml", "yaml", "jsonnet"}

type FinderConfig struct {
	startingPaths     []string // At least one directory path to start the search at
	traverseDirection TraverseDirection

	filenames            []string // When empty - will return all paths mathing extensions. Else will match names (excluding extensions) against provided values
	extensions           []string // Like above but for extensions
	multiLevelExtensions bool     // If true extensions will be separated by the first dot in the path, else last one
}

func newFinderConfig() FinderConfig {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Unable to get current working directory!")
	}

	r := FinderConfig{
		startingPaths:     []string{cwd},
		traverseDirection: TowardsChildren,

		filenames:            []string{},
		extensions:           []string{},
		multiLevelExtensions: false,
	}

	return r
}

func contains[T comparable](needle T, haystack []T) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func isMatching(config FinderConfig, filename string) bool {
	index := -1
	if config.multiLevelExtensions {
		index = strings.Index(filename, ".")
	} else {
		index = strings.LastIndex(filename, ".")
	}

	left, right := "", ""
	if index < 0 {
		left = filename
		right = ""
	} else {
		left = filename[:index]
		right = filename[index+1:]
	}

	if len(config.filenames) > 0 {
		if !contains(left, config.filenames) {
			return false
		}
	}

	if len(config.extensions) > 0 {
		if !contains(right, config.extensions) {
			return false
		}
	}

	return true
}

func getParentDir(dir string) (string, bool) {
	candidate := filepath.Dir(dir)

	// Check whether you can't go up further
	if candidate != dir {
		return candidate, true
	}

	return "", false
}

func FindFiles(config FinderConfig) ([]string, error) {
	r := NewSet[string]()
	q := NewVisitableQueue[string]()
	q.Import(config.startingPaths)

	for !q.IsEmpty() {
		dir, ok := q.Get()
		if !ok {
			return nil, fmt.Errorf("Unable to pop item from the queue")
		}

		if config.traverseDirection == TowardsRoot {
			parent, ok := getParentDir(dir)
			if ok {
				q.Put(parent)
			}
		}

		dirContent, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range dirContent {
			fullPath := filepath.Join(dir, entry.Name())

			if entry.IsDir() {
				if config.traverseDirection == TowardsChildren {
					q.Put(fullPath)
				}
			} else {
				if isMatching(config, entry.Name()) {
					r.Insert(fullPath)
				}
			}
		}
	}

	return r.GetSlice(), nil
}
