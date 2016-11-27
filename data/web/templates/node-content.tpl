{{ if not .HasTitle }}
	{{ if .Filtered }}
		<h1><a href="{{.Reference.URL}}">{{ .Item.Name }}</a></h1>
	{{ else }}
		<h1>{{ .Item.Name }}</h1>
	{{ end }}
{{ end }}
{{ .Content }}
