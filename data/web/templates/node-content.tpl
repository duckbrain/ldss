{{ if not .HasTitle }}
	{{ if .Filtered }}
		<h1><a href="{{.Reference.URL}}">{{ .Item.Name }}</a></h1>
	{{ else }}
		<h1>{{ .Item.Name }}</h1>
		<h2>{{ .Item.Subtitle }}</h2>
	{{ end }}
{{ end }}
{{ .Content }}