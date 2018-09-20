{{define "head"}}
{{end}}
{{define "title"}}{{.Title}} {{end}}
{{define "body"}}
{{.Body}}
<form class="form-horizontal" id="monsvc-control" method="post" action="javascript:void(null);">
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
  <p>Управление обработчиком правил монитора</p> 
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="state">Состояние обработчика:</label>  
  <div class="col-md-4">
     {{.SvcState}}	
  </div>
</div>

<div class="form-group">
<br><br><br><br>
</div>

<!-- Buttons -->
<div class="form-group">
  <label class="col-md-4 control-label" for="startbutton">Управление обработчиком:</label>
  <div class="col-md-8">    
	<button id="startbutton" name="startbutton" class="btn btn-success" value="save">Запустить</button>
	<button id="stopbutton" name="stopbutton" class="btn btn-danger" value="save">Остановить</button>
    <button id="cancelbutton" name="cancelbutton" class="btn btn-default">Отмена</button>	
  </div>
</div>	
<div class="form-group">
<br><br>
</div>
<div class="form-group">
  <label class="col-md-4 control-label" for="shutdownbutton">Завершение работы приложения:</label>
  <div class="col-md-8">    	
	<button id="shutdownbutton" name="stopbutton" class="btn btn-danger" value="save">Завершить</button>
  </div>
</div>	

	<div class="col-md-2">
	</div>
</fieldset>
</form>
{{end}}
{{define "scripts"}}
<script type="text/javascript" language="javascript">
$('#stopbutton').click(function () {
$('#stopbutton').prop('disabled', true);
var data = {};//$("#monsvc-control").serializeObject();
data["Post"]="StopButton";
$.ajax({                 /* start ajax function to send data */
        url: "/mon/svc",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {			
			if (JSON.parse(data) == "StopOK") {
				//alert("Обработка правил монитора отключена.");
				//$('#startbutton').prop('disabled', false);
				window.location = '/mon/svc';
			}						
        }
    }); 
});

$('#startbutton').click(function () {
$('#startbutton').prop('disabled', true);
var data = {};
data["Post"]="StartButton";
$.ajax({                 /* start ajax function to send data */
        url: "/mon/svc",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {
			//alert("REG: "+data);
			// handle AJAX redirection
			if (JSON.parse(data) == "StartOK") {				
				//$('#stopbutton').prop('disabled', false);
				window.location = '/mon/svc';
			}						
        }
    }); 
});

$('#cancelbutton').click(function () {
   $('#cancelbutton').prop('disabled', true);
   window.location = '/';
});

$('#shutdownbutton').click(function () {
$('#shutdownbutton').prop('disabled', true);
var data = {};//$("#monsvc-control").serializeObject();
data["Post"]="ShutdownButton";
$.ajax({                 /* start ajax function to send data */
        url: "/mon/svc",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {			
			if (JSON.parse(data) == "ShutOK") {
				alert("Работа приложения была завершена.");					
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

document.addEventListener('DOMContentLoaded', function () {
var data = {};
data["Post"]="DocumentReady";
//alert(JSON.stringify(data));
$.ajax({                 /* start ajax function to send data */
        url: "/mon/svc",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {
			if (JSON.parse(data) == "StatusON") {				
				$('#startbutton').prop('disabled', true);
				$('#stopbutton').prop('disabled', false);
			}			
			if (JSON.parse(data) == "StatusOFF") {
				$('#stopbutton').prop('disabled', true);
				$('#startbutton').prop('disabled', false);
			}
        }
    });
});

</script>
{{end}}