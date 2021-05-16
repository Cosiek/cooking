{{template "base" .}}

{{define "title"}}Edytuj Produkt{{end}}

{{define "content"}}
<h3>Edytuj Produkt</h3>
<form action="" method=post>
    Nazwa: <input type="text" name="name" value="{{ .Name }}">
    Jednostka: <input type="text" name="mesure" value="{{ .Mesure }}">
    Ilość początkowa: <input type="text" name="quantity" value="2">
    <input type="submit" value="Zapisz">
</form>
{{end}}