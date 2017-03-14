// JSONConf реализует работу с конфигурационными файлами в формате JSON
// Обеспечивает сохранение и загрузку основных конфигурационных параметров
// приложения в формате JSON
package courgo

import (
	"os"
	"encoding/json"
	"log"
	"errors"
	"fmt"
	"strings"
	"strconv"
)

type Config struct {
	jsonFile   string // Путь к файлу если структура будет экспортироваться в JSON
	tmpDir     string // Директория для временных файлов
	managerSrv managerSrv
	smtpSrv    srvSMTP
}

//Глобальная переменная для хранения настроек
var GlobalConfig Config = Config{};

const GlobalConfigFile = "config.json"

type configJSON struct {
	JSONFile   string //путь к файлу если структура будет экспортироваться в JSON
	TempDir    string // Директория для временных файлов
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
	FromName string //Имя отправителя, например "Информатор GO"
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

func (c *Config) SetTempDir(tempdir string) {
	c.tmpDir = tempdir
}

func (c *Config) SetManagerSrv(addr string, port uint16) {
	c.managerSrv = managerSrv{
		Addr: addr,
		Port: port,
	}
}

func (c *Config) SetSMTPSrv(addr string, port uint, username string, password string, from string, fromname string, tls bool) {
	c.smtpSrv = srvSMTP{
		Addr: addr,
		Port: port,
		Account: username,
		Password: password,
		From: from,
		FromName:fromname,
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
		FromName:c.smtpSrv.FromName,
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
		log.Printf("Ошибка чтения JSON файла конфигурации: %v\n", err)
		return err
	}

	//Готовим для импорта структуру JSON
	var jsonConfig configJSON

	//Читаем файл JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonConfig)
	if err != nil {
		log.Printf("JSON decoder error: %v", err)
		return err
	}

	var aeserr error
	encStr := strings.Split(jsonConfig.SMTPSrv.Password, "*")

	var encBytes []byte
	b := 0
	for _, elm := range encStr {
		if len(elm) == 0 {
			continue
		}
		b, err = strconv.Atoi(elm)
		if err != nil {
			break
		}
		encBytes = append(encBytes, byte(b))
	}

	jsonConfig.SMTPSrv.Password, aeserr = AesDecript(encBytes)
	if aeserr != nil {
		log.Println("readJSON decript error:",aeserr)
		return aeserr
	}

	c.jsonFile = jsonConfig.JSONFile
	c.tmpDir = jsonConfig.TempDir
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
		log.Printf("Ошибка открытия для записи JSON файла конфигурации: %v\n", err)
		return err
	}

	//Готовим данные JSON (конвертируем в экспортируемый вид)
	jsonConfig := configJSON{
		JSONFile: c.jsonFile,
		ManagerSrv: c.managerSrv,
		SMTPSrv: c.smtpSrv,
		TempDir: c.tmpDir,
	}
	var aeserr error
	// Зашифруем строку пароля
	var encBytes []byte
	encBytes, aeserr = AesEncrypt(jsonConfig.SMTPSrv.Password)
	if aeserr != nil {
		return aeserr
	}

	str := ""
	for i := range encBytes {
		str += fmt.Sprintf("%d*", encBytes[i])
	}

	jsonConfig.SMTPSrv.Password = str

	// пишем в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&jsonConfig)
	if err != nil {
		log.Printf("JSON encoder error: %v", err)
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

// Инициализация переменной Config и проверка
// параметров на ошибки
func (c *Config) ConfigInit(jsonfile string, tempdir, webaddr string, webport uint16, smtpaddr string,
smtpport uint, user string, passwd string, emailfrom string, fromname string, usetls bool) (error) {

	// Если директории не существует
	if _, err := os.Stat(tempdir); os.IsNotExist(err) {
		return errors.New("ConfigInit:Временный каталог tempdir не существует.")
	}
	if len(jsonfile) < 1 {
		return errors.New("ConfigInit:Имя jsonfile слишком короткое.")
	}
	if len(webaddr) < 1 {
		return errors.New("ConfigInit:Адрес webaddr указан неверно.")
	}
	if webport <= 1024 {
		return errors.New("ConfigInit:Порт webport должен быть в пределах 1025-65535.")
	}
	if len(smtpaddr) < 1 {
		return errors.New("ConfigInit:Адрес smtpaddr указан неверно.")
	}
	if smtpport <= 1 {
		return errors.New("ConfigInit:Порт smtpport указан неверно.")
	}
	if len(fromname) < 1 {
		return errors.New("ConfigInit:Имя отправителя писем заполнено неверно.")
	}
	if len(user) < 1 {
		return errors.New("ConfigInit:Email адрес от имени которого ведется рассылка заполнен неверно.")
	}
	if len(emailfrom) == 0 {
		return errors.New("ConfigInit:Адрес emailfrom не может быть пустым.")
	}

	c.jsonFile = jsonfile
	c.tmpDir = tempdir
	c.managerSrv.Addr = webaddr
	c.managerSrv.Port = webport
	c.smtpSrv.Addr = smtpaddr
	c.smtpSrv.Port = smtpport
	c.smtpSrv.Account = user
	c.smtpSrv.Password = passwd
	c.smtpSrv.From = emailfrom
	c.smtpSrv.FromName = fromname
	c.smtpSrv.UseTLS = usetls

	return nil
}
