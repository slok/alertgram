# Alertgram [![Build Status][github-actions-image]][github-actions-url] [![Go Report Card][goreport-image]][goreport-url]

Alertgram is the easiest way to forward alerts to [Telegram] (Supports [Prometheus alertmanager] alerts).

<p align="center">
    <img src="https://i.imgur.com/4jdOFj9.jpg" width="40%" align="center" alt="alertgram">
</p>

## Introduction

Everything started as a way of forwarding [Prometheus alertmanager] alerts to [Telegram] because the solutions that I found where too complex, I just wanted to forward alerts to channels without trouble. And Alertgram is just that, a simple app that forwards alerts to Telegram groups and channels.

## Input alerts

Alertgram is developed in a decoupled way so in a future may be extended to more inputs apart from Alertmanager's webhook API (ask for a new input if you want).

## Options

Use `--help` flag to show the options.

The configuration of the app is based on flags that also can be set as env vars prepending `ALERTGRAM` to the var. e.g: the flag `--telegram.api-token` would be `ALERTGRAM_TELEGRAM_API_TOKEN`. You can combine both, flags have preference.

## Run

To forward alerts to Telegram the minimum options that need to be set are `--telegram.api-token` and `--telegram.chat-id`

### Simple example

```bash
docker run -p8080:8080 -p8081:8081 slok/alertgram:latest --telegram.api-token=XXXXX --telegram.chat-id=YYYYY
```

### Production

- [Get telegram API token][telegram-token]
- [Get telegram chat IDs][telegram-chat-id]
- [Configure Alertmanager][alertmanager-configuration]
- [Deploy on Kubernetes][kubernetes-deployment]

## Metrics

The app comes with [Prometheus] metrics, it measures the forwarded alerts, HTTP requests, errors... with rate and latency.

By default are served on `/metrics` on `0.0.0.0:8081`

## Development and debugging

You can use the `--notify.dry-run` to show the alerts on the terminal instead of forwarding them to telegram.

Note that the required options are required, so I would suggest to do this before starting to develop with dry-run mode:

```bash
export ALERTGRAM_TELEGRAM_API_TOKEN=fake
export ALERTGRAM_TELEGRAM_CHAT_ID=1234567890
```

Also remember that you can use `--debug` flag.

## FAQ

### Can I use custom template?

Yes!, use the flag `--notify.template-path`. You can check [testdata/templates](testdata/templates) for examples.

The templates are [HTML Go templates] with [Sprig] functions, so you can use these also.

You can use also the notification dry run mode to check your templates without the need
to notify on telegram:

```bash
export ALERTGRAM_TELEGRAM_API_TOKEN=fake
export ALERTGRAM_TELEGRAM_CHAT_ID=1234567890

go run ./cmd/alertgram/ --notify.template-path=./testdata/templates/simple.tmpl --debug --notify.dry-run
```

To send an alert easily and check the template rendering without an alertmanager, prometheus, alerts... you can use the test alerts that are on [testdata/alerts](testdata/alerts):

```bash
curl -i http://127.0.0.1:8080/alerts -d @./testdata/alerts/base.json
```

[github-actions-image]: https://github.com/slok/alertgram/workflows/CI/badge.svg
[github-actions-url]: https://github.com/slok/alertgram/actions
[goreport-image]: https://goreportcard.com/badge/github.com/slok/alertgram
[goreport-url]: https://goreportcard.com/report/github.com/slok/alertgram
[prometheus alertmanager]: https://github.com/prometheus/alertmanager
[prometheus]: https://prometheus.io/
[telegram]: https://telegram.org/
[telegram-token]: https://core.telegram.org/bots#6-botfather
[telegram-chat-id]: https://github.com/GabrielRF/telegram-id
[alertmanager-configuration]: docs/alertmanager
[kubernetes-deployment]: docs/kubernetes
[html go templates]: https://golang.org/pkg/html/template/
[sprig]: http://masterminds.github.io/sprig
