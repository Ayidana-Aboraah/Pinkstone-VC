package Graph

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"net/http"
)

var Labels *[]string
var Categories *[]string
var Inputs *[]float64
var LineInputs *[][]float64

func generateLineItems(label string, r []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i, _ := range r{
		items = append(items, opts.LineData{Name:label, Value: r[i]})
	}
	return items
}

func generatePieItems(tags []string, data []float64) []opts.PieData {
	items := make([]opts.PieData, 0)


	for i, _ := range tags {
		items = append(items, opts.PieData{Name: tags[i], Value: data[i]})
	}

	return items
}

func CreateLineGraph(w http.ResponseWriter) {
	line := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Bronze Hermes Data", Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Profit Line Chart",
			Subtitle: "Dark blue is Revenue, Light blue is Cost, Profit is pink;",
		}),
	)

	line.SetXAxis(Labels)

	for _, v := range *LineInputs{
		line.AddSeries("", generateLineItems("data", v)).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	}

	line.Render(w)
}

func CreatePieGraph(w http.ResponseWriter) {
	pie := charts.NewPie()

	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Bronze Hermes Data", Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Item Popularity",
			Subtitle: "Hover Over them to see how much they take up the wheel",
		}))

	pie.AddSeries("Tree", generatePieItems(*Labels, *Inputs))

	pie.Render(w)
}