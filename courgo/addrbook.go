package courgo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

//Адресная книга содержит все аккаунты пользователей-подписчиков
type AddressBook struct {
	JSONFile string //путь к файлу если структура будет экспортироваться в JSON
	account  []Acc
}


//Адресная книга приложения
var GlobalBook AddressBook

//Заполняет глобальную книгу из указанного файла
//и меняет имя файла в поле JSONFile
func InitGlobalBook(jsonfile string) (err error) {
	GlobalBook.JSONFile = jsonfile
	if err = GlobalBook.readJSON(); err != nil {
		return err
	}
	GlobalBook.JSONFile = jsonfile
	return err
}


//Возвращает наибольший номер id для account
func (a *AddressBook) MaxID() uint64 {
	var max uint64
	for i := 0; i < len(a.account); i++ {
		if a.account[i].id < max {
			continue
		}
		max = a.account[i].id
	}
	return max
}

//Удаляет аккаунт с заданным id
func (a *AddressBook) RemoveAccount(id uint64) {
	i := a.indexByID(id) //index of record by id
	a.account = append(a.account[:i],a.account[i+1:]...)
}

//Возвращает индекс подписчика (Acc) по его id
//если не найдено = -1
func (a *AddressBook) indexByID(id uint64) (i int) {
	for i = 0; i < len(a.account); i++ {
		if a.account[i].id != id {
			continue
		}
		return i
	}
	return -1 //не найдено
}


//Возвращает копию адресной книги
func (a AddressBook) Copy() AddressBook {
	return a
}

//Добавляет аккаунт в адресную книгу
func (a *AddressBook) Add(account Acc) int {
	a.account = append(a.account, account)
	return len(a.account)
}

//Дамп адресной книги в JSON
//если id записи = 0 выполняется дамп всей книги
//возвращаемый указатель удовлетворяет io.Writer
func (a *AddressBook) dumpJSON(id uint64) (*bytes.Buffer, error) {
	//Структуры для эскпорта данных в JSON
	type AccJSON struct {
		Id   uint64
		Name string
		Dept string
		Mail []string
	}
	type addressBookJSON struct {
		Account []AccJSON
	}
	// Экземпляр структуры с именами полей для выгрузки JSON
	var JSONBook addressBookJSON

	//Копируем поля Aсс в поля AccJSON
	//конвертация из неэкспортируемой структуры Acc
	for i, elem := range a.account {
		if id > 0 && a.account[i].id != id {
			continue
		}
		//Добавляем новый эелемент
		JSONBook.Account = append(JSONBook.Account, AccJSON{
			Id:elem.id,
			Name:elem.name,
			Dept:elem.dept,
			Mail:elem.mail,
		})
	}

	//Буфер для записи строки результата
	//удовлетворяет io.Writer
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&JSONBook)
	if err != nil {
		log.Println(err)
		return buffer, err
	}
	return buffer, err
}

//Возвращает запись JSON из адресной книги в string
//если id записи = 0 выполняется дамп всей книги
func (a *AddressBook) StringJSON(id uint64) (string, error) {
	var err error
	buffer := bytes.NewBufferString("")
	if buffer, err = a.dumpJSON(id); err != nil {
		log.Println(err)
		return "", err
	}
	return buffer.String(), err
}

//Возвращает запись JSON из адресной книги в buffer
//если id записи = 0 выполняется дамп всей книги
//возвращаемый указатель удовлетворяет io.Writer
func (a *AddressBook) BufferJSON(id uint64) (*bytes.Buffer, error) {
	var err error
	buffer := bytes.NewBufferString("")
	if buffer, err = a.dumpJSON(id); err != nil {
		log.Println(err)
		return nil, err
	}
	return buffer, err
}

//Запись данных книги в файл JSON
//файл будет перезаписан
func (a *AddressBook) writeJSON() (error) {
	buffer, err := a.dumpJSON(0)
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(a.JSONFile, buffer.Bytes(), 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

//Заполнение (чтение) адресной книги из JSON конфигурации
//возвращаемый указатель удовлетворяет io.Writer
func (a *AddressBook) readJSON() (error) {
	//Структуры для эскпорта данных в JSON
	type AccJSON struct {
		Id   uint64
		Name string
		Dept string
		Mail []string
	}
	type addressBookJSON struct {
		Account []AccJSON
	}
	// Экземпляр структуры с именами полей для выгрузки JSON
	var JSONBook addressBookJSON

	buffer, err := ioutil.ReadFile(a.JSONFile)
	if err != nil {
		log.Println(err)
		return err
	}

	//Буфер для записи строки результата удовлетворяет io.Reader
	decoder := json.NewDecoder(bytes.NewReader(buffer))
	err = decoder.Decode(&JSONBook)
	if err != nil {
		log.Println(err)
		return err

	}

	//Копируем данные из addressBookJSON в AddressBook
	//перенос импортируемых данных в приватную структуру
	for _, elem := range JSONBook.Account {
		/*if id > 0 && JSONBook.Account[i].id != id {
			continue
		}*/
		//Добавляем новый эелемент
		a.account = append(a.account, Acc{
			id:elem.Id,
			name:elem.Name,
			dept:elem.Dept,
			mail:elem.Mail,
		})
	}

	//res,_:=a.StringJSON(0)
	//fmt.Println("DEBUG:",res)
	return err
}
//buf := bytes.NewBuffer(make([]byte, 0, capacity))