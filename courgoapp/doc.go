package main

/*
	Courier Go
*/

/*
	TO DO

	//- Для работы правил монитора организовать очередь с интервалами ожидания до повторного запуска правила
	- Правила в которых есть пересечения по маскам с другими правилами должны передавать таковым вслед за собой
	  управление в очередь (ближайшему по номеру правилу) без присвоения статуса файлу "обработан"
	- Оперативная история обработанных сообщений хранится в памяти приложения и очищается в начале дня. (Окно хранения - сутки.)
	  История:
	      		Дата,
	  		Время,
	  		Номер правила,
	  		Маска,
	  		Статус (обработан, ожидается дальнейшая обработка),
	  		Полное имя файла включая путь
	- История обработанных сообщений сохраняется в процессе работы программы в файл для просмотра статистики работы программы.



	- добавить генерацию настроек defaults флагом -init
	- добавить редактирование правил через webctl
	- добавить шифрование/дешифрование конфигурационных файлов
	- добавить авторизацию webctl

*/


/* Folders monitor

  MONITOR -> Monitor search rules (Monitor Search Rule ID = MSRID)
          -> Check ALL FOLDERS from MSRID's
          -> Send notification to subscribers from address book (Subscriber ID = SID)
          -> Do possible Actions (Action ID = AID)          
 
  MONITOR RULE	-> Folder
              	-> Mask
                -> Subscribers (SID) []
                -> Notification e-mail message text header and body
                -> Reader (RID)
                -> Actions (AID) []

*/

/* Address book

   ID
   Name
   Department
   []Mail

*/

/* Actions

   Send e-mail message
   Copy file from SRC to DST folder
*/

/* Reader

   Unpack
   Unpack and get text
*/

/* Courier (e-mail, may be other way) 
*/

/* Archiver "Archi"
*/

/* Web-control

   Statistics
   Subscribers book manager
   Monitor rules
*/

/* JSON config */

//    Path PTK_IN
//    Path PTK_OUT
//    WebInt Address
//    WebInt Port

//    Если файла конфигурации нет, то он будет создан со значениями по умолчанию.



/*

*/






