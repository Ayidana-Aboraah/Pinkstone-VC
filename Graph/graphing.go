package Graph

import (
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

var Labels []string
var Inputs []float32
var LineInputs [][]float32

func generateLineItems(label string, r []float32) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := range r {
		items = append(items, opts.LineData{Name: label, Value: r[i]})
	}
	return items
}

func generatePieItems(tags []string, data []float32) []opts.PieData {
	items := make([]opts.PieData, 0)

	for i := range tags {
		items = append(items, opts.PieData{Name: tags[i], Value: data[i]})
	}

	return items
}

func CreateLineGraph(w http.ResponseWriter) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Bronze Hermes Data", Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line Chart",
			Subtitle: "The line labeled Series is the total of everything on the chart.",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Save as image",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
					Lang:  []string{"Number view", "turn off", "refresh"},
				}},
		}),
		charts.WithLegendOpts(opts.Legend{Bottom: "0%", Show: true, SelectedMode: "multiple", Orient: "horizontal"}),
	)

	line.SetXAxis([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14",
	"15", "16", "17", "18", "19", "20", "21", "22", "23", "24",
	"25", "26", "27", "28", "29", "30", "31"})

	for i, v := range LineInputs {
		line.AddSeries(Labels[i], generateLineItems(Labels[i], v)).SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
			charts.WithMarkPointNameTypeItemOpts(
				opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
				opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
				opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
			),
			charts.WithMarkPointStyleOpts(
				opts.MarkPointStyle{Label: &opts.Label{Show: true}}),
		)
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
		}),
	)

	pie.SetSeriesOptions(charts.WithPieChartOpts(opts.PieChart{Radius: 50}))

	pie.AddSeries("Tree", generatePieItems(Labels, Inputs)).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
		)

	pie.Render(w)
}
