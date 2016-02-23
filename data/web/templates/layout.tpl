<!DOCTYPE html>
<html>
<head>
	<link rel="icon" href="/favicon.ico" sizes="16x16 32x32 128x128" type="image/vnd.microsoft.icon">
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
	{{ .Content }}
	</div>
</body>
</html>
