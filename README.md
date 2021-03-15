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
