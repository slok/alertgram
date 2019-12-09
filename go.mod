module github.com/slok/alertgram

require (
	github.com/Masterminds/sprig/v3 v3.0.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/oklog/run v1.0.0
	github.com/prometheus/alertmanager v0.19.0
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/common v0.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/slok/go-http-metrics v0.5.0
	github.com/stretchr/testify v1.4.0
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

// k8s.io/client-go v1.12 subdependency is broken with Go mod.
// WTF... Why should I need to fix a subdependency? where does it come from?
// More info: https://github.com/kubernetes/client-go/issues/670
replace k8s.io/client-go v12.0.0+incompatible => k8s.io/client-go v0.0.0-20191204082520-bc9b51d240b2

go 1.13
