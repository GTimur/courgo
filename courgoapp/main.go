package main

import (
	"fmt"
	"github.com/Gtimur/courgo"
	"net"
	"os"
	"log"
	"strconv"
	"time"
	"github.com/GeertJohan/go.rice"
)

func main() {
	var web courgo.WebCtl

	/* Приветствие */
	fmt.Println(courgo.BannerString)
	if err := courgo.InitGlobal(); err != nil {
		fmt.Println("Ошибка запуска программы: ", err)
	}
	fmt.Println("Web control configured: " + "http://" + courgo.GlobalConfig.ManagerSrvAddr() + ":" + strconv.Itoa(int(courgo.GlobalConfig.ManagerSrvPort())))

	_, err := rice.FindBox("static")
	if err != nil {
		log.Fatalln("go.rice embedding error: ", err)
	}

	/* Запускаем сервер обслуживания WebCtl */
	web.SetHost(net.ParseIP(courgo.GlobalConfig.ManagerSrvAddr()))
	web.SetPort(courgo.GlobalConfig.ManagerSrvPort())
	err = web.StartServe()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	/* Запускаем обработчик правил монитора */
	courgo.MonSvcState = false
	ticker := time.NewTicker(time.Second * 1)
	i := 0;
	for _ = range ticker.C {
		// Запускаем обработчик каждую минуту
		if i != courgo.Interval {
			i++
			courgo.TimeRemain = courgo.Interval - i
			//fmt.Println("tick", i)
			if courgo.WaitExit {
				time.Sleep(1 * time.Second)
				break
			}
			continue
		}
		i = 0
		if err := courgo.StartMonitor(courgo.GlobalMonCol, courgo.GlobalBook, courgo.GlobalConfig.SMTPCred(), courgo.MonSvcState); err != nil {
			log.Fatal("Ошибка: Не удалось запустить диспетчер правил.", err)
		}
	}

	ticker.Stop()

	/* WebCtl stop */
	web.Close() //stop web-server gently
	if !courgo.WaitExit {
		os.Exit(0)
	}
	log.Println("Работа программы была завершена по требованию пользователя.")
}
