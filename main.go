package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/Shikugawa/nature-chart-slack/pkg/chart_builder"
	"github.com/Shikugawa/nature-chart-slack/pkg/nature"
	"github.com/Shikugawa/nature-chart-slack/pkg/slack"
	"github.com/Shikugawa/nature-chart-slack/pkg/store"
)

func main() {
	checkInterval := 30                        // minutes
	bufferEntrySize := 60 * 24 / checkInterval // 24h
	postInterval := 1 * time.Hour

	c := nature.NewClient(os.Getenv("NATURE_TOKEN"))
	store, err := store.NewDataStore(int64(bufferEntrySize), "./mem")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := slack.NewSlackClient(os.Getenv("SLACK_TOKEN"))

	timer := time.NewTicker(time.Duration(checkInterval) * time.Minute)
	postTimer := time.NewTicker(postInterval)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

	go func() {
		for {
			select {
			case <-timer.C:
				timestamp := time.Now().Unix()
				resp, err := c.Request()
				if err != nil {
					continue
				}

				if len(resp) == 0 {
					continue
				}

				store.WriteRecord(timestamp, resp[0].NewestEvents.Te.Value)
			case <-postTimer.C:
				records := store.ReadAll()
				sort.Slice(records, func(i, j int) bool {
					return records[i].Timestamp < records[j].Timestamp
				})

				dstPath := "output.png"

				times := make([]int64, len(records))
				points := make([]float64, len(records))

				for i := 0; i < len(records); i++ {
					times[i] = records[i].Timestamp
					points[i] = records[i].Point
				}

				if err := chart_builder.BuildTempChart(times, points, dstPath); err != nil {
					continue
				}

				s.Post(dstPath, "#bot-test")
			}
		}
	}()

	<-sig

	if err := store.Close(); err != nil {
		return
	}
}
