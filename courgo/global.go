// Инициализация глобальных переменных
package courgo

const (
	Version = "0.2"
	BannerString = "Courier Go notification utility. " + Version + " (C) 2017 UMK BANK (GTG)"
)

// Инициализация глобальных переменных с обработкой ошибок
func InitGlobal() error {

	/* Экземпляр структуры для хранения настроек программы */
	GlobalConfig.SetJSONFile("config.json")
	if err := GlobalConfig.ReadJSON(); err != nil {
		return err
	}
	/* Если конфига нет, создаем базовый */
	if len(GlobalConfig.managerSrv.Addr) == 0 {
		GlobalConfig.SetManagerSrv("127.0.0.1", 8000)
		GlobalConfig.SetSMTPSrv("blank", 25, "enter@your.mail", "blank", "enter@noti.mail", "Информатор GO", false)
		WriteJSONFile(&GlobalConfig)
	}




	//GlobalConfig.smtpSrv.Password="test"
	//WriteJSONFile(&GlobalConfig)
	//fmt.Println(GlobalConfig.smtpSrv.Password)

	/* Мониторы */
	GlobalMonCol.SetJSONFile("col.json")
	if err := ReadJSONFile(&GlobalMonCol); err != nil {
		return err
	}

	/* Каталог подписчиков (адресная книга) */
	if err := InitGlobalBook("test.book"); err != nil {
		return err
	}
	/*
	if err:=InitGlobalBook("./static/data/test1.book"); err !=nil{
	return err
	}*/

	/* История работы монитора */
	GlobalHist.SetFilename("history.dat")

	return nil
}
