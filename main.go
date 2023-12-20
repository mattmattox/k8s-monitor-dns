package main

import (
	"flag"
	"os/exec"
	"strings"
	"time"

	"github.com/mattmattox/k8s-monitor-dns/pkg/config"
	"github.com/mattmattox/k8s-monitor-dns/pkg/logging"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logging.SetupLogging()

var (
	dnsStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dns_status",
			Help: "DNS up/down status. 1 for up, 0 for down.",
		},
		[]string{"type"}, // type can be "internal" or "external"
	)
	dnsResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_response_time_seconds",
			Help:    "Histogram of response times for DNS checks in seconds.",
			Buckets: prometheus.DefBuckets, // Default buckets, adjust as needed
		},
		[]string{"type"}, // type can be "internal" or "external"
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(dnsStatus)
	prometheus.MustRegister(dnsResponseTime)
}

func main() {
	flag.Parse() // Parse command-line flags

	config := config.LoadConfigFromEnv()

	logging.SetupLogging()

	if config.Debug {
		log.Debug("Debug mode enabled")
	} else {
		log.Info("Debug mode disabled")
	}

	log.Info("Starting k8s-monitor-dns")

	// Set up HTTP server for Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		MetircsPort := ":" + config.MetircsPort
		if err := http.ListenAndServe(MetircsPort, nil); err != nil {
			log.Errorf("Error starting Prometheus HTTP server: %v", err)
		}
	}()

	// Parse timeout and delay
	timeout := time.Duration(config.Timeout) * time.Second
	delay := time.Duration(config.Delay) * time.Second

	for {

		log.Infof("Checking Internal DNS: %s", config.InternalHost)
		checkDNS(config.InternalIP, config.InternalHost, config.InternalIP, timeout, "internal")

		log.Infof("Checking External DNS: %s", config.ExternalHost)
		checkDNS(config.ExternalIP, config.ExternalHost, config.ExternalIP, timeout, "external")

		time.Sleep(delay)
	}

}

func checkDNS(ip, host, expectedIP string, timeout time.Duration, dnsType string) {
	start := time.Now() // Start time for measuring response time

	// Execute the DNS check
	cmd := exec.Command("timeout", timeout.String(), "dig", "+short", "@"+ip, host)
	output, err := cmd.Output()

	elapsed := time.Since(start) // Calculate the elapsed time

	// Update Prometheus metrics for response time
	dnsResponseTime.WithLabelValues(dnsType).Observe(elapsed.Seconds())

	if err != nil {
		// DNS check failed (down)
		dnsStatus.WithLabelValues(dnsType).Set(0) // DNS is down
		log.Errorf("DNS check for %s failed: %v", host, err)
		return
	}

	// Trim the output and compare with expected IP
	actualIP := strings.TrimSpace(string(output))
	if actualIP == expectedIP {
		// DNS check succeeded (up)
		dnsStatus.WithLabelValues(dnsType).Set(1) // DNS is up
		log.Infof("%s DNS is OK", host)
	} else {
		// DNS check returned an unexpected result (down)
		dnsStatus.WithLabelValues(dnsType).Set(0) // DNS is down
		log.Errorf("%s DNS returned a bad IP: %s", host, actualIP)
	}
}
