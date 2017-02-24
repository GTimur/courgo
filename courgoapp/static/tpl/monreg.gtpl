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
  <label class="col-md-4 control-label" for="dept">Отдел</label>  
  <div class="col-md-5">
  <input type="file" id="flup" onchange="getfolder(event)" webkitdirectory mozdirectory msdirectory odirectory directory multiple />
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