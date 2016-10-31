{{ if not .HasTitle }}
<h1>{{ .Item.Name }}</h1>
{{ end }}
{{ .Content }}
