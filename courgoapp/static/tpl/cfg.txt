<!-- Special version of Bootstrap that only affects content wrapped in .bootstrap-iso -->
<link rel="stylesheet" href="https://formden.com/static/cdn/bootstrap-iso.css" /> 

<!-- Inline CSS based on choices in "Settings" tab -->
<style>.bootstrap-iso .formden_header h2, .bootstrap-iso .formden_header p, .bootstrap-iso form{font-family: Arial, Helvetica, sans-serif; color: black}.bootstrap-iso form button, .bootstrap-iso form button:hover{color: white !important;} .asteriskField{color: red;}</style>

<!-- HTML Form (wrapped in a .bootstrap-iso div) -->
<div class="bootstrap-iso">
 <div class="container-fluid">
  <div class="row">
   <div class="col-md-6 col-sm-6 col-xs-12">
    <div class="formden_header">
     <h2>
      Настройка параметров программы
     </h2>
     <p>
      Форма настройки основных параметров COURIER GO
     </p>
    </div>
    <form method="post">
     <div class="form-group ">
      <label class="control-label " for="web-addr">
       Адрес
      </label>
      <input class="form-control" id="web-addr" name="web-addr" placeholder="localhost" type="text"/>
      <span class="help-block" id="hint_web-addr">
       IP-адрес или запись DNS
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="web-port">
       Порт
      </label>
      <input class="form-control" id="web-port" name="web-port" placeholder="8000" type="text"/>
      <span class="help-block" id="hint_web-port">
       TCP-порт
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-addr">
       Адрес
      </label>
      <input class="form-control" id="smtp-addr" name="smtp-addr" placeholder="smtp.mydomain.ru" type="text"/>
      <span class="help-block" id="hint_smtp-addr">
       IP-адрес или запись DNS
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-port">
       Порт
      </label>
      <input class="form-control" id="smtp-port" name="smtp-port" placeholder="25" type="text"/>
      <span class="help-block" id="hint_smtp-port">
       TCP-порт
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label requiredField" for="from-email">
       Email
       <span class="asteriskField">
        *
       </span>
      </label>
      <input class="form-control" id="from-email" name="from-email" placeholder="notificator@mydomain.ru" type="text"/>
      <span class="help-block" id="hint_from-email">
       Email-адрес от имени которого ведется рассылка
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-login">
       Имя пользователя
      </label>
      <input class="form-control" id="smtp-login" name="smtp-login" type="text"/>
      <span class="help-block" id="hint_smtp-login">
       login для сервера smtp
      </span>
     </div>
     <div class="form-group ">
      <label class="control-label " for="smtp-password">
       Пароль
      </label>
      <input class="form-control" id="smtp-password" name="smtp-password" type="text"/>
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
         <input name="use-tls" type="checkbox" value="Использовать TLS"/>
         Использовать TLS
        </label>
       </div>
       <span class="help-block" id="hint_use-tls">
        Если включено, то почта будет отправляться с использование протокола TLS
       </span>
      </div>
     </div>
     <div class="form-group">
      <div>
       <button class="btn btn-primary " name="submit" type="submit">
        Submit
       </button>
      </div>
     </div>
    </form>
   </div>
  </div>
 </div>
</div>
