package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace  = "circleci"
	insightURL = "https://circleci.com/api/v2/insights"
)

var (
	listenAddress = flag.String("web.listen-address", ":9101",
		"Address on which to expose metrics and web interface.")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics.")
	projectSlug = flag.String("project-slug", "",
		"Project slug in the form vcs-slug/org-name/repo-name.")
	vcsBranch = flag.String("vcs-branch", "master",
		"VCS branch name.")

	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of CircleCI successful.",
		nil, nil,
	)
	successRate = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "success_rate"),
		"Success builds' rate.",
		[]string{"name"}, nil,
	)
	totalRuns = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "total_runs"),
		"Total number of running builds.",
		[]string{"name"}, nil,
	)
	failedRuns = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "failed_runs"),
		"Total number of failed builds.",
		[]string{"name"}, nil,
	)
	successfulRuns = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "successful_runs"),
		"Total number of successful builds.",
		[]string{"name"}, nil,
	)
	throughput = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "throughput"),
		"Builds' throughput metric.",
		[]string{"name"}, nil,
	)
	mttr = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mttr"),
		"Mean time to recovery.",
		[]string{"name"}, nil,
	)
	durationMin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "duration_min"),
		"Minimal duration of builds.",
		[]string{"name"}, nil,
	)
	durationMax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "duration_max"),
		"Maximal duration of builds.",
		[]string{"name"}, nil,
	)
	durationMedian = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "duration_median"),
		"Median duration of builds.",
		[]string{"name"}, nil,
	)
	durationMean = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "duration_mean"),
		"Mean duration of builds.",
		[]string{"name"}, nil,
	)
	durationP95 = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "duration_p95"),
		"95th percentile duration of builds.",
		[]string{"name"}, nil,
	)
	durationStandardDeviation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "duration_standard_deviation"),
		"Duration standard deviation of builds.",
		[]string{"name"}, nil,
	)
)

// Response https://circleci.com/docs/api/v2/#operation/getProjectWorkflowMetrics
type Response struct {
	Items         []Item `json:"items"`
	NextPageToken string `json:"next_page_token"`
}

// Item structure
type Item struct {
	Name        string `json:"name"`
	WindowStart string `json:"window_start"`
	WindowEnd   string `json:"window_end"`
	Metrics     Metric `json:"metrics"`
}

// Metric structure
type Metric struct {
	SuccessRate      float64        `json:"success_rate"`
	TotalRuns        int64          `json:"total_runs"`
	FailedRuns       int64          `json:"failed_runs"`
	SuccessfulRuns   int64          `json:"successful_runs"`
	Throughput       float64        `json:"throughput"`
	MTTR             int64          `json:"mttr"`
	TotalCreditsUsed int64          `json:"total_credits_used"`
	DurationMetrics  DurationMetric `json:"duration_metrics"`
}

// DurationMetric structure
type DurationMetric struct {
	Min               int64   `json:"min"`
	Max               int64   `json:"max"`
	Median            int64   `json:"median"`
	Mean              int64   `json:"mean"`
	P95               int64   `json:"p95"`
	StandardDeviation float64 `json:"standard_deviation"`
}

// Exporter structure
type Exporter struct {
	insightURL string
	slug       string
	token      string
	vcsBranch  string
}

// NewExporter func
func NewExporter(insightURL string, slug string, token string, vcsBranch string) *Exporter {
	return &Exporter{
		insightURL: insightURL,
		slug:       slug,
		token:      token,
		vcsBranch:  vcsBranch,
	}
}

// Describe func
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- successRate
	ch <- totalRuns
	ch <- failedRuns
	ch <- successfulRuns
	ch <- throughput
	ch <- mttr
	ch <- durationMin
	ch <- durationMax
	ch <- durationMean
	ch <- durationMedian
	ch <- durationP95
	ch <- durationStandardDeviation
}

// LoadMetrics func
func (e *Exporter) LoadMetrics() (body []byte, err error) {
	url := e.insightURL + "/" + e.slug + "/workflows?circle-token=" + e.token + "&reporting-window=last-24-hours" + "&branch=" + e.vcsBranch

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ = ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return body, nil
}

// Collect func
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	resp, err := e.LoadMetrics()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	if err != nil {
		log.Println(err)
	}

	var buildMetrics Response
	err = json.Unmarshal(resp, &buildMetrics)
	if err != nil {
		log.Println(err)
	}
	for i := range buildMetrics.Items {
		name := buildMetrics.Items[i].Name
		ch <- prometheus.MustNewConstMetric(
			successRate, prometheus.GaugeValue, buildMetrics.Items[i].Metrics.SuccessRate, name,
		)
		ch <- prometheus.MustNewConstMetric(
			totalRuns, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.TotalRuns), name,
		)
		ch <- prometheus.MustNewConstMetric(
			failedRuns, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.FailedRuns), name,
		)
		ch <- prometheus.MustNewConstMetric(
			successfulRuns, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.SuccessfulRuns), name,
		)
		ch <- prometheus.MustNewConstMetric(
			throughput, prometheus.GaugeValue, buildMetrics.Items[i].Metrics.Throughput, name,
		)
		ch <- prometheus.MustNewConstMetric(
			mttr, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.MTTR), name,
		)
		ch <- prometheus.MustNewConstMetric(
			durationMin, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.DurationMetrics.Min), name,
		)
		ch <- prometheus.MustNewConstMetric(
			durationMax, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.DurationMetrics.Max), name,
		)
		ch <- prometheus.MustNewConstMetric(
			durationMean, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.DurationMetrics.Mean), name,
		)
		ch <- prometheus.MustNewConstMetric(
			durationMedian, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.DurationMetrics.Median), name,
		)
		ch <- prometheus.MustNewConstMetric(
			durationP95, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.DurationMetrics.P95), name,
		)
		ch <- prometheus.MustNewConstMetric(
			durationStandardDeviation, prometheus.GaugeValue, float64(buildMetrics.Items[i].Metrics.DurationMetrics.StandardDeviation), name,
		)
	}
}

func main() {
	flag.Parse()

	token := os.Getenv("CIRCLECI_TOKEN")

	if *projectSlug == "" {
		log.Fatal("[FATAL] Argument project-slug is empty")
	}

	log.Printf("[INFO] Starting CircleCI Insights Exporter on port " + *listenAddress)
	exporter := NewExporter(insightURL, *projectSlug, token, *vcsBranch)
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>CircleCI Insights Exporter</title></head>
             <body>
             <h1>CircleCI Insights Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>	
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
