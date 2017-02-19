{{define "main"}}
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{template "title" .}}</title>
    <link rel="stylesheet" type="text/css" href="/static/bootstrap/css/bootstrap.min.css">
    <!-- Optional Bootstrap theme -->
    <link rel="stylesheet" href="/static/bootstrap/css/bootstrap-theme.min.css">
    <link rel="stylesheet" href="/static/main.css">
    <script src="/static/js/jquery-3.1.1.min.js"></script>
    <script src="/static/bootstrap/js/bootstrap.min.js"></script>
  {{template "head"}}
</head>
<body>
{{template "body" .}}        
    
{{template "scripts"}}
</body>
</html>
{{end}}
