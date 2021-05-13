// +build ignore

// This program generates datas.go. It can olny be invoked by running `go generate`
package main

import (
	_ "embed"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	ldraw "github.com/zzjin/ldraw_explosion"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Param error,pls spec ldraw dir.\nAuthor: zzjin tczzjin#gmail.com\n")
	}

	ldrawRoot := os.Args[len(os.Args)-1]
	if ldrawRoot[len(ldrawRoot)-1:] != "/" {
		ldrawRoot += "/"
	}

	pFiles := walkDatDir(ldrawRoot, ldraw.PLocation, false)
	partFiles := walkDatDir(ldrawRoot, ldraw.PartsLocation, true)

	log.Printf("p:%d,part:%d\n", len(pFiles), len(partFiles))

	f, err := os.Create("ldraw_aio.gob")
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	defer f.Close()

	pGob := map[string]struct{}{}
	for k := range pFiles {
		pGob[k] = struct{}{}
	}
	partsGob := map[string][2][3]float64{}
	for k, v := range partFiles {
		partsGob[k] = v.ToGob()
	}

	filesAIO := &ldraw.LdrInfo{P: pGob, Parts: partsGob}
	if err := gob.NewEncoder(f).Encode(filesAIO); err != nil {
		log.Fatalf("Write failed: %v", err)
	}
}

var (
	l       sync.Mutex
	wg      sync.WaitGroup
	numCPUs = runtime.NumCPU()
)

func walkDatDir(ldrawRoot, entryPath string, parseBounding bool) map[string]*ldraw.BoundingBox {
	walkPath := ldrawRoot + entryPath
	ch := make(chan string)

	files := map[string]*ldraw.BoundingBox{}
	worker := func(ch chan string) {
		for path := range ch {
			relaPath := strings.ToLower(strings.TrimPrefix(path, walkPath)) // cast all to lower case

			boundingBox := &ldraw.BoundingBox{}
			if parseBounding {
				log.Printf("parse: %s\n", strings.ReplaceAll(path, ldrawRoot, ""))
				boundingBox = ldraw.ParseDatFile(path, ldraw.InitMatrix, ldrawRoot)
			}

			l.Lock()
			files[relaPath] = boundingBox
			l.Unlock()
		}

		wg.Done()
	}

	// start the workers
	for t := 0; t < numCPUs; t++ {
		wg.Add(1)
		go worker(ch)
	}

	if wErr := filepath.WalkDir(walkPath, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".dat" {
			if strings.Contains(path, "/textures/") || strings.Contains(path, "/s/") {
				return nil
			}

			ch <- path
		}
		return nil
	}); wErr != nil {
		log.Fatal(wErr)
	}

	close(ch)
	wg.Wait()

	return files
}
