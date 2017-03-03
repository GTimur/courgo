/* Реализует хранение сведений об обработанных файлах
- История обработанных сообщений сохраняется в процессе работы программы в файл для просмотра статистики работы программы.
- Оперативная история обработанных сообщений хранится в памяти приложения и очищается в начале дня. (Окно хранения - сутки.)
История состав:
	Дата/время,
	Номер правила,
	Тип действия (ID),
	Статус (обработан, ожидается дальнейшая обработка),
	Маска,
	Полное имя файла включая путь,
	Список получателей (для кого выполнялось действие),
	Признак сброса данных в файл
*/
package courgo

import (
	"time"
	"os"
	"log"
	"strconv"
	"strings"
	"fmt"
)

type Hist struct {
	Filename string
	Events   []Event
}

type Event struct {
	Date      time.Time
	RuleID    uint64
	//Тип действия (ID)
	ActType   uint64
	// Признак завершения обработки правилом
	Completed bool
	Mask      string
	File      string
	//Получатели
	Rcpt      []string
	// Признак записи события в файл
	IsWritten bool
}

var GlobalHist Hist

func (h *Hist) SetFilename(filename string) {
	h.Filename = filename
}

func (h *Hist) AddEvt(Date time.Time, RuleID uint64, ActType uint64, Mask string, Completed bool, File string, Rcpt []string) {
	evt := Event{Date:Date, RuleID:RuleID, ActType:ActType, Mask:Mask, Completed:Completed, File:File, Rcpt:Rcpt}
	h.Events = append(h.Events, evt)
}

//Создает новый файл для записи файла истории
//если таковой отсутствует (не перезаписывает его)
func (h *Hist) MakeHistFile() (err error) {
	if _, err = os.Stat(h.Filename); err == nil {
		//Файл существует и не будет перезаписан
		//fmt.Println(os.IsExist(err),err)
		return err
	}
	file, err := os.Create(h.Filename)
	defer file.Close()
	return err
}

// Добавляет историю в файл
func (h *Hist) Write() (err error) {
	if err := h.MakeHistFile(); err != nil {
		return err
	}
	file, err := os.OpenFile(h.Filename, os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("Ошибка записи файла истории: %v\n", err)
		return err
	}
	var hst []string
	var line string
	var idx []int
	for i, evt := range h.Events {
		if evt.IsWritten {
			continue
		}
		// Сохраняем в историю события только из текущей даты (с 00 часов текущего дня)
		if evt.Date.Before(BeginOfDay(time.Now())) {
			continue
		}
		line = evt.Date.Format("2006-01-02 15:04:05") + "\t" +
			strconv.Itoa(int(evt.RuleID)) + "\t" +
			strconv.Itoa(int(evt.ActType)) + "\t" +
			strconv.FormatBool(evt.Completed) + "\t" +
			evt.Mask + "\t" +
			evt.File + "\t" +
			strings.Join(evt.Rcpt, ",") + "\r\n"
		hst = append(hst, line)
		// Соберем индексы событий
		// для разметки признака isWritten
		idx = append(idx, i)
	}
	if len(hst) == 0 {
		return err
	}
	for i, ln := range hst {
		_, err := file.WriteString(ln)
		if err != nil {
			log.Printf("Ошибка записи файла истории: %v\n", err)
			return err
		}
		h.Events[i].IsWritten = true
	}

	return err
}

// Проверяет есть ли событие которое выполнялось правилом
// для данного файла (все события актуальны в течение суточного окна)
func (h *Hist) IsEventExist(Date time.Time, RuleID uint64, ActType uint64, File string) bool {
	// Проверим есть ли файлы обработанные нашим правилом
	for _, evt := range h.Events {
		//fmt.Println("====:", evt.File, " ", File, " ", evt.RuleID, " ", RuleID, " ", evt.Date, " ", BeginOfDay(Date))
		if evt.Date.Before(BeginOfDay(Date)) {
			continue
		}
		if strings.Compare(evt.File, File) != 0 || evt.RuleID != RuleID || evt.ActType != ActType {
			continue
		}
		fmt.Println("Event already exist.")
		return true
	}
	return false
}

func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// Удаляет всю историю из памяти до 0 часов указаннго дня Day.
func (h *Hist) CleanUntilDay(Day time.Time) error {
	for i, evt := range h.Events {
		if !evt.Date.Before(BeginOfDay(Day)) {
			continue
		}
		// Удалим событие из списка если его дата меньше указанной
		if len(h.Events) == 0 {
			return nil // нет данных для обработки
		}
		h.Events = append(h.Events[:i], h.Events[i + 1:]...)
	}
	return nil
}

// Проверяет настал ли новый операционный день.
// Если последнее событие было до начала текущего дня - новый день настал.
func (h *Hist) IsNewDay(Day time.Time) bool {
	if len(h.Events) == 0 {
		return false // нет данных для обработки
	}
	// Последнее событие датировано не раньше начала текущего дня
	if !h.Events[len(h.Events) - 1].Date.Before(Day) {
		return false
	}
	// Последнее событие датировано РАНЬШЕ начала текущего дня
	return true
}