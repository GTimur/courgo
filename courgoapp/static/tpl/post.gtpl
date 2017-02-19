{{define "head"}}{{end}}
{{define "title"}}{{.Title}} {{end}}{{define "body"}}
<div class="container" id="container">
<div class="row" style="text-align:left;">
<div class="col-md-8 col-md-offset-2">
<h1>{{.Title}}</h1>
{{.Body}}
</div>
</div>
</div>
{{end}}
{{define "scripts"}}{{end}}