{{template "base" .}}

{{define "title"}}Spiżarnia{{end}}

{{define "content"}}
<h3>Produkty!</h3>
<a href="/products/new">Dodaj nowy</a>
<ul>
    <li><a href="/products/edit">Ogórki na mizerię</a></li>
</ul>
{{end}}