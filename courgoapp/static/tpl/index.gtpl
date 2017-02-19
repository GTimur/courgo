{{define "head"}}
<style>
ul {
    width: 50%;
    margin: auto;
}
</style>
{{end}}
{{define "title"}}{{.Title}} {{end}}
{{define "body"}}
{{.Body}}
<form class="form-horizontal" id="register-subscribers" action="/" method="post">
<fieldset>
<!-- Form Name -->
<div class="jumbotron text-center">
  <h1>COURIER GO</h1>
  <p>Утилита уведомления о файловых событиях</p> 
</div>
<ul class="list-group">
  <li class="list-group-item"><a href="http://127.0.0.1:8000/">HOME</a></li>
  <li class="list-group-item"><a href="http://127.0.0.1:8000/">Статистика работы программы</a></li>
  <li class="list-group-item"><a href="http://127.0.0.1:8000/acc">Управление получателями рассылки</a></li>
  <li class="list-group-item"><a href="http://127.0.0.1:8000/acc/register">Регистрация получателя рассылки</a></li>
  <li class="list-group-item"><a href="http://127.0.0.1:8000/">Регистрация правила монитора</a></li>
</ul>
{{end}}
{{define "scripts"}}{{end}}