//Реализует сервис файлового монитора
package courgo

import (
	"os"
	"log"
	"fmt"
	"path/filepath"
	"errors"
	"time"
)

//Запускает выполнение правил монитора по всему списку
func StartMonitor(rules MonitorCol, accbook AddressBook, auth EmailCredentials) {
	for _, r := range rules.collection {
		fmt.Println("RULE ID =", r.id, r.mask)
		runRule(r, accbook, auth)
	}
	for _, r := range rules.collection {
		fmt.Println("RULE ID =", r.id, r.mask)
		runRule(r, accbook, auth)
	}
	fmt.Println("StartMonitor: DONE.")
}

/*
   Обработка заданного правила монитора
*/
func runRule(rule Monitor, accbook AddressBook, auth EmailCredentials) error {
	//Отображение id получателей подписки
	uid := map[uint64]bool{}
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
			uid[acc] = true
		}
	}

	if len(uid) == 0 {
		return errors.New("Не найден подходящий для правила получатель")
	}

	//Поиск файлов согласно масок, указнных в правиле монитора
	files := findFiles(rule.folder, rule.mask)
	if len(files) == 0 {
		return errors.New("Файлы указанные в правиле не найдены")
	}
	//fmt.Println(fl)
	//fmt.Println(uid)

	/*Выполнение действия согласно списка действий action
	//Обработка кодов действий
	//10 = отправить найденные вложения по email, с указанием subject и body в сообщении.
	*/
	for _, code := range rule.action {
		//Код 10 = отправка email уведомления о поступлении файла
		if code.id == 10 {
			// Убедимся что по данному файлу данным правилом не выполнялось действий

			now := time.Now()
			for name := range files {
				if GlobalHist.IsEventExist(now, rule.id, code.id, name) {
					// иначе - удалим его из списка
					if len(files) == 0 {
						return errors.New("Файлы указанные в правиле не найдены")
					}
					delete(files, name)
					continue
				}
			}
			if len(files) == 0 {
				return errors.New("Файлы указанные в правиле не найдены")
			}
			//создадим новое извещение
			msg := NewHTMLMessage(rule.msgSubject, rule.msgBody)
			//вложим все найденные файлы
			for file := range files {
				if err := msg.Attach(file); err != nil {
					log.Println("Ошибка прикрепления файла \"", file, "\":", err)
				}
			}
			//добавим всех получателей правила, для
			//которых нашлись записи в адресной книге в получатели сообщений
			for k := range uid {
				if !uid[k] {
					continue
				}
				msg.To = append(msg.To, GlobalBook.account[GlobalBook.indexByID(k)].mail...)
				msg.Body += "<br>Письмо для " + GlobalBook.account[GlobalBook.indexByID(k)].name
				//Уберем повторения адресов если таковые случатся
				msg.To = Dedup(msg.To)
			}

			// Отправим сообщение
			if err := SendEmailMsg(auth, msg); err != nil {
				log.Println("Ошибка отправки сообщения для \"", msg.To, "\":", err)
				return err
			}
			// Добавим в историю события об обработке для каждого файла
			for name, mask := range files {
				GlobalHist.AddEvt(time.Now(), rule.id, code.id, mask, true, name, msg.To)
			}
			// Записываем имеющуюся несохранненную историю на диск
			GlobalHist.Write()
		} /*if code.id == 10*/
	}
	return nil
}

func ifPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		//если путь не существует
		return false
	}
	return true
}

//Выполняет поиск файлов в каталоге согласно списка масок
func findFiles(dir string, mask []string) (files map[string]string) {
	var err error
	var list []string
	files = make(map[string]string)

	for i := range mask {
		list, err = filepath.Glob(dir + "/" + mask[i])
		if err != nil {
			log.Println("findFiles error: ", err)
			return nil
		}
		//files = append(files, list...)
		for _, f := range list {
			files[f] = mask[i]
		}
	}
	// Удаляем дубликаты из результата
	// return Dedup(files)
	return files
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
	for key := range checked {
		result = append(result, key)
	}

	return result
}

/*
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
		//if file.IsDir() {
		//	fmt.Println(file.Name()," <DIR>")
		//	continue
		//}
		fmt.Println(file)
	}
	ioutil.ReadDir(dirpath)
}
*/
