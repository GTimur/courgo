//Реализует хранение сведений об обработанных файлах
package courgo

import (
	"time"
	"os"
	"log"
)

type Hist struct {
	Filename string
	Events   []Event
}

type Event struct {
	Date      time.Time
	RuleID    uint
	// Признак завершения обработки = true
	Completed bool
	Mask      string
	File      string
}

var GlobalHist Hist

func (h *Hist) SetFilename(filename string) {
	h.Filename = filename
}

func (h *Hist) Add(Date time.Time, RuleID uint, Mask string, Completed bool, File string) {
	evt := Event{Date:Date, RuleID:RuleID, Mask:Mask, Completed:Completed, File:File}
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
	//file.WriteString(h.Events)
	return err
}