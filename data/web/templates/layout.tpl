<!DOCTYPE html>
<html>
<head>

<link rel="icon" href="/favicon.ico" type="image/x-icon">
<link rel="stylesheet" href="/css/stylesheet.css">
<title>{{ .Title }}</title>

</head>
<body>

<div class="toolbar">
	<form action="/lookup" method="GET">
		<input type="text" name="q">
		<button>Lookup</button>
	</form>
</div>

<div class="main-content">
{{ .Content -}}
</div>

</body>
</html>
