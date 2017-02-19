{{define "head"}}
{{end}}
{{define "title"}}{{.Title}} {{end}}
{{define "body"}}
{{.Body}}
<form class="form-horizontal" id="register-subscribers" method="post">
<fieldset>
<fieldset>
<!-- Form Name -->
<div class="jumbotron text-center">
  <h1>COURIER GO</h1>
  <p>Регистрация подписчика</p> 
</div>

<!-- Form Name -->
<!-- legend align="center">Регистрация подписчика</legend -->
<div class="row">
		<div class="col-md-2">
		</div>
		<div class="col-md-8">
<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="fio">Ф.И.О.</label>  
  <div class="col-md-5">
  <input id="fio" name="fio" type="text" placeholder="Александр Сергеевич Пушкин" class="form-control input-md">
  <span class="help-block">Ф.И.О сотрудника</span>  
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="dept">Отдел</label>  
  <div class="col-md-5">
  <input id="dept" name="dept" type="text" placeholder="Юридический отдел" class="form-control input-md">
  <span class="help-block">Название отдела, включая слово "отдел"</span>  
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="email">E-mail:</label>  
  <div class="col-md-4">
  <input id="email" name="email" type="text" placeholder="aspushkin@ymkbank.ru, pushkinas@ymkbank.ru" class="form-control input-md">
  <span class="help-block">Адреса электронной почты, разделитель: запятая</span>  
  </div>
</div>

<!-- Button (Double) -->
<div class="form-group">
  <label class="col-md-4 control-label" for="savebutton"></label>
  <div class="col-md-8">
    <button id="savebutton" name="savebutton" class="btn btn-success" value="save">Добавить</button>
    <button id="cancelbutton" name="cancelbutton" class="btn btn-danger">Отмена</button>
  </div>
</div>	
		</div>
		<div class="col-md-2">
		</div>
</div>	
</fieldset>
</form>
{{end}}
{{define "scripts"}}
{{end}}