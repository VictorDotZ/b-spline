package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/VictorDotZ/bspline/pkg/bspline"
	"github.com/VictorDotZ/bspline/pkg/points"
)

func usage() {
	fmt.Println("usage:")
	flag.PrintDefaults()
	os.Exit(2)
}

var x_0 float64
var x_1 float64
var degree int
var N int
var h float64

func init() {
	flag.Float64Var(&x_0, "x_0", 0.0, "[x_0; x_1]")
	flag.Float64Var(&x_1, "x_1", math.Pi*2, "[x_0; x_1]")
	flag.Float64Var(&h, "h", 1e-2, "t_0 = 0; t_{i+1} = t_i + h; t_n = 1")
	flag.IntVar(&N, "N", 7, "(x_1 - x_0) / (N - 1)")
	flag.IntVar(&degree, "degree", 4, "degree")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	f := func(x float64) float64 {
		return math.Sin(x)
	}

	controlPoints := make([]points.Point2d, N)
	hh := (x_1 - x_0) / float64(N-1)
	for i := 0; i < N; i++ {
		x_i := hh*float64(i) + x_0
		controlPoints[i] = points.Point2d{X: x_i, Y: f(x_i)}
	}

	knots := make([]float64, degree+N+1)

	// left pad
	for i := 0; i <= degree; i++ {
		knots[i] = 0.0
	}
	// right pad
	for i := N; i <= degree+N; i++ {
		knots[i] = 1.0
	}

	// internal knots
	for i := degree + 1; i < N; i++ {
		knots[i] = float64(i-degree) / float64(N-degree+1)
	}

	bspline := bspline.BSpline{
		Degree:        degree,
		ControlPoints: controlPoints,
		Knots:         knots,
		H:             h}

	x, y := bspline.InterpolateOnRange()

	for i := 0; i < len(x); i++ {
		fmt.Printf("%.10f\t%.10f\n", x[i], y[i])
	}
}
