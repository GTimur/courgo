//Web-сервер реализует интерфейс управления приложением
package courgo

import (
	"net"
	"net/http"
	"log"
	"fmt"
	"html/template"
	"github.com/braintree/manners"
	"path"
	"strings"
	"encoding/json"
)

type WebCtl struct {
	host     net.IP
	port     uint16
	islisten bool
}

type Page struct {
	Title string
	Body  template.HTML
	LnkHome string
}

var (
	// компилируем шаблоны
	home_template = template.Must(template.ParseFiles("main.gtpl", path.Join("static", "tpl", "index.gtpl")))
	acc_template = template.Must(template.ParseFiles("main.gtpl", path.Join("static", "tpl", "acc.gtpl")))
	register_template = template.Must(template.ParseFiles("main.gtpl", path.Join("static", "tpl", "register.gtpl")))
	monitor_template = template.Must(template.ParseFiles("main.gtpl", path.Join("static", "tpl", "mon.gtpl")))
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
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", urlhome) //Каждый запрос вызывает обработчик
	http.HandleFunc("/acc", urlacc) //Страница с таблицей подписчиков
	http.HandleFunc("/acc/register", urlregister) //Страница регистрации подписчика
	http.HandleFunc("/mon", urlmon) //Страница с таблицей правил монитора
	go func() {
		log.Fatal(manners.ListenAndServe(w.connString(), http.DefaultServeMux))
	}()
	w.islisten = true
	return err
}

//Обработчик запросов для home
func urlhome(w http.ResponseWriter, r *http.Request) {
	title := "COURIER GO"
	body := ""
	lnkhome := "http://127.0.0.1:8000"
	page := Page{title, template.HTML(body), lnkhome }
	if err := home_template.ExecuteTemplate(w, "main", page); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

//обработчик для /acc/register
func urlregister(w http.ResponseWriter, r *http.Request) {
	//var accnt Acc
	//fmt.Println("method:", r.Method) //get request method
	title := "Регистрация подписчика"
	body := ""
	lnkhome := "http://127.0.0.1:8000"
	page := Page{title, template.HTML(body),lnkhome}
	if r.Method == "GET" {
		if err := register_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
	} else {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		if strings.Contains(r.Form.Encode(), "cancelbutton") {
			fmt.Println(r.Form.Encode())
		} else {
			log.Println("POST")
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

			switch jh["post"] {
			case "SaveButton" :

				if err:=RegisterAccount(jh["fio"],jh["dept"],strings.Split(jh["email"], ","));err!=nil{
					log.Println("Ошибка регистрации подписчика:",err)
					enc := json.NewEncoder(w)
					enc.Encode("SaveNotOk")
				} else {
					if err := WriteJSONFile(&GlobalBook); err != nil {
						log.Println("Не удалось сохранить в файл нового подписчика:", err)
						enc := json.NewEncoder(w)
						enc.Encode("SaveNotOk")
					}else{
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
	lnkhome:="http://127.0.0.1:8000"
	page := Page{title, template.HTML(body),lnkhome}
	if r.Method == "GET" {
		if err := acc_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
	} else {
		log.Println("POST")
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
		log.Println(jh[0])

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
	//var accnt Acc
	//fmt.Println("method:", r.Method) //get request method
	title := "Управление правилами монитора"
	body := ""
	lnkhome:="http://127.0.0.1:8000"
	page := Page{title, template.HTML(body),lnkhome}
	if r.Method == "GET" {
		if err := monitor_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
	} else {
		log.Println("POST")
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		// Массив данных JSON для получения данных из формы (ajax)
		// Первый элемент должен содержать код действия
		type jsonPOSTData []string

		var jh jsonPOSTData

		err := dec.Decode(&jh)
		if err != nil {
			log.Println("Handshake error: ", err)
		}
		log.Println(jh[0])

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
	}
}