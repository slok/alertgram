# Alertgram [![Build Status][github-actions-image]][github-actions-url] [![Go Report Card][goreport-image]][goreport-url]

Alertgram is the easiest way to forward alerts to [Telegram] (Supports [Prometheus alertmanager] alerts).

<p align="center">
    <img src="https://i.imgur.com/4jdOFj9.jpg" width="40%" align="center" alt="alertgram">
</p>

## Table of contents

- [Introduction](#introduction)
- [Features](#features)
- [Input alerts](#input-alerts)
- [Options](#options)
- [Run](#run)
  - [Simple example](#simple-example)
  - [Production](#production)
- [Metrics](#metrics)
- [Development and debugging](#development-and-debugging)
- [FAQ](#faq)
  - [Only alertmanager alerts are supported?](#only-alertmanager-alerts-are-supported)
  - [Where does alertgram listen to alertmanager alerts?](#where-does-alertgram-listen-to-alertmanager-alerts)
  - [Can I notify to different chats?](#can-i-notify-to-different-chats)
  - [Can I use custom templates?](#can-i-use-custom-templates)
  - [Dead man's switch?](#dead-mans-switch)

## Introduction

Everything started as a way of forwarding [Prometheus alertmanager] alerts to [Telegram] because the solutions that I found were too complex, I just wanted to forward alerts to channels without trouble. And Alertgram is just that, a simple app that forwards alerts to Telegram groups and channels and some small features that help like metrics and dead man's switch.

## Features

- Alertmanager alerts webhook receiver compatibility.
- Telegram notifications.
- Metrics in Prometheus format.
- Optional dead man switch endpoint.
- Optional customizable templates.
- Configurable notification chat ID targets (with fallback to default chat ID).
- Easy setup and flexible.
- Lightweight.
- Perfect for any environment, from a company cluster to home cheap clusters (e.g [K3S]).

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

### Only alertmanager alerts are supported?

At this moment yes, but we can add more input alert systems if you want, create an issue
so we can discuss and implement.

### Where does alertgram listen to alertmanager alerts?

By default in `0.0.0.0:8080/alerts`, but you can use `--alertmanager.listen-address` and
`--alertmanager.webhook-path` to customize.

### Can I notify to different chats?

There are 3 levels where you could customize the notification chat:

- By default: Using the required `--telegram.chat-id` flag.
- At URL level: using [query string] parameter, e.g. `0.0.0.0:8080/alerts?chat-id=-1009876543210`.
  This query param can be customized with `--alertmanager.chat-id-query-string` flag.
- At alert level: TODO

The preference is in order from the lowest to the highest: Default, URL, Alert.

### Can I use custom templates?

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

### Dead man's switch?

A [dead man's switch][dms] (from now on, DMS) is a technique or process where at regular intervals a signal must be received
so the DMS is disabled, if that signal is not received it will be activated.

In monitoring this would be: If an alert is not received at regular intervals, the switch will be activated and notify that we
are not receiving alerts, this is mostly used to know that our alerting system is working.

For example we would set Prometheus triggering an alert continously, Alertmanager sending this specific alert
every `7m` to the DMS endpoint in Alertgram, and Alertgram would be configured with a `10m` interval DMS.

With this setup if Prometheus fails creating the alert, Alertmanager sending the alert to Alertgram, or Alertgram not receiving
this alert (e.g. network problems), Alertmanager will send an alert to Telegram to notify us that our monitoring system is broken.

You could use the same alertgram or another instance, usually in other machine, cluster... so if the cluster/machine fails, your
is isolated and could notify you.

To Enable Alertgram's DMS use `--dead-mans-switch.enable` to enable. By default it will be listening in `/alert/dms`, with a
`15m` interval and use the telegrams default notifier and chat ID. To customize this settings use:

- `--dead-mans-switch.interval`: To configure the interval.
- `--dead-mans-switch.chat-id`: To configure the notifier chat, is independent of the notifier
  although at this moment is Telegram, if not set it will use the notifier default chat target.
- `--alertmanager.dead-mans-switch-path` To configure the path the alertmanager can send the DMS alerts.

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
[query string]: https://en.wikipedia.org/wiki/Query_string
[k3s]: https://k3s.io/
[dms]: https://en.wikipedia.org/wiki/Dead_man%27s_switch
