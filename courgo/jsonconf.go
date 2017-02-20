// JSONConf реализует работу с конфигурационными файлами в формате JSON
// Обеспечивает сохранение и загрузку основных конфигурационных параметров
// приложения в формате JSON
package courgo

import (
	"os"
	"fmt"
	"encoding/json"
	"log"
)

type Config struct {
	jsonFile   string //путь к файлу если структура будет экспортироваться в JSON
	managerSrv managerSrv
	smtpSrv    srvSMTP
}

type configJSON struct {
	JSONFile   string //путь к файлу если структура будет экспортироваться в JSON
	ManagerSrv managerSrv
	SMTPSrv    srvSMTP
}

//Представляет сведения о настройках SMTP
type srvSMTP struct {
	Addr     string
	Port     uint
	Account  string //smtp аккаунт для отправки сообщений
	Password string
	From     string //Адрес отправителя
	UseTLS   bool   //auth: использовать TLS или plain/text
}

//Представляет адрес сервера управления программой и порт
type managerSrv struct {
	Addr string
	Port uint16
}

type readerWriterJSON interface {
	readJSON() error
	writeJSON() error
}


/*SETTERS*/
func (c *Config) SetJSONFile(jsonfile string) {
	c.jsonFile = jsonfile
}

func (c *Config) SetManagerSrv(addr string, port uint16) {
	c.managerSrv = managerSrv{
		Addr: addr,
		Port: port,
	}
}

func (c *Config) SetSMTPSrv(addr string, port uint, username string, password string, from string, tls bool) {
	c.smtpSrv = srvSMTP{
		Addr: addr,
		Port: port,
		Account: username,
		Password: password,
		From: from,
		UseTLS: tls,
	}
}

/*GETTERS*/
func (c *Config) JSONFile() string {
	return c.jsonFile
}

func (c *Config) ManagerSrvAddr() string {
	return c.managerSrv.Addr
}

func (c *Config) ManagerSrvPort() uint16 {
	return c.managerSrv.Port
}

//Возвращает данные smtpSrv в виде EmailCredentials
func (c *Config) SMTPCred() EmailCredentials {
	creds := EmailCredentials{
		Server: c.smtpSrv.Addr,
		Port:int(c.smtpSrv.Port),
		Username:c.smtpSrv.Account,
		Password:c.smtpSrv.Password,
		From:c.smtpSrv.From,
		UseTLS:c.smtpSrv.UseTLS,
	}
	return creds
}

func (c *Config) SMTPSrvAddr() string {
	return c.smtpSrv.Addr
}

func (c *Config) SMTPSrvPort() uint {
	return c.smtpSrv.Port
}

func (c *Config) SMTPSrvUser() string {
	return c.smtpSrv.Account
}

func (c *Config) SMTPSrvPassword() string {
	return c.smtpSrv.Password
}

func (c *Config) SMTPSrvUseTLS() bool {
	return c.smtpSrv.UseTLS
}


//Создает новый файл для записи конфигурации JSON если
//таковой отсутствует (не перезаписывает его)
func (c *Config) MakeConfig() (err error) {
	if _, err = os.Stat(c.jsonFile); err == nil {
		//Файл существует и не будет перезаписан
		//fmt.Println(os.IsExist(err),err)
		return err
	}
	file, err := os.Create(c.jsonFile)
	defer file.Close()
	return err
}

//Создает новый файл для записи конфигурации JSON
//стирает данные если файл существует
func (c *Config) NewConfig() (err error) {
	if _, err = os.Stat(c.jsonFile); err == nil {
		//Файл существует и не будет перезаписан
		//fmt.Println(os.IsExist(err),err)
	}
	file, err := os.Create(c.jsonFile)
	defer file.Close()
	return err
}

//Читает из файла конфигурации JSON-структуру
func (c *Config) readJSON() (err error) {
	file, err := os.Open(c.jsonFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("Ошибка чтения JSON файла конфигурации: %v\n", err)
		return err
	}
	//Готовим для импорта структуру JSON
	var jsonConfig configJSON

	//Читаем файл JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonConfig)
	if err != nil {
		fmt.Errorf("JSON decoder error: %v", err)
		//return err
	}
	c.jsonFile = jsonConfig.JSONFile
	c.managerSrv = jsonConfig.ManagerSrv
	c.smtpSrv = jsonConfig.SMTPSrv

	//fmt.Println("readJSON: ", c)
	return err
}

func (c *Config) ReadJSON() (err error) {
	err = c.readJSON()
	return
}


//Пишет в файл конфигурации JSON структуру
//Файл должен быть создан и доступен для записи
func (c *Config) writeJSON() (err error) {
	file, err := os.OpenFile(c.jsonFile, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Fatalf("Ошибка открытия для записи JSON файла конфигурации: %v\n", err)
		return err
	}

	//Готовим данные JSON (конвертируем в экспортируемый вид)
	jsonConfig := configJSON{
		JSONFile: c.jsonFile,
		ManagerSrv: c.managerSrv,
		SMTPSrv: c.smtpSrv,
	}

	// пишем в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&jsonConfig)
	if err != nil {
		log.Fatalf("JSON encoder error: %v", err)
		//return err
	}
	return err
}

//Функция записи данных в формате JSON в файл
//путь для файла будет взят из поля JSONFile структуры
func (c *Config) WriteJSON() (err error) {
	err = c.writeJSON()
	return
}

//Интерфейсная функция чтения данных JSON в структурную переменную
//путь для файла будет взят из поля JSONFile структуры
func ReadJSONFile(data readerWriterJSON) error {
	err := data.readJSON()
	return err
}

//Интерфейсная функция записи данных в формате JSON в файл
//путь для файла будет взят из поля JSONFile структуры
func WriteJSONFile(data readerWriterJSON) error {
	err := data.writeJSON()
	return err
}

//configrw := Config{PTKPSD{PathIN:"Y:\\test\\IN",PathOUT:"Y:\\test\\OUT"},managersrv{Addr:"127.0.0.1",Port:9090},}
