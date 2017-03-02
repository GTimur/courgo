package main

import (
	"fmt"
	"github.com/Gtimur/courgo"
	"net"
	"bufio"
	"os"
	"log"
)


func main() {
	var book courgo.AddressBook = courgo.AddressBook{}
	var web courgo.WebCtl


	//Экземпляр структуры для хранения настроек программы
	cfg := courgo.Config{}
	cfg.SetJSONFile("config.json")
	cfg.SetManagerSrv("127.0.0.1",9090)
	cfg.SetSMTPSrv("smtp.yandex.ru",465,"to-timur@yandex.ru","blank","to-timur@yandex.ru","Информатор GO",true)
	//cfg.WriteJSON()
	cfg.ReadJSON()

	//Экземпляры мониторов
	var mon courgo.Monitor
	var mons []courgo.Monitor = []courgo.Monitor{}
	var acts []courgo.Action = []courgo.Action{}
	var act courgo.Action = courgo.Action{}
	mon.SetJSONFile("monitor.json")
	mon.SetFolder("Y:\\tmp")
	mon.SetMask([]string{"0z*.*","*.txt"})
	mon.SetMsgSubject("Посылка вида 0z от ГУ ЦБ")
	mon.SetMsgBody("Данное сообщение сформировано автоматически и не нуждается в ответе.\nИнформирую вас о получении посылки из ГУ ЦБ.")
	mon.SetSID([]uint64{10,20})
	act.SetId(10)
	act.SetDesc("Send email")
	acts = append(acts,act)
	mon.SetAction(acts)
	mons = append(mons,mon)




	//fmt.Println("MONS:",mons[0])
	//mons[0].MakeConfig()
	//courgo.WriteJSONFile(&mons[0])
	//courgo.ReadJSONFile(&mons[0])


	//book.Add(account)
	//book.Add(account1)
//	res, _ := book.StringJSON(0)
	// fmt.Println(res)
	book.JSONFile = "test.book"
	//book.Add(account1)
	//courgo.WriteJSONFile(&book)
	courgo.ReadJSONFile(&book)
	//res, _ := book.StringJSON(0)
	//fmt.Println("MAIN:",res)

/*
	if err:=courgo.InitGlobalBook("./static/data/test1.book"); err !=nil{
		log.Println(err.Error())
	}
*/
	if err:=courgo.InitGlobalBook("test.book"); err !=nil{
		log.Println(err.Error())
	}

	//fmt.Println(courgo.GlobalBook.MaxID())
	//courgo.GlobalBook.RemoveAccount(20)
	//courgo.GlobalBook.RemoveAccount(21)
	//res, _ := courgo.GlobalBook.StringJSON(0)
	//fmt.Println("MAIN2:",res)


	courgo.GlobalMonCol.SetJSONFile("col.json")
	//courgo.WriteJSONFile(&courgo.GlobalMonCol)
	courgo.ReadJSONFile(&courgo.GlobalMonCol)
	courgo.GlobalHist.SetFilename("history.dat")
	//courgo.StartMonitor(courgo.GlobalMonCol,courgo.GlobalBook, cfg.SMTPCred())

	/* CHECK ARCDIR
	courgo.GlobalArchi.SetDateNow()
	courgo.GlobalArchi.SetSrc("Y:\\TEMP\\SRC")
	courgo.GlobalArchi.SetDst("Y:\\TEMP\\TEST")

	if err:=courgo.GlobalArchi.FullCopy(); err!=nil{
		log.Println(err)
	}

	fmt.Println(courgo.ArcDir(time.Now()))
	/****************/



	/*Запускаем сервер обслуживания "MENU"*/
	web.SetHost(net.ParseIP(cfg.ManagerSrvAddr()))
	web.SetPort(cfg.ManagerSrvPort())
	err := web.StartServe()
	if err!=nil{
		log.Println(err)
		os.Exit(1)
	}
	/*MENU stop*/


	//Ожидаем ввода новой строки
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Println()
	web.Close() //stop web-server gently
}
