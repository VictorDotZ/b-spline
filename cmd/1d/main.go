package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	bspline "github.com/VictorDotZ/bspline/pkg/bspline"
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
	flag.Float64Var(&x_1, "x_1", math.Pi*2.0, "[x_0; x_1]")
	flag.Float64Var(&h, "h", 1e-2, "scaled step")
	flag.IntVar(&N, "N", 4, "(x_1 - x_0) / N")
	flag.IntVar(&degree, "degree", 3, "(x_1 - x_0) / N")

}

func main() {
	flag.Usage = usage
	flag.Parse()

	D := 2

	f := func(x float64) float64 {
		return math.Sin(x)
	}

	points := make([][]float64, N)
	hh := (x_1 - x_0) / float64(N-1)
	for i := 0; i < N; i++ {
		x_i := hh*float64(i) + x_0
		points[i] = []float64{x_i, f(x_i)}
	}

	// for i := 0; i < N; i++ {
	// fmt.Printf("%.10f\t%.10f\n", points[i*D+0], points[i*D+1])
	// }

	knots := make([]float64, N+degree+1)
	for i := 0; i < degree+1; i++ {
		knots[i] = float64(0)
	}

	for i := degree + 1; i < N; i++ {
		knots[i] = float64(i - degree)
	}

	for i := N; i < degree+N+1; i++ {
		knots[i] = knots[N] + 1
	}

	fmt.Println(knots)

	// for i := 0; i < n+degree+1; i++ {
	// 	knots[i] = float64(i)
	// }

	bspline := bspline.NewBSpline(degree, N, D, points, make([]float64, 0), make([]float64, 0), h)

	x, y := bspline.Interpolate()

	// fmt.Println(len(x))

	for i := 0; i < len(x); i++ {
		fmt.Printf("%.10f\t%.10f\n", x[i], y[i])
	}
}
