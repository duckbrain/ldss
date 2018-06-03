<h1 class="item-name">{{ .Item.Name }}</h1>
<ul class="item-children">
{{ $sections := groupSections .Children }}
{{ if $sections }}
{{- range $section, $children := $sections }}
  <h3>{{ $section }}</h3>
  {{- range $key, $child := $children }}
  <li>
    <a href="{{ $child.Path }}?lang={{ $.LangCode }}">{{ $child.Name }}</a>
	{{ if subtitle $child }}
        {{ subtitle $child }}
    {{ end }}
  </li>
{{ end }}
{{- end }}
{{ else }}
{{- range $key, $child := .Children }}
  <li>
    <a href="{{ $child.Path }}?lang={{ $.LangCode }}">{{ $child.Name }}</a>
	{{ if subtitle $child }}
        {{ subtitle $child }}
    {{ end }}
  </li>
{{- end }}
{{ end }}
</ul>
