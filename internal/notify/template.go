package notify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/slok/alertgram/internal/model"
)

// ErrRenderTemplate will be used when there is an error rendering the alerts
// to a template.
var ErrRenderTemplate = errors.New("error rendering template")

// TemplateRenderer knows how to render an alertgroup to get the
// final notification message.
type TemplateRenderer interface {
	Render(ctx context.Context, ag *model.AlertGroup) (string, error)
}

// TemplateRendererFunc is a helper function to use funcs as TemplateRenderer types.
type TemplateRendererFunc func(ctx context.Context, ag *model.AlertGroup) (string, error)

// Render satisfies TemplateRenderer interface.
func (t TemplateRendererFunc) Render(ctx context.Context, ag *model.AlertGroup) (string, error) {
	return t(ctx, ag)
}

// NewHTMLTemplateRenderer returns a new template renderer using the go HTML
// template renderer.
// The templates use https://github.com/Masterminds/sprig to render.
func NewHTMLTemplateRenderer(tpl string) (TemplateRenderer, error) {
	t, err := template.New("tpl").Funcs(sprig.FuncMap()).Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("error rendering template: %w", err)
	}

	return TemplateRendererFunc(func(_ context.Context, ag *model.AlertGroup) (string, error) {
		return renderAlertGroup(ag, t)
	}), nil
}

// renderAlertGroup takes an alertGroup and renders on the given template.
func renderAlertGroup(ag *model.AlertGroup, t *template.Template) (string, error) {
	var b bytes.Buffer
	err := t.Execute(&b, ag)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrRenderTemplate, err)
	}

	return b.String(), nil
}

type defRenderer int

// DefaultTemplateRenderer is the default renderer that will render the
// alerts using a premade HTML template.
const DefaultTemplateRenderer = defRenderer(0)

func (defRenderer) Render(_ context.Context, ag *model.AlertGroup) (string, error) {
	return renderAlertGroup(ag, defTemplate)
}

var defTemplate = template.Must(template.New("def").Funcs(sprig.FuncMap()).Parse(`
🚨🚨 FIRING ALERTS 🚨🚨
{{- range .Alerts }}
{{- if .IsFiring }}

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
{{- end }}
`))
