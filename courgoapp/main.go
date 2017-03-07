package main

import (
	"fmt"
	"github.com/Gtimur/courgo"
	"net"
	"bufio"
	"os"
	"log"
	"time"
	"strconv"
)

func main() {
	var web courgo.WebCtl

	fmt.Println(courgo.BannerString)
	if err := courgo.InitGlobal(); err != nil {
		fmt.Println("Ошибка запуска программы: ", err)
	}
	fmt.Println("Web control configured: "+"http://" + courgo.GlobalConfig.ManagerSrvAddr() + ":" + strconv.Itoa(int(courgo.GlobalConfig.ManagerSrvPort())))









	/* CHECK ARCDIR */
	courgo.GlobalArchi.SetDateNow()
	courgo.GlobalArchi.SetTmp("Y:\\TEMP\\TMP")
	courgo.GlobalArchi.SetSrc("Y:\\TEMP\\SRC")
	courgo.GlobalArchi.SetDst("Y:\\TEMP\\TEST")

	/*if err:=courgo.GlobalArchi.FullCopy(); err!=nil{
		log.Println(err)
	}*/

	fmt.Println(courgo.ArcDir(time.Now()))
	/****************/

	courgo.StartMonitor(courgo.GlobalMonCol, courgo.GlobalBook, courgo.GlobalConfig.SMTPCred())




	/*Запускаем сервер обслуживания "MENU"*/
	web.SetHost(net.ParseIP(courgo.GlobalConfig.ManagerSrvAddr()))
	web.SetPort(courgo.GlobalConfig.ManagerSrvPort())
	err := web.StartServe()
	if err != nil {
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
