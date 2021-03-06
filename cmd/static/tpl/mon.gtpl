{{define "head"}}	
	<link rel="stylesheet" href="/static/bootstrap-table/bootstrap-table.css">	
	<script src="/static/bootstrap-table/bootstrap-table.js"></script>
	<script src="/static/bootstrap-table/locale/bootstrap-table-ru-RU.js"></script>
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
  <p>Управление правилами монитора</p> 
</div>
<!-- Form Name -->
<!-- legend align="center">Управление правилами монитора</legend -->
	<div class="row">
		<div class="col-md-2">
		</div>
		<div class="col-md-8">		 
		<div id="toolbar">
         <button id="btnremove" class="btn btn-danger" data-toggle="modal" data-target="#confirm-delete">
		 Удалить</button>
		 <button id="btnadd" class="btn btn-default" >
         Добавить</button>
		 </div>
		</div>
		<div class="col-md-2">
		</div>
	</div>
<div class="row">
		<div class="col-md-2">
		</div>
		<div class="col-md-8">		
		<table id="table"
		class="table table-striped"
		data-toggle="table"		
       data-height="460"
	   data-toolbar="#toolbar"
       data-click-to-select="true"
	   data-select-item-name="selectItemName"
	   data-id-field="Id"
	   data-pagination="true"   
	   data-page-list="[10, 25, 50, 100, ALL]"
       data-url="/gen/data/moncol.json"	   
	   data-method="post"   
       data-single-select="true"
       data-content-type="application/x-www-form-urlencoded">
    <thead>
        	<tr>
                <th data-field="state" data-checkbox="true"></th>
                <th data-field="Id">ID</th>
                <th data-field="Desc">Описание</th>
                <th data-field="Folder">Папка</th>
				<th data-field="Mask" data-formatter="arrayStringFormatter">Маска</th>
				<th data-field="Action" data-formatter="actionStringFormatter">Действия</th>
				<th data-field="Sid" data-formatter="arrayStringFormatter">Получатели</th>
            </tr>
    </thead>   
</table>		</div>
		<div class="col-md-2">
		</div>
	</div>
</fieldset>
<body>
    <div class="modal fade" id="confirm-delete" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
            
                <div class="modal-header">
                    <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
                    <h4 class="modal-title" id="myModalLabel">Удаление подписчика</h4>
                </div>
            
                <div class="modal-body">
                    <p>Вы собираетесь удалить правило из базы данных.</p>					
					<p>К удалению:</p>
					<div id="recid"></div>
					<p><br>Вы уверены что это необходимо сделать?</p>
                    <p class="debug-url"></p>
                </div>
                
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">Отменить</button>
                    <a class="btn btn-danger btn-ok">Удалить</a>
                </div>
            </div>
        </div>
    </div>
{{end}}
{{define "scripts"}}
<script>
$('#btnremove').prop('disabled', true);
/*Register new account*/
 $('#btnadd').click(function () {
	//Handshake function for JSON request!    
    $.ajax({                 /* start ajax function to send data */
        url: "/mon",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',        
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify({Post:"Register"}),
        success: function (data) {
			// alert("SUCCESS: "+JSON.parse(data));
			// handle AJAX redirection
			if (JSON.parse(data) == "RegisterOK") {
				//alert("REG: "+data);
				window.location = '/mon/register';
			}		
        }		
    });	
 });

var $table = $('#table');

// Handle string array json fields [{},{}]
function arrayStringFormatter(value) {
    var Field = '';

    // Loop through the authors object
    for (var index = 0; index < value.length; index++) {
        Field += value[index];

        // Only show a comma if it's not the last one in the loop
        if (index < (value.length - 1)) {
            Field += ',<br>';
        }
    }
    return Field;
}   

// Handle string array json fields Action [{A,B},{A,B}]
function actionStringFormatter(value) {
    var Field = '';

    // Loop through the authors object
    for (var index = 0; index < value.length; index++) {
        Field += value[index].Desc;

        // Only show a comma if it's not the last one in the loop
        if (index < (value.length - 1)) {
            Field += ',<br>';
        }
    }
    return Field;
}  

/*Delete confirmation*/
 $('#confirm-delete').on('click', '.btn-ok', function(e) {
  var $modalDiv = $(e.delegateTarget);
  $modalDiv.addClass('loading');

	$modalDiv.modal('hide').removeClass('loading');    
	var $table = $('#table');	
	var vals = [];
            $('input[name="selectItemName"]:checked').each(function () {
                //index.push($(this).data('index'));			
				vals.push($(this).val());			
            });
		var jsonstring = JSON.stringify($table.bootstrapTable('getSelections'));
		var json = JSON.parse(jsonstring);	
		var data = {};
		data["Post"]="RemoveButton";
		data["Id"]=String(json[0]['Id']);
		alert("REG: "+JSON.stringify(data));
				
		$.ajax({                 /* start ajax function to send data */
        url: "/mon",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',        
        error: function () { alert("handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {			
			// handle AJAX redirection
			if (JSON.parse(data) == "RemoveOK") {
				//alert("REG: "+data);
				window.location = '/mon';				
			}
			if (JSON.parse(data) != "RemoveOK") {
				alert(data);
				//window.location = '/acc';
			}				
        }		
    }); 
	alert("Вы только что удалили запись! Id="+vals.join(', '));
});

// Fires when user check a row
$('#table').bootstrapTable({
	//onCheck - при отметке чекбокса
    onCheck: function (row, $element) {
		$("#btnremove").prop('disabled', false);
        var jsonstring = JSON.stringify($table.bootstrapTable('getSelections'));
		var json = JSON.parse(jsonstring);
		///$(".modal-body #userid").val( json[0]['Name'] );
		$(".modal-body #recid").html("ID: "+json[0]['Id']+"<br>Описание: "+json[0]['Desc']+"<br>Папка: "+json[0]['Folder']+"<br>Маска:"+json[0]['Mask']);
		//console.log(row, $element);
    },
	
	onUncheck: function (row, $element) {
		$(".modal-body #userid").text("");
		$('#btnremove').prop('disabled', true);
		
		
		console.log("uncheck");
    }
});

</script>
{{end}}