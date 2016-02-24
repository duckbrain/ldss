<!DOCTYPE html>
<html>
<head>

<link rel="icon" href="/favicon.ico" type="image/x-icon">
<link rel="stylesheet" href="/css/stylesheet.css">
<title>{{ .Title }}</title>

</head>
<body>

<div class="toolbar">
	<a class="button" id="previous" {{if .Item.Previous }}href="{{.Item.Previous.Path}}"{{else}}disabled{{end}}>
		<img src="/svg/chevron-left.svg" alt="Previous">
	</a>
	<a class="button" id="up-level" {{if .Item.Parent }}href="{{.Item.Parent.Path}}"{{else}}disabled{{end}}>
		<img src="/svg/chevron-top.svg" alt="Up">
	</a>
	<form action="/lookup" method="GET">
		<input type="text" name="q">
		<button>Lookup</button>
	</form>
	<a class="button" id="next" {{if .Item.Next }}href="{{.Item.Next.Path}}"{{else}}disabled{{end}}>
		<img src="/svg/chevron-right.svg" alt="Next">
	</a>
</div>

<div class="main-content">
{{ .Content -}}
</div>

</body>
</html>
