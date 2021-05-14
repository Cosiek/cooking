{{define "base"}}
<!DOCTYPE html>
<html lang="pl">
    <head>
        <meta charset="utf-8"/>
        <title>{{template "title" .}}</title>
    </head>
    <body>
        {{template "content" .}}
    </body>
</html>
{{end}}