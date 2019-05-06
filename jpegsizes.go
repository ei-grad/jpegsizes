package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	nworkers := flag.Int("workers", 4, "number of images to process in parallel")
	pattern := flag.String("pattern", "*.jpg", "glob pattern to match image filenames")
	flag.Parse()

	tasks := make(chan string)
	out := make(chan string)

	processImage := func(imgName string) {
		imgFile, err := os.Open(imgName)
		if err != nil {
			log.Printf("Can't open %s: %s", imgName, err)
			return
		}
		defer imgFile.Close()
		cfg, err := jpeg.DecodeConfig(imgFile)
		if err != nil {
			log.Printf("Can't get %s size: %s", imgName, err)
			return
		}
		out <- fmt.Sprintf("%s\t%d\t%d\n", imgName, cfg.Width, cfg.Height)
	}

	// run workers
	var workersWG sync.WaitGroup
	for i := 0; i < *nworkers; i++ {
		workersWG.Add(1)
		go func() {
			defer workersWG.Done()
			for imgName := range tasks {
				processImage(imgName)
			}
		}()
	}

	go func() {
		defer close(tasks)
		matches, err := filepath.Glob(*pattern)
		if err != nil {
			log.Printf("Can't get files list: %s", err)
			return
		}
		if matches == nil {
			log.Print("No files found.")
			return
		}
		for i := 0; i < len(matches); i++ {
			tasks <- matches[i]
		}
	}()

	var outWG sync.WaitGroup
	outWG.Add(1)
	go func() {
		defer outWG.Done()
		for i := range out {
			fmt.Print(i)
		}
	}()

	workersWG.Wait()
	close(out)

	outWG.Wait()
}
