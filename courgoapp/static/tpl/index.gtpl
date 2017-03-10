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
<nav class="navbar navbar-default">
  <ul class="nav navbar-nav">
    <li><a href="{{.LnkHome}}">COURIER GO</a></li>
  </ul>
  <p class="navbar-text"></p>
</nav>
<fieldset>
<!-- Form Name -->
<div class="jumbotron text-center">
  <h1>COURIER GO</h1>
  <p>Утилита уведомления о файловых событиях</p> 
</div>
<ul class="list-group">
  <li class="list-group-item"><a href="{{.LnkHome}}/">Статистика работы программы</a></li>
  <li class="list-group-item"><a href="{{.LnkHome}}/acc">Управление получателями рассылки</a></li>
  <li class="list-group-item"><a href="{{.LnkHome}}/mon">Управление правилами монитора</a></li>
  <li class="list-group-item"><a href="{{.LnkHome}}/mon/svc">Управление обработчиком</a></li>
</ul>
{{end}}
{{define "scripts"}}{{end}}