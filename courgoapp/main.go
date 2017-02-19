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


	//Экземпляр структуры для хранения настроек программы
	cfg := courgo.Config{"test", courgo.PTKPSD{PathIN:"Y:\\test\\IN", PathOUT:"Y:\\test\\OUT"}, courgo.ManagerSrv{Addr:"127.0.0.1", Port:9090}, }

	cfg.JSONFile = "config.json"
	cfg.ReadJSON()
	//cfg.JSONFile = "config.bak"
	//cfg.NewConfig()
	//cfg.WriteJSON()
	//fmt.Println(cfg.PTKPath,cfg.ManagerSrv)

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

	if err:=courgo.InitGlobalBook("test.book"); err !=nil{
		log.Println(err.Error())
	}
	fmt.Println(courgo.GlobalBook.MaxID())
	courgo.GlobalBook.RemoveAccount(10)
	courgo.GlobalBook.RemoveAccount(21)
	res, _ := courgo.GlobalBook.StringJSON(0)
	fmt.Println("MAIN2:",res)

	var col courgo.MonitorCol
	col.SetJSONFile("col.json")
	courgo.ReadJSONFile(&col)
	courgo.StartMonitor(col,courgo.GlobalBook)


	/*Запускаем сервер обслуживания "MENU"*/
	web.SetHost(net.ParseIP("127.0.0.1"))
	web.SetPort(8000)
	err := web.StartServe()
	if err!=nil{
		log.Println(err)
		os.Exit(1)
	}
	/*MENU stop*/

	courgo.SendEmailMsg()


	//Ожидаем ввода новой строки
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Println()
	web.Close() //stop web-server gently
}
