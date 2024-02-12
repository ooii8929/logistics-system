{{/* Generate the fullname for resources */}}
{{- define "logistic-system-helm.fullname" -}}
{{- default .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Generate common labels for resources */}}
{{- define "logistic-system-helm.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/name: {{ include "logistic-system-helm.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/* Generate selector labels for resources */}}
{{- define "logistic-system-helm.selectorLabels" -}}
app.kubernetes.io/name: {{ include "logistic-system-helm.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
