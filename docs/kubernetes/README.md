# Kubernetes deployment

## Alertgram

First of all check you meet the requirements:

- [Telegram bot token.][telegram-token]
- [Telegram channel or group ID.][telegram-chat-id]
- Add your Telegram bot to the Telegram channel/group.
- Write a message to the bot at least once (not sure if required currently).

The manifest file is in [docs/kubernetes/deploy.yaml](deploy.yaml).

before aplying the manifests check these:

- They will be deployed on `monitoring`
- It will set a ServiceMonitor for the alertgram metrics, the prometheus that
  will scrape the metrics is `prometheus`, change the labels if required.
- Change the string `CHANGE_ME_TELEGRAM_API_TOKEN`for your Telegram API token.
- Change the string `CHANGE_ME_TELEGRAM_CHAT_ID`for your Telegram chat ID.
- Check the alertgram container image version and update to the one you want `image: slok/alertgram:XXXXXXXX`.
- Alertgram will be listening on `http://alertgram:8080/alerts` or `http://alertgram.monitoring.svc.cluster.local:8080/alerts`.

Now you are ready to `kubectly apply -f ./deploy.yaml`.

## Alertmanager

As an example in [docs/kubernetes/alertmanager-cfg.yaml](alertmanager-cfg.yaml),
is an alertmanager configuration example that connects with alertgram.

Before deploying this, we assume:

- You are using prometheus-operator.
- Your alertmanager is `main`, if not change accordingly (labels, names...)
- Is deployed on `monitoring` namespace.

[telegram-token]: https://core.telegram.org/bots#6-botfather
[telegram-chat-id]: https://github.com/GabrielRF/telegram-id
