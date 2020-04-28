package charts

import (
	"fmt"
	"os"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

var (
	// Harmony brand style guide: https://harmony.one/brand
	colors = map[string]string{
		"electric_blue": "00AEE9",
		"mint_green":    "69FABD",
		"midnight_blue": "1B295E",
		"cool_gray":     "758796",

		// custom
		"light_gray": "f9f9f9",
	}
)

// GenerateGraph - generate a graph using supplied data
func GenerateGraph(fileName string, seriesTitle string, xAxisLabel string, yAxisLabel string, xValues []float64, yValues []float64, details []string) {
	mainSeries := chart.ContinuousSeries{
		Name: seriesTitle,
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex(colors["electric_blue"]),
			FillColor:   drawing.ColorFromHex(colors["mint_green"]), //.WithAlpha(80),
		},
		XValues: xValues,
		YValues: yValues,
	}

	extraAxisSeries := mainSeries
	extraAxisSeries.YAxis = chart.YAxisSecondary

	padding := 50
	graph := chart.Chart{
		Width:  1920,
		Height: 1080,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    padding,
				Bottom: padding,
				Left:   padding,
				Right:  padding,
			},
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex(colors["light_gray"]),
		},
		YAxis: chart.YAxis{
			Name: yAxisLabel,
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
		},
		YAxisSecondary: chart.YAxis{
			Name: yAxisLabel,
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
		},
		XAxis: chart.XAxis{
			Name: xAxisLabel,
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
		},
		Series: []chart.Series{
			mainSeries,
			extraAxisSeries,
		},
	}

	//graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	detailsStyle := chart.Style{
		FillColor:   drawing.ColorFromHex(colors["electric_blue"]),
		FontColor:   drawing.ColorFromHex(colors["mint_green"]),
		FontSize:    11.0,
		StrokeColor: drawing.ColorFromHex(colors["electric_blue"]),
		StrokeWidth: chart.DefaultAxisLineWidth,
	}

	if len(details) > 0 {
		graph.Elements = []chart.Renderable{DetailsBox(&graph, details, detailsStyle)}
	}

	f, _ := os.Create(fileName)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

// DetailsBox adds a box with additional text
func DetailsBox(c *chart.Chart, text []string, userDefaults ...chart.Style) chart.Renderable {
	return func(r chart.Renderer, box chart.Box, chartDefaults chart.Style) {
		// default style
		defaults := chart.Style{
			FillColor:   drawing.ColorWhite,
			FontColor:   chart.DefaultTextColor,
			FontSize:    8.0,
			StrokeColor: chart.DefaultAxisColor,
			StrokeWidth: chart.DefaultAxisLineWidth,
		}

		var style chart.Style
		if len(userDefaults) > 0 {
			style = userDefaults[0].InheritFrom(chartDefaults.InheritFrom(defaults))
		} else {
			style = chartDefaults.InheritFrom(defaults)
		}

		yPadding := 15
		xPadding := 20

		contentPadding := chart.Box{
			Top:    box.Height() + yPadding,
			Left:   xPadding,
			Right:  xPadding,
			Bottom: box.Height() + yPadding,
		}

		contentBox := chart.Box{
			Top:  100,
			Left: 100,
		}

		content := chart.Box{
			Top:    contentBox.Top - yPadding,
			Left:   contentBox.Left + contentPadding.Left,
			Right:  contentBox.Left + contentPadding.Left,
			Bottom: contentBox.Top,
		}

		style.GetTextOptions().WriteToRenderer(r)

		// measure and add size of text to box height and width
		for _, t := range text {
			textbox := r.MeasureText(t)
			content.Top -= textbox.Height()
			right := content.Left + textbox.Width()
			content.Right = chart.MaxInt(content.Right, right)
		}

		contentBox = contentBox.Grow(content)
		contentBox.Right = content.Right + contentPadding.Right
		contentBox.Top = content.Top - yPadding

		// draw the box
		chart.Draw.Box(r, contentBox, style)

		style.GetTextOptions().WriteToRenderer(r)

		// add the text
		ycursor := content.Top
		x := content.Left
		for _, t := range text {
			textbox := r.MeasureText(t)
			y := ycursor + textbox.Height()
			r.Text(t, x, y)
			ycursor += textbox.Height()
		}
	}
}
