/*
Реализует структуру и методы файлового монитора

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
package courgo

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
	"errors"
	"strconv"
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

//Коллекция для хранения правил монитора
var GlobalMonCol MonitorCol

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
	desc       string
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
	Desc       string
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
	m.jsonFile = jsonfile
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

func (m *Monitor) SetDesc(desc string) {
	m.desc = desc
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
	mon.Desc = m.desc
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
	mon.desc = m.Desc
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
	col.JSONFile = m.jsonFile


	//Буфер для записи строки результата
	//удовлетворяет io.Writer
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&col)
	if err != nil {
		log.Println("MonitorCol dumpJSON error:", err)
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
	for _, elem := range col.Collection {
		m.collection = append(m.collection, elem.monJSONToMon())
	}

	return err
}

//Пишет в файл конфигурации JSON структуру
//Файл должен быть создан и доступен для записи
func (m *MonitorCol) writeJSON() (err error) {
	buffer, err := m.dumpJSON()
	if err != nil {
		log.Println("MonitorCol writeJSON error:", err)
		return err
	}
	err = ioutil.WriteFile(m.jsonFile, buffer.Bytes(), 0644)
	if err != nil {
		log.Println("MonitorCol writeJSON error:", err)
		return err
	}
	return err
}

// Выполняет возврат JSON только для collection
func (m *MonitorCol) StringJSON() (string, error) {
	// создадим аналог структуры но вместо Sid будем
	// хранить имя и номер подписчика
	type monJSONm struct {
		Id     uint64
		Desc   string
		Folder string
		Mask   []string
		Sid    []string
		Action []actionJSON
	}
	var col  []monitorJSON
	var colm monJSONm
	var colmJSON []monJSONm

	//Копируем данные JSON структуры m
	for _, elem := range m.collection {
		col = append(col, elem.monToMonJSON())
	}

	for _, elem := range col {
		colm.Id = elem.Id
		colm.Desc = elem.Desc
		colm.Folder = elem.Folder
		colm.Action = append(colm.Action, elem.Action...)
		colm.Mask = append(colm.Mask, elem.Mask...)

		for _, item := range elem.Sid {
			elm := strconv.Itoa(int(item)) + " " + GlobalBook.account[GlobalBook.indexByID(item)].name
			colm.Sid = append(colm.Sid, elm)
		}
		colmJSON = append(colmJSON, colm)
		colm.Action = []actionJSON{}
		colm.Mask = []string{}
		colm.Sid = []string{}
	}

	//Буфер для записи строки результата
	//удовлетворяет io.Writer
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&colmJSON)
	if err != nil {
		log.Println("MonitorCol StringJSON error:", err)
		return buffer.String(), err
	}
	return buffer.String(), err
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
	mon.Desc = m.desc
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
		log.Println("Monitor dumpJSON error:", err)
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
	m.desc = mon.Desc
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
		log.Println("Monitor writeJSON error:", err)
		return err
	}
	err = ioutil.WriteFile(m.jsonFile, buffer.Bytes(), 0644)
	if err != nil {
		log.Println("Monitor writeJSON error:", err)
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


// Добавляет новое правило в коллекцию правил мониторов,
// проверяет данные перед изменениями
func (m *MonitorCol) AddMonitor(newrule Monitor) error {
	if len(newrule.desc) == 0 {
		return errors.New("Название/описание правила слишком короткое.")
	}
	if len(newrule.msgSubject) == 0 {
		return errors.New("Тема извещения слишком короткая.")
	}
	if len(newrule.msgBody) == 0 {
		return errors.New("Текст извещения слишком короткий.")
	}
	if !ifPathExist(newrule.folder) {
		return errors.New("Директория наблюдения указана неверно или не существует (" + newrule.folder + ").")
	}
	if len(newrule.sid) == 0 {
		return errors.New("Не указан получатель извещения.")
	}
	if len(newrule.action) == 0 {
		return errors.New("Должно быть указано хотябы одно действие.")
	}
	if len(newrule.mask) == 0 {
		return errors.New("Должна быть указана хотябы одна файловая маска.")
	}

	for _, col := range m.collection {
		if col.id != newrule.id {
			continue
		}
		return errors.New("Правило с таким номером (" + strconv.Itoa(int(newrule.id)) + ") уже существует.")
	}

	m.collection = append(m.collection, newrule)
	return nil
}

//Возвращает наибольший номер id для account
func (m *MonitorCol) MaxID() uint64 {
	var max uint64
	for _, col := range m.collection {
		if col.id < max {
			continue
		}
		max = col.id
	}
	return max
}

//Удаляет правило с заданным id
func (m *MonitorCol) RemoveColElm(id uint64) error {
	if m.indexByID(id) == -1 {
		return errors.New("Не удалось найти заданный элемент.")
	}
	m.collection = append(m.collection[:m.indexByID(id)], m.collection[m.indexByID(id) + 1:]...)
	return nil
}

//Возвращает индекс правила монитора по его id
//если не найдено = -1
func (m *MonitorCol) indexByID(id uint64) (i int) {
	for i = 0; i < len(m.collection); i++ {
		if m.collection[i].id != id {
			continue
		}
		return i
	}
	return -1 //не найдено
}
