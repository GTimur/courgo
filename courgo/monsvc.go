//Реализует сервис файлового монитора
package courgo

import (
	"os"
	"log"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"errors"
)

/*
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

//Структура аккаунта подписчика для адресной книги
type Acc struct {
	id   uint64
	name string
	dept string
	mail []string
}

type AddressBook struct {
	JSONFile string //путь к файлу если структура будет экспортироваться в JSON
	account  []Acc
}


*/

//Запускает правила монитора
func StartMonitor(rules MonitorCol, accbook AddressBook) {
	for _, r := range rules.collection {
		fmt.Println("RULE ID =", r.id, r.mask)
		runRule(r, accbook)
	}
}

//Обработка заданного правила монитора
//возвращает -1 если правило содержало ошибку
func runRule(rule Monitor, accbook AddressBook) error {
	//Отображение id получателей подписки
	id := map[uint64]bool{}
	//Проверям существует ли путь указанный в правиле
	if !ifPathExist(rule.folder) {
		return errors.New("Указанный путь не существует")
	}

	//Готовим список доступных пользователей для данного правила
	//Отберем всех доступных подписчиков
	for _, acc := range rule.sid {
		for _, bac := range accbook.account {
			if bac.id != acc {
				continue
			}
			//сохраним подписчика для рассылки
			id[acc] = true
		}
	}

	if len(id)==0{
		return errors.New("Не найден подходящий получатель")
	}

	//Поиск файлов согласно масок, указнных в правиле монитора
	fl := findFiles(rule.folder, rule.mask)
	fmt.Println(fl)
	fmt.Println(id)

	return nil
}

func ifPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		//если путь не существует
		return false
	}
	return true
}

func printFileName(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return err
	}
	log.Println("Walk: ", path, " ", info.Name())
	return err
}

//Просмотр вайлов в директории
func listDir(dirpath string) {

	match, err := filepath.Glob(dirpath + "/*.*")
	// fileinfo, err := //ioutil.ReadDir(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range match {
		/*if file.IsDir() {
			fmt.Println(file.Name()," <DIR>")
			continue
		}*/
		fmt.Println(file)
	}
	ioutil.ReadDir(dirpath)
}

//Выполняет поиск файлов в каталоге согласно списка масок
func findFiles(dir string, mask []string) (files []string) {
	var err error
	var list []string

	for i := range mask {
		list, err = filepath.Glob(dir + "/" + mask[i])
		if err != nil {
			log.Println("findFiles error: ", err)
			return nil
		}

		files = append(files, list...)

	}
	//Удаляем дубликаты из результата
	return Dedup(files)
}

//Дедупликатор для среза строк
//не сохраняет порядок значений
func Dedup(slice []string) []string {
	checked := map[string]bool{}

	//Сохраним отображение без повторяющихся элементов
	for i := range slice {
		checked[slice[i]] = true
	}

	//Перенесем отображение в результат
	result := []string{}
	for key, _ := range checked {
		result = append(result, key)
	}

	return result
}