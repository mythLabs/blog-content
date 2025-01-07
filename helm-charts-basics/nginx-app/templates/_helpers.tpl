{{- define "app.fullname" -}}
{{- .Release.Name }}-{{ .Chart.Name }}-nginx
{{- end -}}

{{- define "app.labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app: nginx
{{- end -}}