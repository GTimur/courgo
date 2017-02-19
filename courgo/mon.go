//Реализует структуру и методы файлового монитора
package courgo

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
)

/* Folders Monitor

  Monitor -> Monitor search rules (Monitor Search Rule ID = MSRID)
          -> Check ALL FOLDERS from MSRID's
          -> Send notification to subscribers from address book (Subscriber ID = SID)
          -> Do possible Actions (Action ID = AID)

  Monitor RULE	-> Folder
              	-> Mask
                -> Subscribers (SID) []
                -> Notification e-mail message text header and body
                -> Reader (RID)
                -> Actions (AID) []

*/

//Структура для коллекции всех активных мониторов
type MonitorCol struct {
	jsonFile   string
	collection []Monitor
}

//Аналог MonitorCol для конвертации из/в JSON
type monitorColJSON struct {
	JSONFile   string
	Collection []monitorJSON
}

//Структра для описания одного экземпляра правила монитора
type Monitor struct {
	id         uint64   //Monitor rule id
	folder     string
	mask       []string
	sid        []uint64 //subscriber id
	msgSubject string
	msgBody    string
	jsonFile   string
	action     []Action //action id
}

//Аналог Monitor для конвертации из/в JSON
type monitorJSON struct {
	Id         uint64
	Folder     string
	Mask       []string
	Sid        []uint64 //subscriber id
	MsgSubject string
	MsgBody    string
	JSONFile   string
	Action     []actionJSON
}

func (m *MonitorCol) SetCollection(collection []Monitor) {
	for _, elm := range collection {
		m.collection = append(m.collection, elm)
	}
}

func (m *MonitorCol) SetJSONFile(jsonfile string) {
	m.jsonFile=jsonfile
}

//Структура для описания одного действия
type Action struct {
	id   uint64 //id действия
	desc string //описание действия
}

//Аналог Action для конвертации из/в JSON
type actionJSON struct {
	Id   uint64 //id действия
	Desc string //описание действия
}

func (m *Monitor) SetId(id uint64) {
	m.id = id
}

func (m *Monitor) SetFolder(folder string) {
	m.folder = folder
}

func (m *Monitor) SetMask(mask []string) {
	m.mask = mask
}

func (m *Monitor) SetSID(SID []uint64) {
	m.sid = SID
}

func (m *Monitor) SetMsgSubject(subj string) {
	m.msgSubject = subj
}

func (m *Monitor) SetMsgBody(body string) {
	m.msgBody = body
}

func (m *Monitor) SetJSONFile(JSONfile string) {
	m.jsonFile = JSONfile
}

func (m *Monitor) SetAction(action []Action) {
	m.action = action
}
/*Action*/
func (a *Action) SetId(id uint64) {
	a.id = id
}

func (a *Action) SetDesc(desc string) {
	a.desc = desc
}

//Создает новый файл для записи конфигурации JSON если
//таковой отсутствует (не перезаписывает его)
func (m *MonitorCol) MakeConfig() (err error) {
	if _, err = os.Stat(m.jsonFile); err == nil {
		return err
	}
	file, err := os.Create(m.jsonFile)
	defer file.Close()
	return err
}

//Создает новый файл для записи конфигурации JSON
//стирает данные если файл существует
func (m *MonitorCol) NewConfig() (err error) {
	if _, err = os.Stat(m.jsonFile); err == nil {
		//Файл существует и не будет перезаписан
		//fmt.Println(os.IsExist(err),err)
	}
	file, err := os.Create(m.jsonFile)
	defer file.Close()
	return err
}


//Конвертирует Monitor в monitorJSON
func (m *Monitor) monToMonJSON() monitorJSON {
	var mon monitorJSON

	mon.JSONFile = m.jsonFile
	mon.Folder = m.folder
	mon.Mask = m.mask
	mon.MsgSubject = m.msgSubject
	mon.MsgBody = m.msgBody
	mon.Sid = m.sid
	mon.Id = m.id
	for i := range m.action {
		mon.Action = append(mon.Action, actionJSON{m.action[i].id, m.action[i].desc})
	}

	return mon
}

//Конвертирует monitorJSON в Monitor
func (m *monitorJSON) monJSONToMon() Monitor {
	var mon Monitor

	mon.jsonFile = m.JSONFile
	mon.folder = m.Folder
	mon.mask = m.Mask
	mon.msgSubject = m.MsgSubject
	mon.msgBody = m.MsgBody
	mon.sid = m.Sid
	mon.id = m.Id
	for i := range m.Action {
		mon.action = append(mon.action, Action{m.Action[i].Id, m.Action[i].Desc})
	}

	return mon
}


//Дамп структуры мониотра в JSON
//возвращаемый указатель удовлетворяет io.Writer
func (m *MonitorCol) dumpJSON() (*bytes.Buffer, error) {
	//экземпляр структуры для иморта в JSON
	var col monitorColJSON

	for _, elem := range m.collection {
		col.Collection = append(col.Collection, elem.monToMonJSON())
	}
	col.JSONFile=m.jsonFile

	//Буфер для записи строки результата
	//удовлетворяет io.Writer
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&col)
	if err != nil {
		log.Println(err)
		return buffer, err
	}
	return buffer, err
}

//Читает из файла в JSON-структуру
func (m *MonitorCol) readJSON() (err error) {
	//экземпляр структуры для иморта в JSON
	var col monitorColJSON = monitorColJSON{}
	col.JSONFile = m.jsonFile
	file, err := os.Open(col.JSONFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("Коллекция монитора: ошибка чтения JSON файла: %v\n", err)
		return err
	}
	//Читаем файл JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&col)
	if err != nil {
		fmt.Errorf("JSON decoder error: %v", err)
	}
	//Копируем данные JSON структуры в m
	for _, elem := range col.Collection{
		m.collection = append(m.collection, elem.monJSONToMon())
	}

	return err
}

//Пишет в файл конфигурации JSON структуру
//Файл должен быть создан и доступен для записи
func (m *MonitorCol) writeJSON() (err error) {
	buffer, err := m.dumpJSON()
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(m.jsonFile, buffer.Bytes(), 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}


//Создает новый файл для записи конфигурации JSON если
//таковой отсутствует (не перезаписывает его)
func (m *Monitor) MakeConfig() (err error) {
	if _, err = os.Stat(m.jsonFile); err == nil {
		return err
	}
	file, err := os.Create(m.jsonFile)
	defer file.Close()
	return err
}

//Создает новый файл для записи конфигурации JSON
//стирает данные если файл существует
func (m *Monitor) NewConfig() (err error) {
	if _, err = os.Stat(m.jsonFile); err == nil {
		//Файл существует и не будет перезаписан
		//fmt.Println(os.IsExist(err),err)
	}
	file, err := os.Create(m.jsonFile)
	defer file.Close()
	return err
}

//Дамп структуры мониотра в JSON
//возвращаемый указатель удовлетворяет io.Writer
func (m *Monitor) dumpJSON() (*bytes.Buffer, error) {
	//экземпляр структуры для иморта в JSON
	var mon monitorJSON = monitorJSON{}
	var act []actionJSON = []actionJSON{}

	for _, elem := range m.action {
		act = append(act, actionJSON{elem.id, elem.desc})
	}

	mon.JSONFile = m.jsonFile
	mon.Folder = m.folder
	mon.Mask = m.mask
	mon.MsgSubject = m.msgSubject
	mon.MsgBody = m.msgBody
	mon.Sid = m.sid
	mon.Id = m.id
	mon.Action = act

	//Буфер для записи строки результата
	//удовлетворяет io.Writer
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&mon)
	if err != nil {
		log.Println(err)
		return buffer, err
	}
	return buffer, err
}



//Читает из файла в JSON-структуру
func (m *Monitor) readJSON() (err error) {
	//экземпляр структуры для иморта в JSON
	var mon monitorJSON = monitorJSON{}
	mon.JSONFile = m.jsonFile
	file, err := os.Open(mon.JSONFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("Монитор: ошибка чтения JSON файла: %v\n", err)
		return err
	}
	//Читаем файл JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&mon)
	if err != nil {
		fmt.Errorf("JSON decoder error: %v", err)
		//return err
	}

	var act []Action = []Action{}

	for _, elem := range mon.Action {
		act = append(act, Action{elem.Id, elem.Desc})
	}

	m.folder = mon.Folder
	m.mask = mon.Mask
	m.msgSubject = mon.MsgSubject
	m.msgBody = mon.MsgBody
	m.sid = mon.Sid
	m.id = mon.Id
	m.action = act
	return err
}

func (m *Monitor) ReadJSON() (err error) {
	err = m.readJSON()
	return
}

//Пишет в файл конфигурации JSON структуру
//Файл должен быть создан и доступен для записи
func (m *Monitor) writeJSON() (err error) {
	buffer, err := m.dumpJSON()
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(m.jsonFile, buffer.Bytes(), 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

//Функция записи данных в формате JSON в файл
//путь для файла будет взят из поля jsonFile структуры
func (m *Monitor) WriteJSON() (err error) {
	err = m.writeJSON()
	return
}
