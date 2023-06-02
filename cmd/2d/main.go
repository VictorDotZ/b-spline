package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
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
var y_0 float64
var y_1 float64
var degree int
var N int
var h float64

func init() {
	flag.Float64Var(&x_0, "x_0", 0.0, "[x_0; x_1]")
	flag.Float64Var(&x_1, "x_1", math.Pi*2, "[x_0; x_1]")
	flag.Float64Var(&y_0, "y_0", 0.0, "[y_0; y_1]")
	flag.Float64Var(&y_1, "y_1", math.Pi*2, "[y_0; y_1]")
	flag.Float64Var(&h, "h", 1e-2, "t_0 = 0; t_{i+1} = t_i + h; t_n = 1")
	flag.IntVar(&N, "N", 7, "(x_1 - x_0) / (N - 1)")
	flag.IntVar(&degree, "degree", 4, "degree")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	file, err := os.Create("ground_truth.txt")

	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	f := func(x, y float64) float64 {
		return math.Sin(x) * math.Cos(y)
	}

	xScaler := func(x float64) float64 {
		return (x - x_0) / (x_1 - x_0)
	}

	yScaler := func(y float64) float64 {
		return (y - y_0) / (y_1 - y_0)
	}

	rand.Seed(42)

	mask := make([][]bool, N)
	interpolatedGrid := make([][]float64, N)

	hh := (x_1 - x_0) / float64(N-1)
	for i := 0; i < N; i++ {
		mask[i] = make([]bool, N)
		interpolatedGrid[i] = make([]float64, N)

		x_i := hh*float64(i) + x_0

		for j := 0; j < N; j++ {
			y_j := hh*float64(j) + y_0

			if rand.Float64() > 0.8 {
				mask[i][j] = true
			} else {
				mask[i][j] = false
			}

			_, err := fmt.Fprintf(file, "%.10f\t%.10f\t%.10f\n", x_i, y_j, f(x_i, y_j))
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	numNonMaskedJ := 0
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if !mask[i][j] {
				numNonMaskedJ++
			}
		}
	}

	numNonMaskedI := 0
	for j := 0; j < N; j++ {
		for i := 0; i < N; i++ {
			if !mask[i][j] {
				numNonMaskedI++
			}
		}
	}

	numNonMasked := int(math.Min(float64(numNonMaskedI), float64(numNonMaskedJ)))
	if degree > numNonMasked {
		fmt.Println("change degree from", degree, "to", numNonMasked-1)
		degree = numNonMasked - 1
	}

	for i := 0; i < N; i++ {
		x_i := hh*float64(i) + x_0
		controlPoints := make([]points.Point2d, 0)
		for j := 0; j < N; j++ {
			if !mask[i][j] {
				y_j := hh*float64(j) + y_0

				controlPoints = append(controlPoints, points.Point2d{X: y_j, Y: f(x_i, y_j)})
			}
		}

		n := len(controlPoints)

		knots := make([]float64, degree+n+1)

		for k := 0; k <= degree; k++ {
			knots[k] = 0.0
		}

		for k := n; k <= degree+n; k++ {
			knots[k] = 1.0
		}

		for k := degree + 1; k < n; k++ {
			knots[k] = float64(k-degree) / float64(n-degree+1)
		}

		bspline := bspline.BSpline{
			Degree:        degree,
			ControlPoints: controlPoints,
			Knots:         knots,
			H:             h}

		for j := 0; j < N; j++ {
			y_j := hh*float64(j) + y_0
			interpolatedGrid[i][j] = bspline.Interpolate(yScaler(y_j)).Y
		}
	}

	fileInt, err := os.Create("interpolated.txt")

	if err != nil {
		log.Fatalln(err)
	}

	defer fileInt.Close()

	for j := 0; j < N; j++ {
		y_j := hh*float64(j) + y_0
		controlPoints := make([]points.Point2d, 0)

		for i := 0; i < N; i++ {
			x_i := hh*float64(i) + x_0
			if !mask[i][j] {
				controlPoints = append(controlPoints, points.Point2d{X: x_i, Y: f(x_i, y_j)})
			}
		}
		n := len(controlPoints)

		knots := make([]float64, degree+n+1)

		for k := 0; k <= degree; k++ {
			knots[k] = 0.0
		}

		for k := n; k <= degree+n; k++ {
			knots[k] = 1.0
		}

		for k := degree + 1; k < n; k++ {
			knots[k] = float64(k-degree) / float64(n-degree+1)
		}

		bspline := bspline.BSpline{
			Degree:        degree,
			ControlPoints: controlPoints,
			Knots:         knots,
			H:             h}
		for i := 0; i < N; i++ {
			x_i := hh*float64(i) + x_0

			interpolatedGrid[i][j] += bspline.Interpolate(xScaler(x_i)).Y
			interpolatedGrid[i][j] *= 0.5

			_, err := fmt.Fprintf(fileInt, "%.10f\t%.10f\t%.10f\n", x_i, y_j, interpolatedGrid[i][j])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
