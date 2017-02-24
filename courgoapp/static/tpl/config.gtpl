{{define "head"}}
{{end}}
{{define "title"}}{{.Title}} {{end}}
{{define "body"}}
{{.Body}}
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
  <p>Установка основных параметров программы</p> 
</div>
  <div class="row">
   <div class="col-md-2">
   </div>
   <div class="col-md-8">
    <div class="formden_header">
     <h2>
      Настройка параметров программы
     </h2>
    </div>
    <form method="post" id="EditConfig" action="javascript:void(null);">
     <div class="form-group ">
      <label class="control-label " for="web-addr">
       Адрес сервера web-интерфейса управления
      </label>
      <input class="form-control" id="web-addr" name="web-addr" placeholder="localhost" type="text"/>
      <span class="help-block" id="hint_web-addr">
       IP-адрес или запись DNS
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="web-port">
       Порт HTTP-сервера управления
      </label>
      <input class="form-control" id="web-port" name="web-port" placeholder="8000" type="text"/>
      <span class="help-block" id="hint_web-port">
       TCP-порт
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-addr">
       Адрес SMTP сервера через который будет выполняться рассылка
      </label>
      <input class="form-control" id="smtp-addr" name="smtp-addr" placeholder="smtp.mydomain.ru" type="text"/>
      <span class="help-block" id="hint_smtp-addr">
       IP-адрес или запись DNS
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-port">
       Порт SMTP сервера
      </label>
      <input class="form-control" id="smtp-port" name="smtp-port" placeholder="25" type="text"/>
      <span class="help-block" id="hint_smtp-port">
       TCP-порт
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label requiredField" for="from-email">
       Email
      </label>
      <input class="form-control" id="from-email" name="from-email" placeholder="notificator@mydomain.ru" type="text"/>
      <span class="help-block" id="hint_from-email">
       Email-адрес от имени которого ведется рассылка
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-login">
       Имя пользователя SMTP
      </label>
      <input class="form-control" id="smtp-login" name="smtp-login" type="text"/>
      <span class="help-block" id="hint_smtp-login">
       login для сервера smtp
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-password">
       Пароль SMTP
      </label>
      <input class="form-control" id="smtp-password" name="smtp-password" type="password"/>
      <span class="help-block" id="hint_smtp-password">
       пароль для учетной записи smtp
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label ">
       Настройка протокола TLS
      </label>
      <div class=" ">
       <div class="checkbox">
        <label class="checkbox">
         <input name="use-tls" type="checkbox" value="useTLS"/>
         Использовать TLS
        </label>
       </div>
       <span class="help-block" id="hint_use-tls">
        Если включено, то почта будет отправляться с использованием протокола TLS
       </span>
      </div>
     </div>
     <div class="form-group">
      <div>
       <button class="btn btn-primary " name="save" data-toggle="modal" data-target="#confirm-save">
        Сохранить
       </button>
      </div>
     </div>
    </form>
   </div>
   <div class="col-md-2">
   </div>
  </div>
 <body>
    <div class="modal fade" id="confirm-save" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
            
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                    <h4 class="modal-title" id="myModalLabel">Сохранение конфигурации</h4>
                </div>
            
                <div class="modal-body">
                    <p>Конфигурация будет перезаписана новыми данными.</p>					
					<p>Путь к файлу конфигурации:</p>
					<div id="userid"></div>
					<p><br>Вы уверены что это необходимо сделать?</p>
                    <p class="debug-url"></p>
                </div>
                
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Отменить</button>
                    <a class="btn btn-danger btn-ok">Сохранить</a>
                </div>
            </div>
        </div>
    </div>
{{end}}
{{define "scripts"}}
<script>
/*Save confirmation*/
 $('#confirm-save').on('click', '.btn-ok', function(e) {
  var $modalDiv = $(e.delegateTarget);
  $modalDiv.addClass('loading');
	$modalDiv.modal('hide').removeClass('loading');
	var data = $("#EditConfig").serializeObject();
	data["post"]="SaveConfig"
	//alert(JSON.stringify(data));
	console.log(JSON.stringify(data));
	//Handshake function for JSON request!    
    $.ajax({                 /* start ajax function to send data */
        url: "/config",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',        
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {
			//alert("SUCCESS: "+JSON.parse(data));
			// handle AJAX redirection
			if (JSON.parse(data) == "SaveOK") {
				alert("Конфигурация сохранена!"); 
				//alert("REG: "+data);
				window.location = '/';
			}
			if (JSON.parse(data) != "SaveOK") {
				alert("Конфигурация не была сохранена. Причина:"+data);
				//window.location = '/config';
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