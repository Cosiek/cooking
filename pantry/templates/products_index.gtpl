{{template "base" .}}

{{define "title"}}Produkty{{end}}

{{define "content"}}
<h3>Produkty!</h3>
<a href="/products/new">Dodaj nowy</a>
<ul>
    <li><a href="/products/details/1">Ogórki na mizerię</a></li>
</ul>
{{end}}