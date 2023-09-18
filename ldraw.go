package ldraw

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go run ./generate/generate.go /home/zzjin/projects/lego/ldraw/

const (
	PLocation     = `p/`
	PartsLocation = `parts/`

	UnOfficialLocation = `UnOfficial/`
)

var (
	p4Location     = PLocation + `4/`
	p8Location     = PLocation + `8/`
	p48Location    = PLocation + `48/`
	partsSLocation = PartsLocation + `s/`
	pLocations     = []string{
		PLocation,
		p4Location,
		p8Location,
		p48Location,
		PartsLocation,
		partsSLocation,
	}
)

// getSubFileRealLocation getSubFileRealLocation
func getSubFileRealLocation(filePath, ldrawRoot string) string {
	filePath = strings.Replace(filePath, "\\", "/", -1)

	for _, p := range pLocations {
		cp := filepath.Clean(ldrawRoot + p + filePath)
		if _, err := os.Stat(cp); err == nil {
			return cp
		}
	}

	// Additional check UnOfficial dir path
	unOfficialRoot := ldrawRoot + UnOfficialLocation
	for _, p := range pLocations {
		cp := filepath.Clean(unOfficialRoot + p + filePath)
		if _, err := os.Stat(cp); err == nil {
			return cp
		}
	}

	log.Fatalf("sub file not found: %s\n", filePath)
	return ""
}

// LdrInfo Ldr Full Info
type LdrInfo struct {
	P     map[string]struct{}
	Parts map[string][2][3]float64
}

//go:embed ldraw_aio.gob
var ldrawAIOGob []byte

var AllP, AllParts = func() (map[string]struct{}, map[string][2][3]float64) {
	got := &LdrInfo{}
	gob.NewDecoder(bytes.NewBuffer(ldrawAIOGob)).Decode(&got)
	return got.P, got.Parts
}()

// RawFile RawFile
type RawFile struct {
	Name     string
	Parts    map[string]*Part
	SubFiles map[string]*RawFile
}

// NewRawFile NewRawFile
func NewRawFile() *RawFile {
	return &RawFile{Name: "", Parts: map[string]*Part{}, SubFiles: map[string]*RawFile{}}
}

// Part Part
type Part struct {
	ID    string
	Color int
	Count int
}
