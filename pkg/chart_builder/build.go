package chart_builder

import (
	"os"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
)

func BuildTempChart(timestamps []int64, temps []float64, dstPath string) error {
	times := make([]time.Time, len(timestamps))

	for i := 0; i < len(timestamps); i++ {
		times[i] = time.Unix(timestamps[i], 0)
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Tempature (℃)",
		},

		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: times,
				YValues: temps,
			},
		},
	}

	f, _ := os.Create(dstPath)
	defer f.Close()

	err := graph.Render(chart.PNG, f)
	if err != nil {
		return err
	}

	return nil
}
