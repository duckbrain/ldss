<!DOCTYPE html>
<html>
<head>

<link rel="icon" href="/favicon.ico" type="image/x-icon">
<link rel="stylesheet" href="/css/stylesheet.css">
<title>{{ .Title }}</title>

</head>
<body>

{{define "ItemHref"}}
	{{- if .}}href="{{.Path}}?lang={{.Language.GlCode}}"
	{{- else}}disabled{{end -}}
{{end}}

<div class="toolbar">
	<a class="button" id="previous" {{template "ItemHref" .Item.Previous}}>
		<img src="/svg/chevron-left.svg" alt="Previous">
	</a>
	<a class="button" id="up-level" {{template "ItemHref" .Item.Parent}}>
		<img src="/svg/chevron-top.svg" alt="Up">
	</a>
	<form action="/lookup" method="GET">
		<input type="text" name="q">
		<button>Lookup</button>
	</form>
	<a class="button" id="next" {{template "ItemHref" .Item.Next}}>
		<img src="/svg/chevron-right.svg" alt="Next">
	</a>
</div>

<div class="main-content">
{{ .Content -}}
</div>

</body>
</html>
