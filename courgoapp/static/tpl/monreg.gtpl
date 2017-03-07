{{define "head"}}
<script src="/static/duallist/jquery.bootstrap-duallistbox.js"></script>
<link rel="stylesheet" type="text/css" href="/static/duallist/bootstrap-duallistbox.css">
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
  <label class="col-md-4 control-label" for="RuleName">Название правила</label>  
  <div class="col-md-5">
  <input id="RuleName" name="RuleName" type="text" placeholder="Обработка паспортов сделок ФТС (вал. отдел)" class="form-control input-md">
  <span class="help-block">Название или краткое описание правила</span>  
  </div>
</div>

<!-- DIR/MASK -->
<div class="form-group">
  <label class="col-md-4 control-label" for="RuleDir">Директория наблюдения</label>  
  <div class="col-md-5">
  <input id="RuleDir" name="RuleDir" type="text" placeholder="C:\Папка\для\наблюдения\монитором" class="form-control input-md">
  <span class="help-block">Директория в которой монитор будет осуществлять поиск файлов</span>  
  </div>
  <label class="col-md-4 control-label" for="RuleMask">Шаблон поиска файлов</label>  
  <div class="col-md-5">
  <input id="RuleMask" name="RuleMask" type="text" placeholder="*.rar, *.txt, *.doc" class="form-control input-md">
  <span class="help-block">Маска для поиска файлов. Для нескольких значений - разделитель: запятая</span>  
  </div>
</div>

<!-- SMTP MESSAGE -->
<div class="form-group">
	<label class="col-md-4 control-label" for="MsgSubj">Тема (заголовок) извещения</label>  
	<div class="col-md-5">
		<input id="MsgSubj" name="MsgSubj" type="text" placeholder="Файлы из ФТС" class="form-control input-md">
		<span class="help-block">Тема письма при отправке уведомления</span>  
	</div>
	<label class="col-md-4 control-label" for="MsgBody">Текст извещения</label>  
	<div class="col-md-5">
			<textarea id="MsgBody" rows="5" class="form-control k-textbox" data-role="textarea" style="margin-top: 0px; margin-bottom: 0px; height: 120px;" name="MsgBody">Пересылка сообщений, полученных из ГУ Банка России по Краснодарскому краю (Краснодар) по каналам связи СВК.
Данное сообщение было сформировано автоматически.</textarea>
			<span class="help-block">Содержимое письма при отправке уведомления</span>
	</div>
</div>


<!-- RCPT SELECTION -->
<div class="form-group">
	<label class="col-md-4 control-label" for="Rcpts">Получатели извещения</label>  
	<div class="col-md-6">
		<select multiple="multiple" id="Rcpts" name="Rcpts">
			<option>Получатель 1</option>			
		</select>
		<span class="help-block">В правой стороне должны быть указаны получатели извещения о найденных файлах</span>
	</div>
</div>

<!-- ACTION SELECTION -->
<div class="form-group">
	<label class="col-md-4 control-label" for="ActionCode">Действия</label>  
	<div class="col-md-6">
		<select multiple="multiple" id="ActionCode" name="ActionCode">
			<option>Действие 1</option>			
		</select>
		<span class="help-block">В правой стороне должны быть указаны действия монитора с найденными файлами</span>
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
//Получатели - список
$("#Rcpts").bootstrapDualListbox();

$.getJSON("/gen/data/acclist", function (array) {
    $("#Rcpts").children().remove();
    $.each(array, function () {
        $('<option>').text(this).appendTo("#Rcpts");		
    })
    $("#Rcpts").bootstrapDualListbox('refresh', true);
});
//Действия - список
$("#ActionCode").bootstrapDualListbox();

$.getJSON("/gen/data/actslist", function (array) {
    $("#ActionCode").children().remove();
    $.each(array, function () {
        $('<option>').text(this).appendTo("#ActionCode");		
    })
    $("#ActionCode").bootstrapDualListbox('refresh', true);
});

$('#savebutton').click(function () {
$('#savebutton').prop('disabled', true);
var data = $("#register-subscribers").serializeObject();
// IF Rcpts or ActionCode is not array, then cast them to array
if( Object.prototype.toString.call( data["Rcpts"] ) !== '[object Array]' ) {
    data["Rcpts"]=[data["Rcpts"]];	
}
if( Object.prototype.toString.call( data["ActionCode"] ) !== '[object Array]' ) {
    data["ActionCode"]=[data["ActionCode"]];	
}
data["Post"]="SaveButton";
//alert(JSON.stringify(data));
$.ajax({                 /* start ajax function to send data */
        url: "/mon/register",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {
			//alert("REG: "+data);
			// handle AJAX redirection
			if (JSON.parse(data) == "SaveOK") {
				alert("Правило было успешно добавлено.");
				window.location = '/mon';
			}
			if (JSON.parse(data) != "SaveOK"){
				alert("Ошибка: "+data);
				$('#savebutton').prop('disabled', false);
			}	
        }
    }); 
});

$('#cancelbutton').click(function () {
$('#cancelbutton').prop('disabled', true);
$.ajax({                 /* start ajax function to send data */
        url: "/mon/register",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',        
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify({Post:"CancelButton"}),
        success: function (data) {			
			// handle AJAX redirection
			if (JSON.parse(data) == "CancelOK") {
				//alert("REG: "+data);
				//alert($('#actioncode').val());
				window.location = '/mon';
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
</script>
{{end}}