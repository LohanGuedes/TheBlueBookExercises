package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	countsFiles := make(map[string]map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, countsFiles[""])
	} else {
		for _, arg := range files {
			count := make(map[string]int)
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, count)
			countsFiles[arg] = count
			f.Close()
		}
	}
	for fileName, file := range countsFiles {
		for line, n := range file {
			if n > 1 {
				fmt.Printf("%s\t %d\t%s\n", fileName, n, line)
			}
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}
