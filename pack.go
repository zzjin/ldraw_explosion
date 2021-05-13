package ldraw

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"github.com/zzjin/ldraw_explosion/binpack"
)

type LdrPackPart struct {
	Name    string
	Color   int
	X, Y    int
	W, H, T float64
}

func (ldrp *LdrPackPart) CalcSize() (int, int) {
	// give more width and height for spacing
	w := int(math.Ceil((ldrp.W+20)/20) * 20 * 1.5)
	h := int(math.Ceil((ldrp.H+20)/20) * 20 * 1.5)
	return w, h
}

func (ldrp *LdrPackPart) StandLine() string {
	calcW, calcH := ldrp.CalcSize()

	offsetX := ldrp.X + int(calcW/2) // +x
	offsetY := -int(ldrp.T / 2)      // -y is upper
	offSetZ := ldrp.Y + int(calcH/2) // +x

	return fmt.Sprintf("1 %d %d %d %d %s %s\n", ldrp.Color, offsetX, offsetY, offSetZ, DefaultXMatrix, ldrp.Name)
}

type LdrBinPack []*LdrPackPart

// NewPackParts NewPackParts
func NewPackParts(partMap map[string]*Part) *LdrBinPack {
	parts := LdrBinPack{}
	for _, one := range partMap {
		// ignore none offical part
		name := one.ID + ".dat"
		v, ok := AllParts[name]
		if !ok {
			log.Printf("brick not found: %s\n", name)
			continue
		}

		w, h, t := GetBoxWHTByX(v)

		for i := 0; i < one.Count; i++ {
			parts = append(parts, &LdrPackPart{
				Name: name, Color: one.Color,
				X: 0, Y: 0,
				W: w, H: h, T: t,
			})
		}
	}

	// sort size max->min
	sort.Slice(parts, func(i, j int) bool {
		iw, ih := parts[i].CalcSize()
		jw, jh := parts[j].CalcSize()
		return iw*ih > jw*jh
	})

	return &parts
}

func (lbp LdrBinPack) Len() int {
	return len(lbp)
}

func (lbp LdrBinPack) Size(n int) (int, int) {
	return lbp[n].CalcSize()
}

func (lbp LdrBinPack) Place(n, x, y int) {
	lbp[n].X, lbp[n].Y = x, y
}

const defaultLDrGroundHeader = `0 Untitled Model
0 Name:  ground
0 Author: origin author && zzjin#tczzjin@gmail.com
0 CustomBrick
`

func (lbp *LdrBinPack) Save(fileName string) {
	outputW, outputH := binpack.Pack(lbp)
	fmt.Printf("output: %dx%d\n", outputH, outputW)

	var wf *os.File
	var err error
	if wf, err = os.Create(fileName); err != nil {
		log.Fatal(err)
	}
	defer wf.Close()

	if _, err := wf.WriteString(defaultLDrGroundHeader); err != nil {
		log.Fatal(err)
	}
	if _, err := wf.WriteString(fmt.Sprintf("0 NumOfBricks:  %d\n", len(*lbp))); err != nil {
		log.Fatal(err)
	}

	for _, one := range *lbp {
		if _, err := wf.WriteString(one.StandLine()); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := wf.WriteString("\n"); err != nil {
		log.Fatal(err)
	}
}
