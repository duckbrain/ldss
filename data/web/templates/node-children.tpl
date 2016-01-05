<html>
<body>
<h1>{{ .Item.DisplayName }}</h1>
<ul>
{{ range $key, $child := .Children }}
   <li><a href="{{ $child.Path }}?lang={{ $.LangCode }}">{{ $child.DisplayName }}</a></li>
{{ end }}
</ul>
</body>
</html>