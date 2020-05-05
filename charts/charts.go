package charts

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/SebastianJ/harmony-stats/config"
	"github.com/golang/freetype/truetype"
	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	// Harmony brand style guide: https://harmony.one/brand
	colors = map[string]string{
		"electric_blue": "00AEE9",
		"mint_green":    "69FABD",
		"midnight_blue": "1B295E",
		"cool_gray":     "758796",

		// typography
		"nunito_normal":    "1B295E",
		"fira_sans_normal": "758796",

		// custom
		"light_gray":        "f9f9f9",
		"light_gray_stroke": "eeeeee",
		"mint_green_darker": "56dea5",
	}

	dateFormat string = "2006-01-02"
)

func setupChartPath(fileName string) (string, error) {
	filePath := filepath.Join(config.Configuration.Export.Path, "charts", fileName)
	dirPath, _ := filepath.Split(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", err
	}
	return filePath, nil
}

func loadFont(family string, version string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(filepath.Join(config.Configuration.BasePath, "fonts", family, fmt.Sprintf("%s-%s.ttf", family, version)))
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}

// GenerateBarChart - generates a bar chart based on supplied data
func GenerateBarChart(fileName string, title string, bars []chart.Value) error {
	filePath, err := setupChartPath(fileName)
	if err != nil {
		return err
	}

	style := chart.Style{
		StrokeColor: drawing.ColorFromHex(colors["mint_green_darker"]),
		FillColor:   drawing.ColorFromHex(colors["mint_green"]), //.WithAlpha(80),
		StrokeWidth: 1,
	}

	nunitoBold, err := loadFont("Nunito", "Black")
	if err != nil {
		return err
	}

	firaSansRegular, err := loadFont("FiraSans", "Regular")
	if err != nil {
		return err
	}

	printer := message.NewPrinter(language.English)
	padding := 50
	graph := chart.BarChart{
		Title: title,
		TitleStyle: chart.Style{
			Padding: chart.Box{
				Top: 5,
			},
			Font:      nunitoBold,
			FontColor: drawing.ColorFromHex(colors["electric_blue"]),
		},
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
			FillColor:   drawing.ColorFromHex(colors["light_gray"]),
			StrokeColor: drawing.ColorFromHex(colors["light_gray_stroke"]),
			StrokeWidth: 1,
		},
		YAxis: chart.YAxis{
			Name: "Rewards",
			ValueFormatter: func(v interface{}) string {
				return printer.Sprintf("%d ONE", int(math.RoundToEven(v.(float64))))
			},
			Style: chart.Style{
				Font:      firaSansRegular,
				FontColor: drawing.ColorFromHex(colors["fira_sans_normal"]),
			},
		},
		XAxis: chart.Style{
			Font:     nunitoBold,
			TextWrap: 0,
		},
		BarWidth:   50,
		BarSpacing: 150,
	}

	styledBars := []chart.Value{}
	for _, bar := range bars {
		styledBar := bar
		styledBar.Style = style
		styledBars = append(styledBars, styledBar)
	}
	graph.Bars = styledBars

	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	graph.Render(chart.PNG, file)

	return nil
}

// GenerateTimeSeriesChart - generate a chart for a continous series using supplied data
func GenerateTimeSeriesChart(fileName string, seriesTitle string, xAxisLabel string, yAxisLabel string, xValues []time.Time, yValues []float64, details []string) error {
	filePath, err := setupChartPath(fileName)
	if err != nil {
		return err
	}

	mainSeries := chart.TimeSeries{
		Name: seriesTitle,
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex(colors["electric_blue"]).WithAlpha(5),
			FillColor:   drawing.ColorFromHex(colors["mint_green"]), //.WithAlpha(80),
		},
		XValues: xValues,
		YValues: yValues,
	}

	extraAxisSeries := mainSeries
	extraAxisSeries.YAxis = chart.YAxisSecondary

	xTicks := make([]chart.Tick, len(xValues))
	for index, date := range xValues {
		xTicks[index] = chart.Tick{
			Value: float64(date.UnixNano()),
			Label: date.Format(dateFormat),
		}
	}

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
			Name:  xAxisLabel,
			Ticks: xTicks,
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

	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	graph.Render(chart.PNG, file)

	return nil
}

// GenerateContinousChart - generate a chart for a continous series using supplied data
func GenerateContinousChart(fileName string, seriesTitle string, xAxisLabel string, yAxisLabel string, xValues []float64, yValues []float64, details []string) error {
	filePath, err := setupChartPath(fileName)
	if err != nil {
		return err
	}

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

	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	graph.Render(chart.PNG, file)

	return nil
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
