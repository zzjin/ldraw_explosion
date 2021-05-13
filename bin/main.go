package main

import (
	"log"
	"os"
	"path"
	"strings"

	ldraw "github.com/zzjin/ldraw_explosion"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Param error,pls drag file on.\nAuthor: zzjin tczzjin#gmail.com\n")
	}
	fileName := os.Args[len(os.Args)-1]

	if path.Ext(fileName) != ".ldr" {
		log.Fatal("file not supportted, pls drag ldr file on.\nAuthor: zzjin tczzjin#gmail.com\n")
	}

	// parse ldraw file
	mainFile := ldraw.NewRawFile()
	ldraw.ParseLdrContent(fileName, mainFile)
	// merge sub inline files into parts
	allParts := ldraw.ReplaceSubFiles(mainFile, &mainFile.SubFiles)

	outName := strings.Replace(fileName, ".ldr", "_ground.ldr", 1)
	ldraw.NewPackParts(allParts).Save(outName)
}
