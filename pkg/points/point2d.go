package points

type Point2d struct {
	X float64
	Y float64
}

func (p *Point2d) Sum(pt Point2d) Point2d {
	return Point2d{p.X + pt.X, p.Y + pt.Y}
}

func (p *Point2d) Times(x float64) Point2d {
	return Point2d{p.X * x, p.Y * x}
}

func (p *Point2d) Add(x float64) Point2d {
	return Point2d{p.X + x, p.Y + x}
}
