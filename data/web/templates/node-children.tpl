<h1 class="item-name">{{ .Item.Name }}</h1>
<ul class="item-children">
{{- range $key, $child := .Children }}
   <li><a href="{{ $child.Path }}?lang={{ $.LangCode }}">{{ $child.Name }}</a></li>
{{- end }}
</ul>
