package main

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var (
	RED         = color.RGBA{R: 255, A: 255}
	BLUE        = color.RGBA{B: 255, A: 255}
	GREEN       = color.RGBA{B: 50, G: 190, R: 50, A: 255}
	ORANGE      = color.RGBA{B: 15, G: 175, R: 255, A: 255}
	PINK        = color.RGBA{B: 255, R: 200, A: 255}
	YELLOW      = color.RGBA{G: 255, R: 255, A: 255}
	ROSE        = color.RGBA{B: 200, R: 255, A: 255}
	LIGHT_BLUE  = color.RGBA{B: 255, G: 255, A: 255}
	LIGHT_GREEN = color.RGBA{B: 150, G: 255, R: 150, A: 255}
	BLACK       = color.RGBA{B: 0, G: 0, R: 0, A: 255}
)

var (
	N      int // the maximum number of plots in the same graphics
	colors = []color.Color{RED, BLUE, GREEN, ORANGE, PINK, YELLOW, ROSE, LIGHT_BLUE, LIGHT_GREEN, BLACK}
)

func init() {
	N = len(colors)
}

// Get a color from the pre-defined palette
func getColor(n int) color.Color {
	l := len(colors)
	if N > l {
		colors = append(colors, palette.Reverse(moreland.SmoothBlueRed()).Palette(N+1-l).Colors()...)
	}
	if n >= N {
		return colors[0]
	}
	return colors[n]
}

// Create a plot with title and axis labels
func NewPlot(title, xlabel, ylabel string) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel
	return p, nil
}

type commaTicks struct{}

// Ticks computes the default tick marks, and define the label for the major tick marks.
func (commaTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		tks[i].Label = fmt.Sprintf("%.2f", tks[i].Value)
	}
	return tks
}

// AddWithPoints Draw the data with points
func AddWithPointsXY(x, y []float64, legend string, n int, p *plot.Plot) error {
	points, err := plotter.NewScatter(createPointsXY(x, y))
	if err != nil {
		return err
	}
	points.Radius = 2
	points.Shape = draw.CircleGlyph{}
	points.Color = getColor(n)
	p.Add(points)
	addLegend(legend, p, points, 10)
	p.Y.Tick.Marker = commaTicks{}
	return nil
}

// AddWithLineXY Draw the data with line
func AddWithLineXY(x, y []float64, legend string, n int, p *plot.Plot) error {
	line, err := plotter.NewLine(createPointsXY(x, y))
	if err != nil {
		return err
	}
	line.Color = getColor(n)
	p.Add(line)
	addLegend(legend, p, line, 10)
	p.Y.Tick.Marker = commaTicks{}
	return nil
}

// Add a legend with some position tuned
func addLegend(legend string, p *plot.Plot, thumb plot.Thumbnailer, yoff vg.Length) {
	if legend != "" {
		p.Legend.Add(legend, thumb)
		p.Legend.Padding = -1.
		p.Legend.YOffs = yoff
		p.Legend.YAlign = 0.
		p.Legend.YPosition = -1
	}
	p.Legend.Top = YTOPLEGEND
}

// createPointsXY Transform the x, y slices into plotter
func createPointsXY(x, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	return pts
}
