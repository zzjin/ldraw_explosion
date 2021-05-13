package ldraw

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	parsedFile sync.Map
	pLock      sync.Mutex
)

// ParseDatFile ParseDatFile
func ParseDatFile(fileName string, matrix *TransMatrix, ldrawRoot string) *BoundingBox {
	resp := NewBoundingBox()

	// sync.Map parse sub file once
	if got, ok := parsedFile.Load(fileName); ok {
		bb := got.(*BoundingBox)
		// transform to new bounding box
		resp.MergeMinMaxVector(MultipleVector(matrix, bb.Min, bb.Max)...)

		return resp
	}

	oneReader, errF := os.Open(fileName)
	if errF != nil {
		log.Fatalf("Open ldr file failed: %v.\n", errF)
	}
	defer oneReader.Close()

	lineNum := 0
	// start reading from the file with a reader.
	reader := bufio.NewReader(oneReader)
	var line string
	var err error
	for {
		lineNum++

		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		line = strings.TrimSpace(strings.ReplaceAll(line, "\r\n", ""))
		if line == "" {
			if err != nil {
				// fail break all
				break
			}
			continue
		}

		values := parseOneLine(line)
		if len(values) == 0 {
			log.Fatalf("file %s wrong line %d: %q\n", fileName, lineNum, line)
		}

		if values[0] == "1" {
			if len(values) < 15 {
				continue
			} else if len(values) > 15 {
				// as file name can contain space, join together
				values[14] = strings.Join(values[14:], " ")
			}

			//cast all to lower
			values[14] = strings.ToLower(values[14])

			fileRealLocation := getSubFileRealLocation(values[14], ldrawRoot)
			subFileBoundingBox := ParseDatFile(fileRealLocation, InitMatrix, ldrawRoot)
			parsedFile.Store(fileRealLocation, subFileBoundingBox)

			// calc all parent matrix
			subFileMatrix := NewTransMatrixFromStrs(values[2:14])
			subFileMatrixMulti := MultipleMatrix(matrix, subFileMatrix)

			// apply to sub file bouding-box
			subFileVectorMulti := MultipleVector(subFileMatrixMulti, subFileBoundingBox.Min, subFileBoundingBox.Max)

			resp.MergeMinMaxVector(subFileVectorMulti...)

		} else if values[0] == "2" {
			vectors := NewVectorsFromLine(values[2:8], 2)
			resp.MergeMinMaxVector(MultipleVector(matrix, vectors...)...)
		} else if values[0] == "3" {
			vectors := NewVectorsFromLine(values[2:11], 3)
			resp.MergeMinMaxVector(MultipleVector(matrix, vectors...)...)
		} else if values[0] == "4" {
			vectors := NewVectorsFromLine(values[2:14], 4)
			resp.MergeMinMaxVector(MultipleVector(matrix, vectors...)...)
		}
		// currently do not need parse type "5"

		if err != nil && err == io.EOF {
			break
		}
	}
	if err != io.EOF {
		log.Fatalf("Parse file failed with error: %s\n", err)
	}

	return resp.TransEmpty()
}

// ParseLdrContent ParseLdrContent
func ParseLdrContent(fileName string, mainFile *RawFile) {
	oneReader, errF := os.Open(fileName)
	if errF != nil {
		log.Fatalf("Open ldr file failed: %v.\n", errF)
	}
	defer oneReader.Close()
	// Start reading from the file with a reader.
	reader := bufio.NewReader(oneReader)

	isFirstFile := true
	lfNum := 1
	lineNum := 1
	workingFileName := ""

	var line string
	var err error
	for {
		lineNum++

		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		line = strings.TrimSpace(strings.ReplaceAll(line, "\r\n", ""))
		if line == "" {
			if err != nil {
				// fail break all
				break
			}
			continue
		}

		values := parseOneLine(line)
		if len(values) == 0 {
			log.Fatalf("file %s wrong line %d: %q\n", fileName, lineNum, line)
		}

		if lfNum == 1 {
			// remove optional utf8-bom
			if len(values[0]) != 1 {
				values[0] = values[0][len(values[0])-1:]
			}

			if len(values) < 2 || values[0] != `0` {
				log.Fatalf("parse line fail: %s", line)
			}

			if values[1] == "FILE" && len(values) >= 3 {
				// contains sub file, mpd?
				workingFileName = strings.Join(values[2:], " ")
			} else {
				workingFileName = strings.Join(values[1:], " ")
			}
			workingFileName = strings.ToLower(workingFileName)

			if isFirstFile {
				mainFile.Name = workingFileName
			} else {
				// is inline sub-file
				mainFile.SubFiles[workingFileName] = NewRawFile()
			}

			lfNum++
			continue
		}

		if values[0] == "0" && values[1] == "NOFILE" {
			// main file end, start parse all files
			isFirstFile = false
			workingFileName = ""

			lfNum = 1 // reset inline file numer
			continue
		}

		// TODO: custom hose?

		// other command push lines to data
		if isFirstFile {
			parseInlineFilePart(values, mainFile)
		} else {
			parseInlineFilePart(values, mainFile.SubFiles[workingFileName])
		}

		lfNum++
	}
	if err != io.EOF {
		log.Fatalf("Parse file failed with error: %s\n", err)
	}
}

func parseInlineFilePart(v []string, target *RawFile) {
	if v[0] == "1" {
		if len(v) < 15 {
			return
		} else if len(v) > 15 {
			// as file name can contain space, join together
			v[14] = strings.Join(v[14:], " ")
		}

		//cast all to lower
		v[14] = strings.ToLower(v[14])

		// in ldraw p(sub) dir do not parse
		if _, ok := AllP[v[14]]; ok {
			return
		}

		// in ldraw parts list
		var id string
		if _, ok := AllParts[v[14]]; ok {
			id = v[14][:len(v[14])-4]
		} else {
			// inline sub file or custom parts
			id = v[14]
		}

		k := id + "-" + v[1]

		pLock.Lock()
		if fp, ok := target.Parts[k]; ok {
			fp.Count++
		} else {
			colorInt, _ := strconv.Atoi(v[1])
			target.Parts[k] = &Part{ID: id, Color: colorInt, Count: 1}
		}
		pLock.Unlock()
	}
}

// parseOneLine parse line to command(s)
func parseOneLine(line string) []string {
	// clean up unwanted `tab` usage
	lineClean := strings.TrimSpace(strings.ReplaceAll(line, "\t", " "))
	return strings.Fields(lineClean)
}

func ReplaceSubFiles(rawFile *RawFile, subFiles *map[string]*RawFile) map[string]*Part {
	resp := map[string]*Part{}
	for k, part := range rawFile.Parts {
		if subFile, ok := (*subFiles)[part.ID]; ok {
			subResp := ReplaceSubFiles(subFile, subFiles)
			for k2, part2 := range subResp {
				if part2.Color == 16 {
					// replace with parent color
					part2.Color = part.Color
					k2 = part2.ID + "-" + strconv.FormatInt(int64(part.Color), 10)
				}

				if _, ok2 := resp[k2]; !ok2 {
					resp[k2] = part2
				} else {
					resp[k2].Count += part2.Count
				}
			}
			delete(rawFile.Parts, k) // remove self
		} else {
			if _, ok2 := resp[k]; !ok2 {
				resp[k] = part
			} else {
				resp[k].Count += part.Count
			}
		}
	}
	return resp
}
