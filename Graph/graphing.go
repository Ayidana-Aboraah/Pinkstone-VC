package Graph

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"net/http"
)

var Labels *[]string
var Categories *[]string
var Inputs *[]float64

func generateLineItems(r []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(r); i++ {
		items = append(items, opts.LineData{Value: r[i]})
	}
	return items
}

func generatePieItems(data []float64) []opts.LiquidData {
	items := make([]opts.LiquidData, 0)
	for i := 0; i < len(data); {
		items = append(items, opts.LiquidData{Value: data[i]})
	}
	return items
}


func CreateLineGraph(w http.ResponseWriter){
	// create a new line instance
	line := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Bronze Hermes Data",Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
		Title:    "Profit Line Chart",
		Subtitle: "Only profit",
	}),
		)


	line.SetXAxis(Labels)

	for _, v := range *Categories {
			line.AddSeries(v, generateLineItems(*Inputs)).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
			fmt.Println(Inputs)
	}

	// Where the magic happens
	line.Render(w)
}

func CreatePieGraph(w http.ResponseWriter) {
	pie := charts.NewLiquid()

	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Bronze Hermes Data", Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Item Popularity",
			Subtitle: "Total in Numbers area of statistics menu",
		}))

		pie.AddSeries("",generatePieItems(*Inputs))

	pie.Render(w)
}

/*
func CreateLineOverTime(labels, categories []string, inputs []float64){
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Bronze Hermes Data",Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Graph",
			Subtitle: "Data",
		}))

	// Put data into instance
		line.SetXAxis("labels").
			AddSeries("categories[0]", generateLineValues([]float64{0})).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	//line.Render(render)
}

func generateLineValues(inputs []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(inputs); i++ {
		items = append(items, opts.LineData{Value: inputs[i]})
	}
	return items
}

// generate random data for bar chart
func generateRandomBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < len(*Categories); i++ {
		items = append(items, opts.BarData{Value: rand.Intn(999)})
	}
	return items
}
 */