// Инициализация глобальных переменных и структур
package courgo

import (
	"fmt"
	"time"
)

const (
	Version = "0.3.1"
	BannerString = "Courier Go notification utility. " + Version + " (C) 2017 UMK BANK (GTG)" + "\n" +
		"USE: courgo.exe [start]"

	configFile = "config.json"
)

// Инициализация глобальных переменных с обработкой ошибок
func InitGlobal() error {

	/* Экземпляр структуры для хранения настроек программы */
	GlobalConfig.SetJSONFile(configFile)
	if err := GlobalConfig.ReadJSON(); err != nil {
		fmt.Println("Configuration file " + configFile + " not found. Creating config from template.\nPlease visit " + "http://127.0.0.1:8000/config")
	}
	/* Если конфига нет, создаем базовый */
	GlobalConfig.SetJSONFile("config.json")
	if len(GlobalConfig.managerSrv.Addr) == 0 {
		GlobalConfig.SetManagerSrv("127.0.0.1", 8000)
		GlobalConfig.SetSMTPSrv("blank", 25, "enter@your.mail", "blank", "enter@noti.mail", "Информатор GO", false)
		if err := GlobalConfig.MakeConfig(); err != nil {
			return err
		}
		if err := WriteJSONFile(&GlobalConfig); err != nil {
			return err
		}
	}

	//GlobalConfig.smtpSrv.Password="test"
	//WriteJSONFile(&GlobalConfig)
	//fmt.Println(GlobalConfig.smtpSrv.Password)

	/* Мониторы */
	GlobalMonCol.SetJSONFile("moncol.json")
	if err := ReadJSONFile(&GlobalMonCol); err != nil {
		return err
	}

	/* Каталог подписчиков (адресная книга) */
	if err := InitGlobalBook("adrbook.json"); err != nil {
		return err
	}

	/* История работы монитора */
	GlobalHist.SetFilename("history.dat")
	GlobalHist.SetJSONFile("history.json")
	if err := ReadJSONFile(&GlobalHist); err != nil {
		fmt.Println("History: history.json not found and will be created.")
		if err := GlobalHist.MakeJSONFile(); err != nil {
			return err
		}
	} else {
		fmt.Println("History: cleaning history file until today.")
		t := time.Now()
		t.Add(-24*time.Hour)
		GlobalHist.CleanUntilDay(t)
		if err := GlobalHist.RewriteJSONFile(); err != nil {
			return err
		}
	}

	return nil
}
