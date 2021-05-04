// Package report contains all the chart functionality to render visualization.
package report

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/philipglazman/lnviz/data"
	"math"
	"strconv"
	"time"
)

// RouteFeePerChanIn returns a pie chart with the sum of routing fees for routes
// that passed through each inbound channel.
func RouteFeePerChanIn(events data.ForwardEvents, nodes data.Nodes) *charts.Pie {
	bar := charts.NewPie()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Sum of Route Fee by Inbound Channel",
		Subtitle: "Routing",
		Left:     "center",
	}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "item",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
					Lang:  []string{"data view", "turn off", "refresh"},
				},
			}},
		),
	)

	series := make(map[uint64]uint64, 0)
	y := make([]opts.PieData, 0)

	for _, e := range events {
		series[e.ChanIdIn] += e.Fee
	}

	for k, v := range series {
		var nodeAlias string
		if alias, exists := nodes[k]; exists {
			nodeAlias = alias.Alias
		}

		y = append(y, opts.PieData{Name: fmt.Sprintf("%s : %d", nodeAlias, k), Value: v})
	}

	bar.AddSeries("fee", y).SetSeriesOptions(charts.WithLabelOpts(
		opts.Label{
			Show:      false,
			Formatter: "{b}: {c}",
		}),
	)

	return bar
}

// CumulativeRoutingFees returns a line chart with the cumulative fees over time.
func CumulativeRoutingFees(events data.ForwardEvents) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Cumulative Sum of Routing Fees",
		Subtitle: "Routing",
	}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
			Scale:       true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	x := make([]string, 0)
	y := make([]opts.LineData, 0)

	var sumFee uint64
	for _, e := range events {
		x = append(x, time.Unix(int64(e.Timestamp), 0).String())
		sumFee += e.Fee
		y = append(y, opts.LineData{Value: sumFee})
	}

	line.SetXAxis(x).AddSeries("line", y)

	return line
}

// DailyRoutesProcessed returns a line chart with the daily count of routes processed.
func DailyRoutesProcessed(events data.ForwardEvents) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Daily Count of Routes Processed",
		Subtitle: "Routing",
	}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	x := make([]string, 0)        // day
	y := make([]opts.LineData, 0) // # of routes

	// day => # of routes
	series := make(map[string]int)

	for _, e := range events {
		y, m, d := time.Unix(int64(e.Timestamp), 0).Date()
		date := fmt.Sprintf("%d/%02d/%02d", y, m, d)

		series[date]++

		if val := series[date]; val > 0 {
			x = append(x, date)
		}
	}

	for _, date := range x {
		y = append(y, opts.LineData{Value: series[date]})
	}

	line.SetXAxis(x).AddSeries("routing", y).SetSeriesOptions(charts.WithLineChartOpts(
		opts.LineChart{
			Smooth: true,
		}),
	)
	return line
}

// DailyRouteFees returns a line chart with daily route fees.
func DailyRouteFees(events data.ForwardEvents) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Daily Fees from Routes Processed",
		Subtitle: "Routing",
	}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	x := make([]string, 0)        // day
	y := make([]opts.LineData, 0) // # of routes

	// day => # of routes
	series := make(map[string]uint64)

	for _, e := range events {
		y, m, d := time.Unix(int64(e.Timestamp), 0).Date()
		date := fmt.Sprintf("%d/%02d/%02d", y, m, d)

		series[date] += e.Fee

		if val := series[date]; val > 0 {
			x = append(x, date)
		}
	}

	for _, date := range x {
		y = append(y, opts.LineData{Value: series[date]})
	}

	line.SetXAxis(x).AddSeries("routing", y).SetSeriesOptions(charts.WithLineChartOpts(
		opts.LineChart{
			Smooth: true,
		}),
	)
	return line
}

// DailyRoutingVolume returns a line chart with the daily routing volume.
func DailyRoutingVolume(events data.ForwardEvents) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Daily Volume",
		Subtitle: "Routing",
	}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:        "BTC",
			Type:        "",
			Show:        true,
			SplitNumber: 20,
			Scale:       false,
			Min:         nil,
			Max:         nil,
			GridIndex:   0,
			SplitArea:   nil,
			SplitLine:   nil,
			AxisLabel: &opts.AxisLabel{
				Show:      true,
				Inside:    false,
				Formatter: "{value} BTC",
			},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	x := make([]string, 0)        // day
	y := make([]opts.LineData, 0) // # of routes

	// day => # of routes
	series := make(map[string]uint64)

	for _, e := range events {
		y, m, d := time.Unix(int64(e.Timestamp), 0).Date()
		date := fmt.Sprintf("%d/%02d/%02d", y, m, d)

		series[date] += e.AmtIn

		if val := series[date]; val > 0 {
			x = append(x, date)
		}
	}

	for _, date := range x {
		val := strconv.FormatFloat(float64(series[date])/math.Pow10(int(0+8)), 'f', -int(0+8), 64)
		y = append(y, opts.LineData{Value: val})
	}

	line.SetXAxis(x).AddSeries("routing", y).SetSeriesOptions(charts.WithLineChartOpts(
		opts.LineChart{
			Smooth: true,
		}),
	)
	return line
}

// RouteFeePerChanOut returns a pie chart with the route fees per channel out.
func RouteFeePerChanOut(events data.ForwardEvents, nodes data.Nodes) *charts.Pie {
	bar := charts.NewPie()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Routing: Sum of Route Fee by Outbound Channel",
		Subtitle: "Routing",
	}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "item",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
					Lang:  []string{"data view", "turn off", "refresh"},
				},
			}},
		),
	)

	series := make(map[uint64]uint64, 0)
	y := make([]opts.PieData, 0)

	for _, e := range events {
		series[e.ChanIdOut] += e.Fee
	}

	for k, v := range series {
		var nodeAlias string
		if alias, exists := nodes[k]; exists {
			nodeAlias = alias.Alias
		}

		y = append(y, opts.PieData{Name: fmt.Sprintf("%s : %d", nodeAlias, k), Value: v})
	}

	bar.AddSeries("fee", y).SetSeriesOptions(charts.WithLabelOpts(opts.Label{
		Show: false,
	}))

	return bar
}
