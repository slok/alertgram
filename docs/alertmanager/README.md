# Configure alertmanager

To configure the alertmanager with alertgram you need to use the alertmanager's
webhook receiver.

For example if we have our alertgram service listening to requests in
`http://alertgram:8080/alerts` our alertmanager configuration `receiver` would be:

```yaml
receivers:
  # ...
  - name: "telegram"
    webhook_config:
      - url: "http://alertgram:8080/alerts"
```

Now we could use `telegram` in the `receiver` setting to forward the alerts to
Telegram.

A full example:

```yaml
global:
    resolve_timeout: 5m
route:
    group_wait: 30s
    group_interval: 5m
    repeat_interval: 3h
    receiver: telegram
      routes:
      # Only important alerts.
      - match_re:
          severity: ^(oncall|critical)$
        receiver: telegram-oncall

receivers:
- name: telegram
    webhook_configs:
    - url: 'http://alertgram:8080/alerts'
      send_resolved: false

- name: telegram-oncall
    webhook_configs:
    - url: 'http://alertgram:8080/alerts?chat-id=-1001111111111'
```
