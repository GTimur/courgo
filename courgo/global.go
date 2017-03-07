// Инициализация глобальных переменных
package courgo

const (
	Version = "0.1"
	BannerString = "Courier Go notification utility. " + Version + " (C) 2017 UMK BANK (GTG)"
)

// Инициализация глобальных переменных с обработкой ошибок
func InitGlobal() error {

	/* Экземпляр структуры для хранения настроек программы */
	GlobalConfig.SetJSONFile("config.json")
	//GlobalConfig.SetManagerSrv("127.0.0.1",9090)
	//GlobalConfig.SetSMTPSrv("smtp.yandex.ru",465,"to-timur@yandex.ru","blank","to-timur@yandex.ru","Информатор GO",true)
	if err := GlobalConfig.ReadJSON(); err != nil {
		return err
	}

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
