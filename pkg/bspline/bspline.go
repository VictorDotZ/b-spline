package bspline

import "fmt"

type BSpline struct {
	Degree  int
	Points  [][]float64
	Knots   []float64
	Weights []float64
	N       int
	D       int
	domain  []int
	low     float64
	high    float64
	v       [][]float64
	H       float64
}

func NewBSpline(degree, n, d int, points [][]float64, weights, knots []float64, h float64) BSpline {
	if degree < 1 || degree > n-1 {
		panic("degree must be at least 1 and less than or equal to point count - 1")
	}

	if len(weights) == 0 {
		weights = make([]float64, n)
		for i := 0; i < n; i++ {
			weights[i] = 1.0
		}
	}

	if len(knots) == 0 {
		knots = make([]float64, n+degree+1)
		for i := 0; i < n+degree+1; i++ {
			knots[i] = float64(i)
		}
	} else {
		if len(knots) != n+degree+1 {
			panic("bad knot vector length")
		}
	}

	fmt.Println(knots)

	domain := make([]int, 2)
	domain[0] = degree
	domain[1] = len(knots) - 1 - degree

	low := float64(knots[domain[0]])
	high := float64(knots[domain[1]])

	return BSpline{
		Degree:  degree,
		Points:  points,
		Knots:   knots,
		Weights: weights,
		N:       n,
		D:       d,
		domain:  domain,
		low:     low,
		high:    high,
		H:       h}
}

func (bs *BSpline) Interpolate() ([]float64, []float64) {
	x := make([]float64, 0)
	y := make([]float64, 0)

	for t := float64(0); t <= 1; t += bs.H {
		// for t := float64(0.5); t < 0.52; t += bs.H {

		pt := bs.interpolate(t)
		// fmt.Println("t:", t, "x:", pt[0])
		x = append(x, pt[0])
		y = append(y, pt[1])
	}

	return x, y
}

func (bs *BSpline) interpolate(t float64) []float64 {
	t = t*(bs.high-bs.low) + bs.low
	// fmt.Println(t)
	if t < bs.low || t > bs.high {
		panic("out of bounds")
	}

	var alpha float64
	s := bs.getSplitSegment(t)
	// fmt.Println("s=", s)

	bs.v = bs.convertPoints()

	for l := 1; l <= bs.Degree+1; l++ {
		for i := s; i > s-bs.Degree-1+l; i-- {
			alpha = (t - bs.Knots[i]) / (bs.Knots[i+bs.Degree+1-l] - bs.Knots[i])

			// fmt.Println("alpha =", alpha)

			for j := 0; j < bs.D+1; j++ {
				bs.v[i][j] = (1-alpha)*bs.v[i-1][j] + alpha*bs.v[i][j]

				// fmt.Println("bs.v[i*bs.D+j] =", bs.v[i*bs.D+j])

			}
		}
	}

	result := make([]float64, bs.D)
	for i := 0; i < bs.D; i++ {
		result[i] = bs.v[s][i] / bs.v[s][bs.D]
	}

	return result
}

func (bs BSpline) getSplitSegment(t float64) int {
	var s int

	for s = bs.domain[0]; s < bs.domain[1]; s++ {
		if t >= float64(bs.Knots[s]) && t <= float64(bs.Knots[s+1]) {
			break
		}
	}

	return s
}

func (bs BSpline) convertPoints() [][]float64 {
	v := make([][]float64, bs.N)
	for i := 0; i < bs.N; i++ {
		v[i] = make([]float64, bs.D+1)
		for j := 0; j < bs.D; j++ {
			v[i][j] = bs.Points[i][j] * bs.Weights[i]
		}
		v[i][bs.D] = bs.Weights[i]
	}

	return v
}
