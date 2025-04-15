package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/grafana/explore-logs/generator/flog"
	"github.com/grafana/explore-logs/generator/log"
	"github.com/grafana/loki/pkg/push"
	"github.com/prometheus/common/model"
)

type LogGenerator func(ctx context.Context, logger *log.AppLogger, metadata push.LabelsAdapter)

var generators = map[model.LabelValue]map[model.LabelValue]LogGenerator{
	"gateway": {
		"grafanacon-json": func(ctx context.Context, logger *log.AppLogger, metadata push.LabelsAdapter) {
			go func() {
				for ctx.Err() == nil {
					level := log.RandLevel()
					t := time.Now()
					logger.LogWithMetadata(level, t, flog.NewJSONLogFormat(t, log.RandURI(), statusFromLevel(level)), metadata)
					time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
				}
			}()
		},
		"grafanacon-otel": func(ctx context.Context, logger *log.AppLogger, metadata push.LabelsAdapter) {
			go func() {
				for ctx.Err() == nil {
					level := log.RandLevel()
					t := time.Now()
					logger.LogWithMetadata(level, t, flog.NewJSONLogFormat(t, log.RandURI(), statusFromLevel(level)), metadata)
					time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
				}
			}()
		},
	},
}

func statusFromLevel(level model.LabelValue) int {
	switch level {
	case log.INFO:
		return 200
	case log.WARN:
		return 400
	case log.ERROR:
		return 500
	default:
		return 200
	}
}
