{{template "base" .}}

{{define "title"}}Dodaj Nowy Produkt{{end}}

{{define "content"}}
<h3>Nowy Produkt</h3>
<form action="" method=post>
    Nazwa: <input type="text" name="name">
    <span>{{index .Ctx.form.Errors "name"}}</span>
    Jednostka: 
    <select name="mesure" value="">
        {{range $key, $mesure := .Mesures}}
            <option value="{{ $key }}">{{ $mesure }}</option>
        {{end}}
    </select>
    <span>{{ index .Ctx.form.Errors "mesure" }}</span>
    Ilość początkowa: <input type="text" name="quantity">
    <input type="submit" value="Zapisz">
</form>
{{end}}