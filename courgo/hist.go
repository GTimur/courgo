// Реализует хранение сведений об обработанных файлах
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
		line = evt.Date.Format("2006-01-02 15:04:05") + "\t" +
			strconv.Itoa(int(evt.RuleID)) + "\t" +
			strconv.Itoa(int(evt.ActType)) + "\t" +
			strconv.FormatBool(evt.Completed) + "\t" +
			evt.Mask + "\t" +
			evt.File + "\t" +
			strings.Join(evt.Rcpt,",") + "\r\n"
		hst = append(hst, line)
		// Соберем индексы событий
		// для разметки признака isWritten
		idx = append(idx, i)
	}
	if len(hst[0]) == 0 {
		return err
	}
	for i, ln := range hst {
		_,err:=file.WriteString(ln)
		if err!=nil{
			log.Printf("Ошибка записи файла истории: %v\n", err)
			return err
		}
		h.Events[i].IsWritten=true
	}

	return err
}

// Проверяет есть ли событие которое выполнялось правилом
// для данного файла
func (h *Hist) IsEventExist(Date time.Time, RuleID uint64, ActType uint64, Mask string, File string) bool {
	// Проверим есть ли файлы обработанные нашим правилом
	for _,evt := range h.Events{
		if evt.File != File && evt.RuleID != RuleID && evt.Mask != Mask && evt.ActType!=ActType && !evt.Date.After(BeginOfDay(Date)){
			fmt.Println("DEBUG EVENT EXIST IN:",BeginOfDay(Date),evt.Date,evt.Date.After(BeginOfDay(Date)))
			continue
		}
		fmt.Println("DEBUG EVENT EXIST OUT:",BeginOfDay(Date),evt.Date,evt.Date.After(BeginOfDay(Date)))
		return true
	}
	return false
}


func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}