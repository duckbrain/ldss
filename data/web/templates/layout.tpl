<!DOCTYPE html>
<html>
<head>

<link rel="icon" href="/favicon.ico" type="image/x-icon">
<link rel="stylesheet" href="/css/stylesheet.css">
<title>{{ .Title }}</title>

<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="manifest" href="/manifest.webmanifest">

</head>
<body>

{{define "ItemHref"}}
	{{- if .}}href="{{.Path}}?lang={{.Language.GlCode}}"
	{{- else}}disabled{{end -}}
{{end}}

<header class="toolbar">
	
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
	<form class="lookup lookup-form" action="/search" method="GET">
		<input type="text" name="q" value="{{ .Query }}">
	</form>
	<div class="navButtons">
		<a class="button" id="previous" {{if .Item}}{{template "ItemHref" .Item.Previous}}{{else}}disabled{{end}}>
			<img src="/svg/chevron-left.svg" alt="Previous">
		</a>
		<a class="button" id="next" {{if .Item}}{{template "ItemHref" .Item.Next}}{{else}}disabled{{end}}>
			<img src="/svg/chevron-right.svg" alt="Next">
		</a>
	</div>
	
</header>

<article class="main-content">
{{ .Content -}}
</article>

<aside class="footnotes-container">
<div class="footnotes-header">
Footnotes
</div>
<div class="footnotes">
{{- range $key, $ref := .Footnotes }}
   <li id="ref-{{ $ref.Name }}" class="footnotes-footnote"><span class="footnotes-footnote-name">{{ $ref.Name }}</span> <span class="footnotes-footnote-linkName">{{$ref.LinkName}}</span> - {{ $ref.Content }}
   </li>
{{- end }}
</ul>
</div>
</aside>

<script src="/js/ldss.js"></script>
</body>
</html>
