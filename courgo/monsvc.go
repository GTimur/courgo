/*
	Реализует сервис файлового монитора

	Реализует действия с вложениями перед отправкой
 	1. Вложить файл как есть.
 	Код действия = 10
 	2. Распаковка архивных файлов с последующим вложением содержимого в письмо.
 	Код действия = 11
*/


package courgo

import (
	"os"
	"log"
	"fmt"
	"path/filepath"
	"errors"
	"time"
	"io/ioutil"
)

//Запускает выполнение правил монитора по всему списку
func StartMonitor(rules MonitorCol, accbook AddressBook, auth EmailCredentials) {
	for _, r := range rules.collection {
		fmt.Println("RULE ID = ", r.id, "MASK = ", r.mask)
		now := time.Now()
		// Проверим настал ли новый день для регистрации событий
		if GlobalHist.IsNewDay(now) {
			// Запишем все назаписанные события на диск
			GlobalHist.Write()
			// Сотрем из памяти программы историю о старых событиях
			GlobalHist.CleanUntilDay(now)
		}
		// Запустим очередное правило
		if r.id > 999 {
			runRule(r, accbook, auth)
		}

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


	// Создадим временную директорию для различных нужд
	dir, err := ioutil.TempDir(GlobalArchi.tmp, "monsvc")
	if err != nil {
		return err
	}
	// Удалим временную директорию
	defer os.RemoveAll(dir)

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
		} /* if code.id == 10 */

		if code.id == 11 {
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
			// Cоздадим новое извещение
			msg := NewHTMLMessage(rule.msgSubject, rule.msgBody)
			// Добавим всех получателей правила, для
			// которых нашлись записи в адресной книге в получатели сообщений
			for k := range uid {
				if !uid[k] {
					continue
				}
				msg.To = append(msg.To, GlobalBook.account[GlobalBook.indexByID(k)].mail...)
				msg.Body += "<br>Письмо для " + GlobalBook.account[GlobalBook.indexByID(k)].name
				//Уберем повторения адресов если таковые случатся
				msg.To = Dedup(msg.To)
			}

			// Распакуем каждый файл и отправим его содержимое отдельным письмом.
			for file, mask := range files {
				// Распакуем файл и получим его содержимое
				funarc, err := prepUnArc(file, dir)
				if err != nil {
					return err
				}
				// Прикрепим к письму распакованные файлы
				for name := range funarc {
					if err := msg.Attach(name); err != nil {
						log.Println("Ошибка прикрепления файла \"", file, "\":", err, " Code.ID=", code.id)
					}
				}
				// Отправим сообщение
				if err := SendEmailMsg(auth, msg); err != nil {
					log.Println("Ошибка отправки сообщения для \"", msg.To, "\":", err)
					// Удалим файлы из временного каталога
					for name := range funarc {
						os.Remove(name)
					}
					return err
				}
				// Удалим файлы из временного каталога
				for name := range funarc {
					if err := os.Remove(name); err != nil {
						log.Println("Не удалось удалить файл ", name, ". Ошибка:", err)
					}
					delete(msg.Attachments, filepath.Base(name))
				}
				// Добавим в историю события об обработке для каждого файла
				GlobalHist.AddEvt(time.Now(), rule.id, code.id, mask, true, file, msg.To)
				for name, mask := range funarc {
					GlobalHist.AddEvt(time.Now(), rule.id, code.id, mask, true, name, msg.To)
				}
				// Записываем имеющуюся несохранненную историю на диск
				GlobalHist.Write()
			}
		} /* if code.id == 11 */
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
