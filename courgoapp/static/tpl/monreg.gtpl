{{define "head"}}
{{end}}
{{define "title"}}{{.Title}} {{end}}
{{define "body"}}
{{.Body}}
<form class="form-horizontal" id="register-subscribers" method="post" action="javascript:void(null);">
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
  <p>Регистрация правила монитора</p> 
</div>

<!-- Form Name -->
<!-- legend align="center">Регистрация правила монитора</legend -->
<div class="row">
		<div class="col-md-2">
		</div>
		<div class="col-md-8">
<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="fio">Правило монитора</label>  
  <div class="col-md-5">
  <input id="fio" name="rule-name" type="text" placeholder="Обработка паспортов сделок ФТС (вал. отдел)" class="form-control input-md">
  <span class="help-block">Описание правила</span>  
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="dir">Директория наблюдения</label>  
  <div class="col-md-5">
  <input id="dir" name="rule-dir" type="text" placeholder="C:\My\Files\To\Monitor" class="form-control input-md">
  <span class="help-block">Директория в которой монитор будет осуществлять поиск файлов</span>  
  </div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="mask">Шаблон поиска файлов</label>  
  <div class="col-md-5">
  <input id="mask" name="Файловая маска" type="text" placeholder="*.rar, *.txt, *.doc" class="form-control input-md">
  <span class="help-block">Маска для поиска файлов, разделитель: запятая</span>  
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="msgsubj">Заголовок извещения</label>  
  <div class="col-md-5">
  <input id="msgsubj" name="Тема письма" type="text" placeholder="Файлы из ФТС" class="form-control input-md">
  <span class="help-block">Тема письма при отправке уведомления</span>  
  </div>
</div>

<!-- Text AREA input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="msgbody">Заголовок извещения</label>  
		<div class="col-md-5">					    
        <textarea id="msgbody" rows="3" class="form-control k-textbox" data-role="textarea" style="margin-top: 0px; margin-bottom: 0px; height: 72px;" name="msgbody"></textarea>
		<span class="help-block">Содержимое письма при отправке уведомления</span>  						        
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
<script type="text/javascript" language="javascript">
$('#savebutton').click(function () {
$('#savebutton').prop('disabled', true);
var data = $("#register-subscribers").serializeObject();
data["post"]="SaveButton"
alert(JSON.stringify(data));
$.ajax({                 /* start ajax function to send data */
        url: "/acc/register",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {
			//alert("REG: "+data);
			// handle AJAX redirection
			if (JSON.parse(data) == "SaveOk") {
				alert("Получатель был успешно зарегистрирован.");
				window.location = '/acc';
			}
			if (JSON.parse(data) == "SaveNotOk"){
				alert("Данные введены с ошибкой. Получатель не был добавлен.");
				$('#savebutton').prop('disabled', false);
			}
						
        }
    }); 
});

$('#cancelbutton').click(function () {
$('#cancelbutton').prop('disabled', true);
$.ajax({                 /* start ajax function to send data */
        url: "/acc/register",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',        
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify({post:"CancelButton"}),
        success: function (data) {			
			// handle AJAX redirection
			if (JSON.parse(data) == "CancelOK") {
				//alert("REG: "+data);
				window.location = '/acc';
			}		
        }		
    }); 
});

$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        } else {
            o[this.name] = this.value || '';
        }
    });
    return o;
};
function getfolder(e) {
    var files = e.target.files;
    var path = files[0].webkitRelativePath;
    var Folder = [path]; //.split("/");
    alert(Folder[0]);
}
</script>
{{end}}