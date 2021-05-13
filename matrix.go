package ldraw

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

// TransMatrix TransMatrix
type TransMatrix [16]float64

/* InitMatrix 1
/ a d g 0 \   / a b c x \
| b e h 0 |   | d e f y |
| c f i 0 |   | g h i z |
\ x y z 1 /   \ 0 0 0 1 /
*/
var InitMatrix = &TransMatrix{
	1.0, 0.0, 0.0, 0.0,
	0.0, 1.0, 0.0, 0.0,
	0.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 1.0,
}

// DefaultXMatrix default ldr file matrix for stand
const DefaultXMatrix = "1 0 0 0 1 0 0 0 1"

// str2F64 inline string to float64
func str2F64(s string) float64 {
	result, _ := strconv.ParseFloat(s, 64)
	return result
}

// NewTransMatrixFromStrs NewTransMatrixFromStrs
func NewTransMatrixFromStrs(d []string) *TransMatrix {
	if len(d) != 12 {
		log.Fatalf("string not match:%v", d)
	}

	return &TransMatrix{
		str2F64(d[3]), str2F64(d[6]), str2F64(d[9]), 0,
		str2F64(d[4]), str2F64(d[7]), str2F64(d[10]), 0,
		str2F64(d[5]), str2F64(d[8]), str2F64(d[11]), 0,
		str2F64(d[0]), str2F64(d[1]), str2F64(d[2]), 1,
	}
}

// NewVectorsFromLine NewVectorsFromLine
func NewVectorsFromLine(d []string, vCount int) []*TransVector {
	num := 3 * vCount
	if len(d) < num {
		return nil
	}

	resp := []*TransVector{}
	for i := 0; i < num; i += 3 {
		resp = append(resp, &TransVector{str2F64(d[i]), str2F64(d[i+1]), str2F64(d[i+2])})
	}
	return resp
}

// MultipleMatrix https://www.migenius.com/articles/3d-transformations-part1-matrices
func MultipleMatrix(l, r *TransMatrix) *TransMatrix {
	result := &TransMatrix{}

	/*
		{\displaystyle \mathbf {C} ={\begin{pmatrix}a_{11}b_{11}+\cdots +a_{1n}b_{n1}&a_{11}b_{12}+\cdots +a_{1n}b_{n2}&\cdots &a_{11}b_{1p}+\cdots +a_{1n}b_{np}\\a_{21}b_{11}+\cdots +a_{2n}b_{n1}&a_{21}b_{12}+\cdots +a_{2n}b_{n2}&\cdots &a_{21}b_{1p}+\cdots +a_{2n}b_{np}\\\vdots &\vdots &\ddots &\vdots \\a_{m1}b_{11}+\cdots +a_{mn}b_{n1}&a_{m1}b_{12}+\cdots +a_{mn}b_{n2}&\cdots &a_{m1}b_{1p}+\cdots +a_{mn}b_{np}\\\end{pmatrix}}}
	*/
	// https://en.wikipedia.org/wiki/Matrix_multiplication
	result[0] = l[0]*r[0] + l[4]*r[1] + l[8]*r[2] + l[12]*r[3]
	result[1] = l[1]*r[0] + l[5]*r[1] + l[9]*r[2] + l[13]*r[3]
	result[2] = l[2]*r[0] + l[6]*r[1] + l[10]*r[2] + l[14]*r[3]
	result[3] = l[3]*r[0] + l[7]*r[1] + l[11]*r[2] + l[15]*r[3]
	result[4] = l[0]*r[4] + l[4]*r[5] + l[8]*r[6] + l[12]*r[7]
	result[5] = l[1]*r[4] + l[5]*r[5] + l[9]*r[6] + l[13]*r[7]
	result[6] = l[2]*r[4] + l[6]*r[5] + l[10]*r[6] + l[14]*r[7]
	result[7] = l[3]*r[4] + l[7]*r[5] + l[11]*r[6] + l[15]*r[7]
	result[8] = l[0]*r[8] + l[4]*r[9] + l[8]*r[10] + l[12]*r[11]
	result[9] = l[1]*r[8] + l[5]*r[9] + l[9]*r[10] + l[13]*r[11]
	result[10] = l[2]*r[8] + l[6]*r[9] + l[10]*r[10] + l[14]*r[11]
	result[11] = l[3]*r[8] + l[7]*r[9] + l[11]*r[10] + l[15]*r[11]
	result[12] = l[0]*r[12] + l[4]*r[13] + l[8]*r[14] + l[12]*r[15]
	result[13] = l[1]*r[12] + l[5]*r[13] + l[9]*r[14] + l[13]*r[15]
	result[14] = l[2]*r[12] + l[6]*r[13] + l[10]*r[14] + l[14]*r[15]
	result[15] = l[3]*r[12] + l[7]*r[13] + l[11]*r[14] + l[15]*r[15]

	return result
}

// TransVector TransVector
type TransVector [3]float64

// MultipleVector https://www.ldraw.org/article/218.html
func MultipleVector(m *TransMatrix, vs ...*TransVector) []*TransVector {
	results := []*TransVector{}

	for _, v := range vs {
		result := &TransVector{}
		// u' = a*u + b*v + c*w + x
		// v' = d*u + e*v + f*w + y
		// w' = g*u + h*v + i*w + z
		result[0] = m[0]*v[0] + m[4]*v[1] + m[8]*v[2] + m[12]
		result[1] = m[1]*v[0] + m[5]*v[1] + m[9]*v[2] + m[13]
		result[2] = m[2]*v[0] + m[6]*v[1] + m[10]*v[2] + m[14]

		results = append(results, result)
	}

	return results
}

func (tv *TransVector) String() string {
	return fmt.Sprintf("[%7.2f,%7.2f,%7.2f]", tv[0], tv[1], tv[2])
}

// BoundingBox BoundingBox
type BoundingBox struct {
	Min, Max *TransVector
}

// NewBoundingBox NewBoundingBox
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		Min: &TransVector{math.Inf(1), math.Inf(1), math.Inf(1)},
		Max: &TransVector{math.Inf(-1), math.Inf(-1), math.Inf(-1)},
	}
}

// TransEmpty TransEmpty
func (bb *BoundingBox) TransEmpty() *BoundingBox {
	if bb.Min[0] == math.Inf(1) && bb.Min[1] == math.Inf(1) && bb.Min[2] == math.Inf(1) &&
		bb.Max[0] == math.Inf(-1) && bb.Max[0] == math.Inf(-1) && bb.Max[0] == math.Inf(-1) {
		return &BoundingBox{Min: &TransVector{0, 0, 0}, Max: &TransVector{0, 0, 0}}
	}
	return bb
}

// ToGob ToGob
func (bb *BoundingBox) ToGob() [2][3]float64 {
	return [2][3]float64{{bb.Min[0], bb.Min[1], bb.Min[2]}, {bb.Max[0], bb.Max[1], bb.Max[2]}}
}

// MergeMinMaxVector MergeMinMaxVector
func (bb *BoundingBox) MergeMinMaxVector(news ...*TransVector) {
	for _, new := range news {
		if new[0] < bb.Min[0] {
			bb.Min[0] = new[0]
		}
		if new[1] < bb.Min[1] {
			bb.Min[1] = new[1]
		}
		if new[2] < bb.Min[2] {
			bb.Min[2] = new[2]
		}

		if new[0] > bb.Max[0] {
			bb.Max[0] = new[0]
		}
		if new[1] > bb.Max[1] {
			bb.Max[1] = new[1]
		}
		if new[2] > bb.Max[2] {
			bb.Max[2] = new[2]
		}
	}
}

// CalcSize Calc LDU Size
func (bb *BoundingBox) CalcSize() [3]float64 {
	return [3]float64{
		bb.Max[0] - bb.Min[0],
		bb.Max[1] - bb.Min[1],
		bb.Max[2] - bb.Min[2],
	}
}

// CalcBrickSize Calc Brick Size
func (bb *BoundingBox) CalcBrickSize() [3]float64 {
	lduSize := bb.CalcSize()
	return [3]float64{lduSize[0] / 20, lduSize[1] / 20, lduSize[2] / 20}
}

// GetBoxWHTByX GetBoundingBoxWidthHeightTall
func GetBoxWHTByX(b [2][3]float64) (float64, float64, float64) {
	width := b[1][0] - b[0][0]  // x
	height := b[1][2] - b[0][2] // z
	tall := b[1][1] - b[0][1]   // y
	return width, height, tall
}
