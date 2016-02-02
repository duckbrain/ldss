<html>
<body>
<h1>{{ .Item.Name }}</h1>
<ul>
{{ range $key, $child := .Children }}
   <li><a href="{{ $child.Path }}?lang={{ $.LangCode }}">{{ $child.Name }}</a></li>
{{ end }}
</ul>
</body>
</html>
