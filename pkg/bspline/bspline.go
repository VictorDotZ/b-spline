package bspline

import (
	"math"

	"github.com/VictorDotZ/bspline/pkg/points"
)

type BSpline struct {
	ControlPoints []points.Point2d
	Knots         []float64
	Degree        int
	H             float64
}

func (bs *BSpline) InterpolateOnRange() ([]float64, []float64) {
	x := make([]float64, 0)
	y := make([]float64, 0)

	for t := 0.0; t <= 1.0; t += bs.H {
		p := bs.Interpolate(t)
		x = append(x, p.X)
		y = append(y, p.Y)
	}

	return x, y
}

func (bs *BSpline) Interpolate(t float64) points.Point2d {
	var sumPoint points.Point2d = points.Point2d{X: 0.0, Y: 0.0}

	for i, controlPoint := range bs.ControlPoints {
		sumPoint = sumPoint.Sum(controlPoint.Times(bs.BasisFunction(t, i, bs.Degree)))
	}

	return sumPoint
}

func (bs BSpline) BasisFunction(t float64, i, j int) float64 {
	if j == 0 {
		if (bs.Knots[i] <= t) && (t < bs.Knots[i+1]) {
			return 1.0
		} else {
			return 0.0
		}
	}

	var alpha, beta float64 = 0.0, 0.0

	if math.Abs(bs.Knots[i+j]-bs.Knots[i]) > math.Nextafter(1, 2)-1 {
		alpha = (t - bs.Knots[i]) / (bs.Knots[i+j] - bs.Knots[i])
	}

	if math.Abs(bs.Knots[i+j+1]-bs.Knots[i+1]) > math.Nextafter(1, 2)-1 {
		beta = (bs.Knots[i+j+1] - t) / (bs.Knots[i+j+1] - bs.Knots[i+1])
	}

	return alpha*bs.BasisFunction(t, i, j-1) + beta*bs.BasisFunction(t, i+1, j-1)
}
