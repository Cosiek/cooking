{{template "base" .}}

{{define "title"}}Produkty{{end}}

{{define "content"}}
<h3>Produkty!</h3>
<a href="/products/new">Dodaj nowy</a>
<ul>
    {{range $product := .}}
        <li><a href="/products/details/{{ $product.ID }}">{{ $product.Name }}</a></li>
    {{end}}
</ul>
{{end}}