package courgo

import (
	"os"
	"fmt"
	"encoding/json"
	"log"
)

type Config struct {
	JSONFile   string //путь к файлу если структура будет экспортироваться в JSON
	PTKPath    PTKPSD
	ManagerSrv ManagerSrv
}

//ПТК ПСД основные параметры
type PTKPSD struct {
	PathIN  string //Путь к файлам исходящих сообщений
	PathOUT string //Путь к файлам входящих сообщений
}

//Адрес сервера управления программой и порт
type ManagerSrv struct {
	Addr string
	Port uint64
}

type readerWriterJSON interface {
	readJSON() error
	writeJSON() error
}


//Создает новый файл для записи конфигурации JSON если
//таковой отсутствует (не перезаписывает его)
func (c *Config) MakeConfig() (err error) {
	if _, err = os.Stat(c.JSONFile); err == nil {
		//Файл существует и не будет перезаписан
		//fmt.Println(os.IsExist(err),err)
		return err
	}
	file, err := os.Create(c.JSONFile)
	defer file.Close()
	return err
}

//Создает новый файл для записи конфигурации JSON
//стирает данные если файл существует
func (c *Config) NewConfig() (err error) {
	if _, err = os.Stat(c.JSONFile); err == nil {
		//Файл существует и не будет перезаписан
		//fmt.Println(os.IsExist(err),err)
	}
	file, err := os.Create(c.JSONFile)
	defer file.Close()
	return err
}

//Читает из файла конфигурации JSON-структуру
func (c *Config) readJSON() (err error) {
	file, err := os.Open(c.JSONFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("Ошибка чтения JSON файла конфигурации: %v\n", err)
		return err
	}
	//Читаем файл JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		fmt.Errorf("JSON decoder error: %v", err)
		//return err
	}
	fmt.Printf("%T %T\n", c.ManagerSrv.Addr, c.ManagerSrv.Port)
	fmt.Printf("%s:%d\n", c.ManagerSrv.Addr, c.ManagerSrv.Port)
	return err
}

func (c *Config) ReadJSON() (err error) {
	err = c.readJSON()
	return
}


//Пишет в файл конфигурации JSON структуру
//Файл должен быть создан и доступен для записи
func (c *Config) writeJSON() (err error) {
	file, err := os.OpenFile(c.JSONFile, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Fatalf("Ошибка открытия для записи JSON файла конфигурации: %v\n", err)
		return err
	}
	//Готовим данные JSON и пишем в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&c)
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
