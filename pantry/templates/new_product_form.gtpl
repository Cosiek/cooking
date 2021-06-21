{{template "base" .}}

{{define "title"}}Dodaj Nowy Produkt{{end}}

{{define "content"}}
<h3>Nowy Produkt</h3>
<form action="" method=post>
    Nazwa: <input type="text" name="name" value="{{ .Ctx.form.GetName }}">
    <span>{{index .Ctx.form.Errors "name"}}{{ .Ctx.form.GetMesure }}</span>
    Jednostka: 
    <select name="mesure" value="{{ .Ctx.form.GetMesure }}">
        {{range $key, $mesure := .Mesures}}
            <option value="{{ $key }}">{{ $mesure }}</option>
        {{end}}
    </select>
    <span>{{ index .Ctx.form.Errors "mesure" }}</span>
    Ilość początkowa: <input type="text" name="quantity">
    <input type="submit" value="Zapisz">
</form>
{{end}}