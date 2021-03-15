# СircleCI Insights API Exporter
Prometheus exporter based on [CircleCI Insights API](https://circleci.com/docs/api/v2/#tag/Insights)

## Build
Build with pre-installed golang:
```sh
$ go build -v
```
Build with Docker:
```sh
$ ./docker_build.sh
```

## Run
Set token via environment variable:
```sh
export CIRCLECI_TOKEN=*************
```

Run exporter with mandatory argument project-slug:
```sh
$ ./circleci-exporter -project-slug bb/foo/bar
```

Check help:
```sh
$ ./circleci-exporter --help
Usage of ./circleci-exporter:
  -project-slug string
        Project slug in the form vcs-slug/org-name/repo-name.
  -vcs-branch string
        VCS branch name. (default "master")
  -web.listen-address string
        Address on which to expose metrics and web interface. (default ":9101")
  -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")
```

Example of metrics:
```
# HELP circleci_duration_max Maximal duration of builds.
# TYPE circleci_duration_max gauge
circleci_duration_max{name="frontend"} 2090
circleci_duration_max{name="backend"} 154
# HELP circleci_duration_mean Mean duration of builds.
# TYPE circleci_duration_mean gauge
circleci_duration_mean{name="frontend"} 1292
circleci_duration_mean{name="backend"} 115
# HELP circleci_duration_median Median duration of builds.
# TYPE circleci_duration_median gauge
circleci_duration_median{name="frontend"} 1292
circleci_duration_median{name="backend"} 114
# HELP circleci_duration_min Minimal duration of builds.
# TYPE circleci_duration_min gauge
circleci_duration_min{name="frontend"} 494
circleci_duration_min{name="backend"} 89
# HELP circleci_duration_p95 95th percentile duration of builds.
# TYPE circleci_duration_p95 gauge
circleci_duration_p95{name="frontend"} 2011
circleci_duration_p95{name="backend"} 135
# HELP circleci_duration_standard_deviation Duration standard deviation of builds.
# TYPE circleci_duration_standard_deviation gauge
circleci_duration_standard_deviation{name="frontend"} 1128
circleci_duration_standard_deviation{name="backend"} 15
# HELP circleci_failed_runs Total number of failed builds.
# TYPE circleci_failed_runs gauge
circleci_failed_runs{name="frontend"} 1
circleci_failed_runs{name="backend"} 4
# HELP circleci_mttr Mean time to recovery.
# TYPE circleci_mttr gauge
circleci_mttr{name="frontend"} 5164
circleci_mttr{name="backend"} 811
# HELP circleci_success_rate Success builds' rate.
# TYPE circleci_success_rate gauge
circleci_success_rate{name="frontend"} 0.5
circleci_success_rate{name="backend"} 0.8
# HELP circleci_successful_runs Total number of successful builds.
# TYPE circleci_successful_runs gauge
circleci_successful_runs{name="frontend"} 1
circleci_successful_runs{name="backend"} 16
# HELP circleci_throughput Builds' throughput metric.
# TYPE circleci_throughput gauge
circleci_throughput{name="frontend"} 2
circleci_throughput{name="backend"} 20
# HELP circleci_total_runs Total number of running builds.
# TYPE circleci_total_runs gauge
circleci_total_runs{name="frontend"} 2
circleci_total_runs{name="backend"} 20
# HELP circleci_up Was the last query of CircleCI successful.
# TYPE circleci_up gauge
circleci_up 1
```
