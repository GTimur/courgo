/*
// Web-сервер реализует интерфейс управления приложением

   Statistics
   Subscribers book manager
   Monitor rules
*/
package courgo

import (
	"net"
	"net/http"
	"log"
	"fmt"
	"html/template"
	"github.com/braintree/manners"
	"github.com/GeertJohan/go.rice"
	"strings"
	"encoding/json"
	"strconv"
)

type WebCtl struct {
	host     net.IP
	port     uint16
	islisten bool
}

type Page struct {
	Title    string
	Body     template.HTML
	LnkHome  string
	SvcState template.HTML
}

var (
	// компилируем шаблоны

	/*conf_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "config.gtpl")))
	home_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "index.gtpl")))
	acc_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "acc.gtpl")))
	register_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "register.gtpl")))
	mon_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "mon.gtpl")))
	monreg_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "monreg.gtpl")))
	monsvc_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "monsvc.gtpl")))*/


	// Переменная сообщающая программе о необходимости при первой же возможности завершить работу.
	WaitExit bool
	// Переменная для отображения инфорамации о том сколько времени осталось до очередного запуска мониторов.
	TimeRemain int
	// Интервал задающий частоту запуска правил монитора в секундах
	Interval = 60
)

//Функции установки значений
func (w *WebCtl) SetHost(host net.IP) {
	w.host = host
}

func (w *WebCtl) SetPort(port uint16) {
	w.port = port
}
/**/
func (w WebCtl) connString() string {
	return fmt.Sprintf("%s:%d", w.host.String(), w.port)
}
func (w *WebCtl) IsListen() bool {
	return w.islisten
}

//возвращает true если операция закрытия сервера выполнялась
func (w *WebCtl) Close() bool {
	if !w.islisten {
		return w.islisten
	}
	manners.Close()
	w.islisten = false
	return !w.islisten
}


/*Сервер*/
//Запускает goroutine ListenAndServe
//Может изменять accbook - справочник подписантов
func (w *WebCtl) StartServe() (err error) {
	// для отдачи сервером статичных файлов из папки public/static

	//fs := http.FileServer(http.Dir("./static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	/* Добавил go.rice для объединения статических файлов в exe-шник */
	box, err := rice.FindBox("static")
	if err != nil {
		log.Fatalln("go.rice embedding error: ", err)
	}
	cssFileServer := http.StripPrefix("/static/", http.FileServer(box.HTTPBox()))
	http.Handle("/static/", cssFileServer)
	http.HandleFunc("/gen/data/test.book", accbook) //Генерация JSON с данными подписчиков для таблицы
	http.HandleFunc("/gen/data/acclist", accduallist) //Генерация JSON с данными подписчиков для duallist
	http.HandleFunc("/gen/data/actslist", actsduallist) //Генерация JSON с данными действий для duallist
	http.HandleFunc("/gen/data/moncol.json", moncol) //Генерация JSON с данными правил монитора
	http.HandleFunc("/", urlhome) //Каждый запрос вызывает обработчик
	http.HandleFunc("/acc", urlacc) //Страница с таблицей подписчиков
	http.HandleFunc("/acc/register", urlregister) //Страница регистрации подписчика
	http.HandleFunc("/mon", urlmon) //Страница с таблицей правил монитора
	http.HandleFunc("/mon/register", urlmonreg) //Страница создания правила монитора
	http.HandleFunc("/mon/svc", urlmonsvc) //Страница управления обработчиком (вкл/выкл)
	http.HandleFunc("/config", urlconfig) //Страница настройки приложения
	go func() {
		log.Fatalln("WebCtl:", manners.ListenAndServe(w.connString(), http.DefaultServeMux))
	}()
	w.islisten = true
	return err
}

//Обработчик запросов для home
func urlhome(w http.ResponseWriter, r *http.Request) {
	title := "COURIER GO"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	page := Page{title, template.HTML(body), lnkhome, "" }

	/*GO.RICE*/
	tmplMessage := goriceTmpl("static", "tpl/main.gtpl", "tpl/index.gtpl")
	if err := tmplMessage.Execute(w, page); err != nil {
		log.Println("Template error:", err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
	/*GO.RICE*/

	//if err := home_template.ExecuteTemplate(w, "main", page); err != nil {
	//	log.Println(err.Error())
	//	http.Error(w, http.StatusText(500), 500)
	//}
}

//Обработчик запросов для /gen/data/actslist
//Генерирует список доступных действий (actions) формате JSON для передачи в duallist
//моделирует работу с файлом
func actsduallist(w http.ResponseWriter, r *http.Request) {
	data := "[" + "\"10:Отправить как вложение\",\"11:Отправить содержимое архива\",\"20:Только уведомить(без вложения)\",\"1000:Не обрабатывать\"" + "]"
	fmt.Fprint(w, data)
}


//Обработчик запросов для /gen/data/acclist
//Генерирует содержимое адресной книги в формате JSON для передачи в duallist
//моделирует работу с файлом
func accduallist(w http.ResponseWriter, r *http.Request) {
	data := "["
	for i, acc := range GlobalBook.account {
		if i > 0 {
			data += ","
		}
		//data += "{\"id\":" + strconv.Itoa(int(acc.Id())) + ",\"name\":" + "\""+acc.Name() + "\"}"
		data += "\"" + strconv.Itoa(int(acc.Id())) + " : " + acc.Name() + "\""
	}
	data += "]"
	fmt.Fprint(w, data)
}

//Обработчик запросов для /gen/data/test.book
//Генерирует содержимое адресной книги в формате JSON для передачи в таблицу
//моделирует работу с файлом для таблицы
func accbook(w http.ResponseWriter, r *http.Request) {
	data := ""
	var err error
	if data, err = GlobalBook.StringJSON(0); err != nil {
		data = ""
	}
	fmt.Fprint(w, data)
}

// Обработчик запросов для /gen/data/moncol.json
// Генерирует в формате JSON имеющиеся правила монитора
// для передачи в таблицу
func moncol(w http.ResponseWriter, r *http.Request) {
	data := ""
	var err error
	if data, err = GlobalMonCol.StringJSON(); err != nil {
		data = ""
	}
	fmt.Fprint(w, data)
}

//обработчик для /acc/register
func urlregister(w http.ResponseWriter, r *http.Request) {
	//var accnt Acc
	//fmt.Println("method:", r.Method) //get request method
	title := "Регистрация подписчика"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	page := Page{title, template.HTML(body), lnkhome, ""}

	/*GO.RICE*/
	tmplMessage := goriceTmpl("static", "tpl/main.gtpl", "tpl/register.gtpl")
	/*GO.RICE*/

	if r.Method == "GET" {
		if err := tmplMessage.Execute(w, page); err != nil {
			log.Println("Template error:", err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		/*if err := register_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}*/
	} else {

		err := r.ParseForm()
		if err != nil {
			log.Println("Form parse error:", err)
		}

		if strings.Contains(r.Form.Encode(), "cancelbutton") {
			fmt.Println(r.Form.Encode())
		} else {
			dec := json.NewDecoder(r.Body)
			defer r.Body.Close()

			// Массив данных JSON для получения данных из формы (ajax)
			// Первый элемент должен содержать код действия:
			// "RemoveAcc" - удалить аккаунт из книги
			// "NewAcc" - добавить аккаунт в книгу
			type jsonPOSTData map[string]string

			var jh jsonPOSTData

			err := dec.Decode(&jh)
			if err != nil {
				log.Println("Handshake error: ", err)
			}
			//log.Println(jh["fio"],jh["dept"],strings.Split(jh["email"], ","))


			//JSONheader
			switch jh["post"] {
			case "SaveButton" :

				if err := RegisterAccount(jh["fio"], jh["dept"], strings.Split(jh["email"], ",")); err != nil {
					log.Println("Ошибка регистрации подписчика:", err)
					enc := json.NewEncoder(w)
					enc.Encode("SaveNotOk")
				} else {
					if err := WriteJSONFile(&GlobalBook); err != nil {
						log.Println("Не удалось сохранить в файл нового подписчика:", err)
						enc := json.NewEncoder(w)
						enc.Encode("SaveNotOk")
					} else {
						enc := json.NewEncoder(w)
						enc.Encode("SaveOk")
					}
				}

			case "CancelButton" :
				enc := json.NewEncoder(w)
				enc.Encode("CancelOK")
			default:
				//Отправляем ответ на POST-запрос
				//для предотвращения ошибки JSON parse error в ajax методе
				enc := json.NewEncoder(w)
				enc.Encode("No actions required")
			}
		}
	}
}

//Обработчик для /acc
//Если в POST передан нужный код, то выполняется действие
func urlacc(w http.ResponseWriter, r *http.Request) {
	//var accnt Acc
	//fmt.Println("method:", r.Method) //get request method
	title := "Подписчики рассылки (адресная книга)"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	page := Page{title, template.HTML(body), lnkhome, ""}

	/*GO.RICE*/
	tmplMessage := goriceTmpl("static", "tpl/main.gtpl", "tpl/acc.gtpl")
	/*GO.RICE*/

	if r.Method == "GET" {
		if err := tmplMessage.Execute(w, page); err != nil {
			log.Println("Template error:", err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		/*if err := acc_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}*/
	} else {
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		// Массив данных JSON для получения данных из формы (ajax)
		// Первый элемент должен содержать код действия:
		// "RemoveAcc" - удалить аккаунт из книги
		// "NewAcc" - добавить аккаунт в книгу
		type jsonPOSTData []string

		var jh jsonPOSTData

		err := dec.Decode(&jh)
		if err != nil {
			log.Println("Handshake error: ", err)
		}
		//log.Println(jh[0])

		switch jh[0] {
		case "RemoveAcc" :
			enc := json.NewEncoder(w)
			enc.Encode("remove")
		case "NewAcc" :
			enc := json.NewEncoder(w)
			enc.Encode("register")
		//http.Redirect(w, r, "/", http.StatusFound)
		default:
			//Отправляем ответ на POST-запрос
			//для предотвращения ошибки JSON parse error в ajax методе
			enc := json.NewEncoder(w)
			enc.Encode("No action requested.")

		}
		//if strings.Compare(jh[1], "RemoveAcc") * strings.Compare(jh[1], "NewAcc") > 0 {

		//}
	}
}

//Обработчик для /mon
//Если в POST передан нужный код, то выполняется действие
func urlmon(w http.ResponseWriter, r *http.Request) {
	title := "Управление правилами монитора"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	page := Page{title, template.HTML(body), lnkhome, ""}
	/*GO.RICE*/
	tmplMessage := goriceTmpl("static", "tpl/main.gtpl", "tpl/mon.gtpl")
	/*GO.RICE*/

	if r.Method == "GET" {
		if err := tmplMessage.Execute(w, page); err != nil {
			log.Println("Template error:", err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		/*if err := mon_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}*/
	} else {
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		// Массив данных JSON для получения данных из формы (ajax)
		// Первый элемент должен содержать код действия
		type jsonPOSTData map[string]string

		var jh jsonPOSTData

		err := dec.Decode(&jh)
		if err != nil {
			log.Println("Handshake error: ", err)
		}

		switch jh["Post"] {
		case "RemoveButton" :
			id, cerr := strconv.Atoi(jh["Id"])
			if cerr != nil {
				enc := json.NewEncoder(w)
				enc.Encode("Не удалось удалить указанное правило.")
				break
			}
			if err := GlobalMonCol.RemoveColElm(uint64(id)); err != nil {
				enc := json.NewEncoder(w)
				enc.Encode("Не удалось удалить указанное правило.")
				break
			}
			// Запишем изменения в файл
			if err := WriteJSONFile(&GlobalMonCol); err != nil {
				// Если записать не удалось, то сообщаем об ошибке
				enc := json.NewEncoder(w)
				enc.Encode("Не удалось удалить указанное правило.")
				break
			}
			enc := json.NewEncoder(w)
			enc.Encode("RemoveOK")
		case "Register" :
			enc := json.NewEncoder(w)
			enc.Encode("RegisterOK")

		default:
			//Отправляем ответ на POST-запрос
			//для предотвращения ошибки JSON parse error в ajax методе
			enc := json.NewEncoder(w)
			enc.Encode("No action requested.")
		}
	}
}

//Обработчик для /mon/register
//Если в POST передан нужный код, то выполняется действие
func urlmonreg(w http.ResponseWriter, r *http.Request) {
	title := "Управление правилами монитора"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	page := Page{title, template.HTML(body), lnkhome, ""}
	/*GO.RICE*/
	tmplMessage := goriceTmpl("static", "tpl/main.gtpl", "tpl/monreg.gtpl")
	/*GO.RICE*/

	if r.Method == "GET" {
		if err := tmplMessage.Execute(w, page); err != nil {
			log.Println("Template error:", err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		/*if r.Method == "GET" {
		if err := monreg_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}*/
	} else {
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		// Массив данных JSON для получения данных из формы (ajax)
		// Первый элемент должен содержать код действия
		type jsonPOSTData struct {
			RuleName, RuleDir, RuleMask, MsgSubj, MsgBody string
			Rcpts, ActionCode                             []string
			Post                                          string
		}

		var jh jsonPOSTData

		err := dec.Decode(&jh)
		if err != nil {
			log.Println("Handshake error: ", err)
		}
		//fmt.Println(jh)
		switch jh.Post {
		case "SaveButton" :
			// Добавим все данные считанные из формы в переменную монитора.
			var mon Monitor = Monitor{
				id:GlobalMonCol.MaxID() + 1,
				desc:jh.RuleName,
				folder:jh.RuleDir,
				mask:strings.Split(strings.TrimRight(strings.Replace(jh.RuleMask, ", ", ",", -1)," "), ","),
				msgSubject:jh.MsgSubj,
				msgBody:strings.Replace(jh.MsgBody, "\r\n", "<br>", -1),
			}
			var id int
			var err error
			for _, str := range jh.Rcpts {
				if len(str) == 0 {
					continue
				}
				if id, err = strconv.Atoi(strings.Trim(str[0:strings.Index(str, ":")], " ")); err != nil {
					enc := json.NewEncoder(w)
					enc.Encode(err.Error())
					break
				}
				mon.sid = append(mon.sid, uint64(id))
			}
			for _, str := range jh.ActionCode {
				if len(str) == 0 {
					continue
				}
				if id, err = strconv.Atoi(strings.Trim(str[0:strings.Index(str, ":")], " ")); err != nil {
					enc := json.NewEncoder(w)
					enc.Encode(err.Error())
					break
				}
				mon.action = append(mon.action, Action{id:uint64(id), desc:str[strings.Index(str, ":") + 1:]})
			}
			// Добавим новый монитор в коллекцию, если возникнут ошибки - пишем в форму
			if err := GlobalMonCol.AddMonitor(mon); err != nil {
				enc := json.NewEncoder(w)
				enc.Encode(err.Error())
				break
			}
			// Запишем изменения в файл
			if err := WriteJSONFile(&GlobalMonCol); err != nil {
				// Если записать не удалось, то удаляем правило из коллекции
				GlobalMonCol.RemoveColElm(mon.id)
				enc := json.NewEncoder(w)
				enc.Encode(err.Error())
				break
			}
			enc := json.NewEncoder(w)
			enc.Encode("SaveOK")
		case "CancelButton" :
			enc := json.NewEncoder(w)
			enc.Encode("CancelOK")
		default:
			//Отправляем ответ на POST-запрос
			//для предотвращения ошибки JSON parse error в ajax методе
			enc := json.NewEncoder(w)
			enc.Encode("No action requested.")
		}
	}
}


//Обработчик для /mon
//Если в POST передан нужный код, то выполняется действие
func urlmonsvc(w http.ResponseWriter, r *http.Request) {
	title := "Управление правилами монитора"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	svcstate := `<h4><span class="label label-success">Выполняется</span></h4>В данный момент обработчик правил монитора запущен и будет выполнять правила каждые ` + strconv.Itoa(Interval) + ` секунд.<Br>Обработчик выполнит запуск правил через ` + strconv.Itoa(TimeRemain) + ` секунд.`
	if !MonSvcState {
		svcstate = `<h4><span class="label label-danger">Остановлен</span></h4>В данный момент монитор остановлен.`
	}
	page := Page{title, template.HTML(body), lnkhome, template.HTML(svcstate)}
	tmplMessage := goriceTmpl("static", "tpl/main.gtpl", "tpl/monsvc.gtpl")
	/*GO.RICE*/

	if r.Method == "GET" {
		if err := tmplMessage.Execute(w, page); err != nil {
			log.Println("Template error:", err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		/*if r.Method == "GET" {
		if err := monsvc_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}*/
	} else {
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		// Массив данных JSON для получения данных из формы (ajax)
		// Первый элемент должен содержать код действия
		type jsonPOSTData map[string]string

		var jh jsonPOSTData

		err := dec.Decode(&jh)
		if err != nil {
			log.Println("Handshake error: ", err)
		}

		enc := json.NewEncoder(w)
		switch jh["Post"] {
		case "DocumentReady" :
			st := "StatusON"
			if !MonSvcState {
				st = "StatusOFF"
			}
			enc.Encode(st)
		case "StopButton" :
			MonSvcState = false
			enc.Encode("StopOK")
		case "StartButton" :
			MonSvcState = true
			enc.Encode("StartOK")
		case "ShutdownButton":
			enc.Encode("ShutOK")
			// Завершаем работу приложения
			WaitExit = true
		default:
			//Отправляем ответ на POST-запрос
			//для предотвращения ошибки JSON parse error в ajax методе
			enc.Encode("No action requested.")
		}
	}
}



//обработчик для /config
func urlconfig(w http.ResponseWriter, r *http.Request) {
	//var accnt Acc
	//fmt.Println("method:", r.Method) //get request method
	title := "Настройка основных параметров"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	page := Page{title, template.HTML(body), lnkhome, ""}
	tmplMessage := goriceTmpl("static", "tpl/main.gtpl", "tpl/config.gtpl")
	/*GO.RICE*/

	if r.Method == "GET" {
		if err := tmplMessage.Execute(w, page); err != nil {
			log.Println("Template error:", err.Error())
			http.Error(w, http.StatusText(500), 500)
		}/*	if r.Method == "GET" {
		if err := conf_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}*/
	} else {
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		// Массив данных JSON для получения данных из формы (ajax)
		// Первый элемент должен содержать код действия:
		type jsonPOSTData map[string]string

		var jh jsonPOSTData

		err := dec.Decode(&jh)
		if err != nil {
			log.Println("Handshake error: ", err)
		}

		//JSONheader
		switch jh["post"] {
		case "SaveConfig" :
			var webport, smtpport int
			var strerr error
			if webport, strerr = strconv.Atoi(jh["web-port"]); err != nil {
				enc := json.NewEncoder(w)
				enc.Encode(strerr.Error())
				break
			}

			if smtpport, strerr = strconv.Atoi(jh["smtp-port"]); err != nil {
				enc := json.NewEncoder(w)
				enc.Encode(strerr.Error())
				break
			}
			//Заполним глобальный конфиг данными из формы
			if err := GlobalConfig.ConfigInit(GlobalConfigFile, jh["temp-dir"], jh["web-addr"], uint16(webport), jh["smtp-addr"], uint(smtpport),
				jh["smtp-login"], jh["smtp-password"], jh["from-email"], jh["from-name"], strings.Contains(jh["use-tls"], "useTLS")); err != nil {
				enc := json.NewEncoder(w)
				enc.Encode(err.Error())
				break
			}
			//Обнулим или создадим файл конфигурации
			if err := GlobalConfig.NewConfig(); err != nil {
				enc := json.NewEncoder(w)
				enc.Encode(err.Error())
				break
			}
			//Запишем данные формы в новый файл конфигурации
			if err := WriteJSONFile(&GlobalConfig); err != nil {
				enc := json.NewEncoder(w)
				enc.Encode(err.Error())
				break
			}
			enc := json.NewEncoder(w)
			enc.Encode("SaveOK")

		default:
			//Отправляем ответ на POST-запрос
			//для предотвращения ошибки JSON parse error в ajax методе
			enc := json.NewEncoder(w)
			enc.Encode("/config - no actions required")
		}
	}
}

// Возвращает tmplMessage для gorice
// сделано чисто для уменьшения текста.
func goriceTmpl(box, tplfile1, tplfile2 string) *template.Template {
	/*GO.RICE*/
	// find/create a rice.Box
	templateBox, err := rice.FindBox(box)
	if err != nil {
		log.Fatal("goriceTmpl error:", err)
	}
	// get file contents as string
	templateString, err := templateBox.String(tplfile1)
	if err != nil {
		log.Fatal("goriceTmpl error:", err)
	}
	// get file contents as string
	templateString2, err := templateBox.String(tplfile2)
	if err != nil {
		log.Fatal("goriceTmpl error:", err)
	}
	templateString += templateString2
	// parse and execute the template
	tmplMessage, err := template.New("main").Parse(templateString)
	if err != nil {
		log.Fatal("goriceTmpl error:", err)
	}
	return tmplMessage
	/*GO.RICE*/
}
