package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/grafana/explore-logs/generator/log"
	"github.com/grafana/loki-client-go/loki"
	"github.com/grafana/loki/pkg/push"
	"github.com/prometheus/common/model"
)

func main() {
	url := flag.String("url", "http://localhost:3100/loki/api/v1/push", "Loki URL")
	tenantId := flag.String("tenant-id", "", "Loki tenant ID")
	flag.Parse()

	cfg, err := loki.NewDefaultConfig(*url)
	if err != nil {
		panic(err)
	}
	cfg.BackoffConfig.MaxRetries = 1
	cfg.BackoffConfig.MinBackoff = 100 * time.Millisecond
	cfg.BackoffConfig.MaxBackoff = 100 * time.Millisecond

	if *tenantId != "" {
		cfg.TenantID = *tenantId
	}

	client, err := loki.New(cfg)
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Creates and starts all apps.
	for namespace, apps := range generators {
		for serviceName, generator := range apps {
			log.ForAllClusters(namespace, serviceName, func(labels model.LabelSet, metadata push.LabelsAdapter) {
				generator(ctx, log.NewAppLogger(labels, log.NewOtelLogger(string(serviceName), labels)), metadata)
			})
		}
	}

	<-ctx.Done()
}
