// Package aggregator – stub for Concurrent File Stats Processor.
package aggregator

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Result mirrors one JSON object in the final array.
type Result struct {
	Path   string `json:"path"`
	Lines  int    `json:"lines,omitempty"`
	Words  int    `json:"words,omitempty"`
	Status string `json:"status"` // "ok" or "timeout"
}

type job struct {
	Index    int
	Path     string // texts/01_zen.txt
	FullPath string // absolute path e.g. data/texts/01_zen.txt
}

// Aggregate must read filelistPath, spin up *workers* goroutines,
// apply a per‑file timeout, and return results in **input order**.
// aggregator --workers=8 --timeout=2 data/filelist.txt  ➜  result.json
func Aggregate(filelistPath string, workers, timeout int) ([]Result, error) {
	// Open filelist.txt
	file, err := os.Open(filelistPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	baseDir := filepath.Dir(filelistPath)

	// Read file paths
	scanner := bufio.NewScanner(file)
	var jobsList []job
	idx := 0
	for scanner.Scan() {
		relPath := scanner.Text()
		fullPath := filepath.Join(baseDir, relPath)
		jobsList = append(jobsList, job{
			Index:    idx,
			Path:     relPath,
			FullPath: fullPath,
		})
		idx++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	results := make([]Result, len(jobsList)) // should be in order
	jobsCh := make(chan job)
	var wg sync.WaitGroup

	// workers receive data from job channel
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobsCh {
				results[j.Index] = processFileWithTimeout(j.Path, j.FullPath, timeout)
			}
		}()

	}

	// send jobs into jobs channel
	go func() {
		for _, j := range jobsList {
			jobsCh <- j
		}
		close(jobsCh)
	}()

	wg.Wait()
	return results, nil
}

func processFileWithTimeout(relPath, fullPath string, inputTimeout int) Result {
	f, err := os.Open(fullPath)
	if err != nil {
		// treat error as timeout
		return Result{Path: relPath, Status: "timeout"}
	}
	defer f.Close()

	// input timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(inputTimeout)*time.Second)
	defer cancel()

	scanner := bufio.NewScanner(f)

	if scanner.Scan() {
		firstLine := scanner.Text()
		if strings.HasPrefix(firstLine, "#sleep=") {
			parts := strings.SplitN(firstLine, "=", 2)
			if len(parts) == 2 {
				if dur, err := time.ParseDuration(parts[1] + "s"); err == nil {
					if dur > time.Duration(inputTimeout)*time.Second {
						// If sleep duration is greater than timeout, return timeout immediately
						return Result{Path: relPath, Status: "timeout"}
					}
					// Otherwise, sleep for the requested duration
					time.Sleep(dur)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Result{Path: relPath, Status: "timeout"}
	}

	resCh := make(chan Result, 1)
	go func() {
		lines := []string{}
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			// treat error as timeout
			resCh <- Result{Path: relPath, Status: "timeout"}
			return
		}

		wordCount := 0
		for _, line := range lines {
			wordCount += len(strings.Fields(line))
		}
		resCh <- Result{
			Path:   relPath,
			Lines:  len(lines),
			Words:  wordCount,
			Status: "ok",
		}
	}()

	select {
	case <-ctx.Done():
		return Result{Path: relPath, Status: "timeout"}
	case res := <-resCh:
		return res
	}
}
