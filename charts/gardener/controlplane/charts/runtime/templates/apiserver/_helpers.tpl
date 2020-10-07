{{- define "gardener-apiserver.featureGates" }}
{{- if .Values.global.apiserver.featureGates }}
- --feature-gates={{ range $feature, $enabled := .Values.global.apiserver.featureGates }}{{ $feature }}={{ $enabled }},{{ end }}
{{- end }}
{{- end -}}

{{- define "gardener-apiserver.hasAdmissionPlugins" -}}
{{- if or .Values.global.apiserver.admission.validatingWebhook.kubeconfig .Values.global.apiserver.admission.mutatingWebhook.kubeconfig .Values.global.apiserver.admission.plugins -}}
true
{{- end -}}
{{- end -}}

{{- define "gardener-apiserver.hasWebhookTokens" -}}
{{- if or .Values.global.apiserver.admission.validatingWebhook.token.enabled .Values.global.apiserver.admission.mutatingWebhook.token.enabled }}
true
{{- end -}}
{{- end -}}


{{- define "gardener-apiserver.hasAdmissionKubeconfig" -}}
{{- if or .Values.global.apiserver.admission.validatingWebhook.kubeconfig .Values.global.apiserver.admission.mutatingWebhook.kubeconfig  }}
true
{{- end -}}
{{- end -}}
