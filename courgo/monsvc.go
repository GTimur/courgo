/*
	Реализует сервис файлового монитора

	Реализует действия с вложениями перед отправкой
 	1. Вложить файл как есть.
 	Код действия = 10
 	2. Распаковка архивных файлов с последующим вложением содержимого в письмо.
 	Код действия = 11
 	3. Не прикреплять вложение к письму, только уведомление
 	Код действия = 20
 	4. Исключить из обработки
 	Код действия = 1000
*/

package courgo

import (
	"os"
	"log"
	"path/filepath"
	"errors"
	"time"
	"io/ioutil"
	"fmt"
	"strings"
)

// Переменная определяющая состояние обработчика правил монитора.
var MonSvcState bool

// Запускает выполнение правил монитора по всему списку
// если state = false, тогда обработка не выполняется
func StartMonitor(rules MonitorCol, accbook AddressBook, auth EmailCredentials, state bool) error {
	evts := len(GlobalHist.Events)
	// Проверим состояние выключателя обработки (до начала работы)
	if !state {
		return nil
	}
	for _, r := range rules.collection {
		now := time.Now()
		// Временно пусть это тут побудет (сдвиг на 3 часа)
		now = now.Add(-3 * time.Hour)
		// Проверим настал ли новый день для регистрации событий
		if GlobalHist.IsNewDay(now) {
			// Запишем все назаписанные события на диск
			if err := GlobalHist.Write(); err != nil {
				return err
			}
			// Сотрем из памяти программы историю о старых событиях
			if err := GlobalHist.CleanUntilDay(now); err != nil {
				return err
			}
		}
		// Проверим состояние выключателя обработки (для отключения во время работы)
		if !MonSvcState {
			return nil
		}
		//fmt.Println("DIAG:", r.id, r.mask, r.sid)
		// Запустим очередное правило
		if r.id > 0 {
			if err := runRule(r, accbook, auth); err != nil {
				//fmt.Println("StartMonitor:", err)
			}
		}
	}
	// Если Количество событий изменилось - сохраним на диск (JSON)
	if evts == len(GlobalHist.Events) {
		return nil
	}
	// Запишем все назаписанные события на диск
	if err := GlobalHist.RewriteJSON(); err != nil {
		return err
	}
	return nil
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
		return errors.New("Поиск файлов: Файлы указанные в правиле не найдены. " + fmt.Sprintf("%s", rule.mask))
	}

	//fmt.Println(files, rule.mask)
	// fmt.Println(uid)

	// Создадим временную директорию для различных нужд
	dir, err := ioutil.TempDir(GlobalConfig.tmpDir, "monsvc")
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
			msg.Body += "<br><br>Направлено:<br>"
			for k := range uid {
				if !uid[k] {
					continue
				}
				msg.To = append(msg.To, GlobalBook.account[GlobalBook.indexByID(k)].mail...)
				msg.Body += GlobalBook.account[GlobalBook.indexByID(k)].name + "<br>"
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
			fmt.Println("Completed: RULE ", rule.id, " FILES:", files, time.Now().Format("02/01/2006 15:04:05"))
			// Записываем имеющуюся несохранненную историю на диск
			if err := GlobalHist.Write(); err != nil {
				return err
			}
			if err := GlobalHist.RewriteJSON(); err != nil {
				return err
			}
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
			msg.Body += "<br><br>Направлено:<br>"
			for k := range uid {
				if !uid[k] {
					continue
				}
				msg.To = append(msg.To, GlobalBook.account[GlobalBook.indexByID(k)].mail...)
				msg.Body += GlobalBook.account[GlobalBook.indexByID(k)].name + "<br>"
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
				// Добавим в историю событие об обработке самого архива
				GlobalHist.AddEvt(time.Now(), rule.id, code.id, mask, true, file, msg.To)
				// Добавим в историю события об обработке для каждого файла
				for name, mask := range funarc {
					GlobalHist.AddEvt(time.Now(), rule.id, code.id, mask, true, name, msg.To)
				}
				fmt.Println("Completed: RULE ", rule.id, " FILES:", files, time.Now().Format("02/01/2006 15:04:05"))
				// Записываем имеющуюся несохранненную историю на диск
				if err := GlobalHist.Write(); err != nil {
					return err
				}
				if err := GlobalHist.RewriteJSON(); err != nil {
					return err
				}
			}
		} /* if code.id == 11 */

		//Код 20 = отправка email уведомления без вложений
		if code.id == 20 {
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

			//добавим всех получателей правила, для
			//которых нашлись записи в адресной книге в получатели сообщений
			msg.Body += "<br><br>Направлено:<br>"
			for k := range uid {
				if !uid[k] {
					continue
				}
				msg.To = append(msg.To, GlobalBook.account[GlobalBook.indexByID(k)].mail...)
				msg.Body += GlobalBook.account[GlobalBook.indexByID(k)].name + "<br>"
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
			fmt.Println("Completed: RULE ", rule.id, " FILES:", files, time.Now().Format("02/01/2006 15:04:05"))
			// Записываем имеющуюся несохранненную историю на диск
			if err := GlobalHist.Write(); err != nil {
				return err
			}
			if err := GlobalHist.RewriteJSON(); err != nil {
				return err
			}
		} /* if code.id == 20 */
		if code.id == 1000 {
			fmt.Println("Completed: RULE ", rule.id)
		} /* if code.id == 1000 */
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
		list, err = filepath.Glob(dir + "\\" + strings.ToUpper(mask[i]))
		if err != nil {
			log.Println("findFiles error: ", err)
			return nil
		}
		//files = append(files, list...)
		for _, f := range list {
			files[f] = mask[i]
		}
	}
	for i := range mask {
		list, err = filepath.Glob(dir + "\\" + strings.ToLower(mask[i]))
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