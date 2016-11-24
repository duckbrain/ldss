<h1 class="item-name">
	Search Results for "{{.SearchString}}"
	in {{ .Item.Name }}
</h1>
<ul class="item-children">
{{- range $key, $child := .SearchResults }}
	<li>
		{{ $child.Weight }}
		<a href="{{ $child.URL }}">{{ $child.URL }}</a>
	</li>
{{- end }}
</ul>
