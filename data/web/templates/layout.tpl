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
	{{ if .Item }}
	<div class="navButtons">
		<a class="button" id="previous" {{template "ItemHref" .Item.Previous}}>
			<img src="/svg/chevron-left.svg" alt="Previous">
		</a>
		<a class="button" id="next" {{template "ItemHref" .Item.Next}}>
			<img src="/svg/chevron-right.svg" alt="Next">
		</a>
	</div>
	{{ end }}
	
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
   <li><a>{{ $ref.Name }}</a> - {{$ref.LinkName}}
      <div>{{ $ref.Content }}</div>
   </li>
{{- end }}
</ul>
</div>
</aside>

<script src="/js/ldss.js"></script>
</body>
</html>
