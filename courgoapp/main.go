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
	"strings"
	"os/signal"
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
		log.Println("HTTP сервер: Ошибка. ", err)
		os.Exit(1)
	}

	/* Перехват CTRL+C */
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("\nReceived %v, shutdown procedure initiated.\n\n", sig)
			courgo.WaitExit = true
		}
	}()

	/* Запускаем обработчик правил монитора */
	if len(os.Args) >= 2 {
		courgo.MonSvcState = strings.Contains(os.Args[1], "start")
	}
	ticker := time.NewTicker(time.Second * 1)
	i := 0; // таймер интервала выполнения правил монитора
	h := 0; // таймер интервала сохранения аварийной истории (JSON)

	for range ticker.C {
		// Запускаем обработчик каждую минуту
		h++
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
		/* Сохраняем системную(JSON) историю раз в час */
		if h <= 3600 {
			continue
		}
		h = 0
		if err := courgo.GlobalHist.RewriteDumpJSON(); err != nil {
			fmt.Println("Ошибка сохранения системной (JSON) истории:", err)
		}
	}
	// Сохраняем системную(JSON) историю в случае штатного завершения программы
	if err := courgo.GlobalHist.RewriteDumpJSON(); err != nil {
		fmt.Println("Ошибка сохранения системной (JSON) истории:", err)
	}
	ticker.Stop()

	/* WebCtl stop */
	// stop web-server gently
	web.Close()
	if !courgo.WaitExit {
		os.Exit(0)
	}
	fmt.Println("Работа программы была завершена по требованию пользователя. ", time.Now())
}
