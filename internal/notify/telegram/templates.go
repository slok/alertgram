package telegram

var defTemplate = `
🚨🚨 FIRING {{ .Alerts | len }} 🚨🚨
{{- range .Alerts }}
💥💥💥 <b>{{ .Labels.alertname }}</b> 💥💥💥
  {{ .Annotations.message }}
  {{- range $key, $value := .Labels }}
	{{- if ne $key "alertname" }}  
	{{- if hasPrefix "http" $value }}
	🔹 <a href="{{ $value }}">{{ $key }}</a>
	{{- else }}
	🔹 {{ $key }}: {{ $value }}
	{{- end}}
	{{-  end }}
  {{- end}}
  {{- range $key, $value := .Annotations }}
  	{{- if ne $key "message" }}  
	{{- if hasPrefix "http" $value }}
	🔸 <a href="{{ $value }}">{{ $key }}</a>
	{{- else }}
	🔸 {{ $key }}: {{ $value }}
	{{- end}}
	{{- end}}
  {{- end}}
{{- end }}
`
