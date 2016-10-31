<!DOCTYPE html>
<html>
<head>

<link rel="icon" href="/favicon.ico" type="image/x-icon">
<link rel="stylesheet" href="/css/stylesheet.css">
<title>{{ .Title }}</title>

<meta name="viewport" content="width=device-width, initial-scale=1">

</head>
<body>

{{define "ItemHref"}}
	{{- if .}}href="{{.Path}}?lang={{.Language.GlCode}}"
	{{- else}}disabled{{end -}}
{{end}}

<div class="toolbar">
	
	<div class="breadcrumbs">
	{{ range .Breadcrumbs }}
		<a class="button" {{template "ItemHref" .}}>
		{{ if eq .Path "/" }}
			<img src="/svg/home.svg" alt="Library">
		{{ else }}
			{{ .Name }}
		{{ end }}
		</a>
	{{ end }}
	</div>
	<form class="lookup lookup-form" action="/lookup" method="GET">
		<input type="text" name="q">
	</form>
	<div class="navButtons">
		<a class="button" id="previous" {{template "ItemHref" .Item.Previous}}>
			<img src="/svg/chevron-left.svg" alt="Previous">
		</a>
		<a class="button" id="next" {{template "ItemHref" .Item.Next}}>
			<img src="/svg/chevron-right.svg" alt="Next">
		</a>
	</div>
	
</div>

<div class="main-content">
{{ .Content -}}
</div>

<script src="/js/ldss.js"></script>
</body>
</html>
