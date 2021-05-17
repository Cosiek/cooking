{{template "base" .}}

{{define "title"}}Produkt - {{ .Ctx.product.Name }}{{end}}

{{define "content"}}
<h3>{{ .Ctx.product.Name }} <a href="/products/edit/{{ .Ctx.product.ID }}">Edytuj</a></h3>
<h5>Mamy: 1</h5>
<ul>
    <li>Kod kreskowy 1234567890 - 2 szt</li>
</ul>
{{end}}
