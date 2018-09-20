/*
	Реализует архивацию обработанных объектов

	Реализует функции работы с архивными файлами посредствам unrar
*/
package courgo

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Глобальная переменная для обслуживания архиватора
var GlobalArchi Archi

type Archi struct {
	date time.Time
	/* Пути */
	src string
	dst string
	// Каталог для хранения временных файлов
	tmp string
}

func (a *Archi) SetDate(date time.Time) {
	a.date = date
}

func (a *Archi) SetDateNow() {
	a.date = time.Now()
}

func (a *Archi) SetSrc(src string) {
	a.src = src
}

func (a *Archi) SetDst(dst string) {
	a.dst = dst
}

func (a *Archi) SetTmp(tmp string) {
	a.tmp = tmp
}

// Возвращает директорию вида ГГГГ\ММ\ДД с учетом даты.
func ArcDir(date time.Time) string {
	res := strconv.Itoa(date.Year())                     //ГГГГ
	res += "\\" + fmt.Sprintf("%02d", int(date.Month())) //ММ
	res += "\\" + fmt.Sprintf("%02d", date.Day())        //ДД
	return res
}

// Возвращает директорию вида ГГГГ\ММ\ДД с учетом даты.
func (a *Archi) ArcDirDstNow() string {
	date := time.Now()
	res := strconv.Itoa(date.Year())                     //ГГГГ
	res += "\\" + fmt.Sprintf("%02d", int(date.Month())) //ММ
	res += "\\" + fmt.Sprintf("%02d", date.Day())        //ДД
	return a.dst + "\\" + res
}

// Создает директорию для архива
// вида dst\yyyy\mm\dd\
func (a *Archi) MakeDir() error {
	dir := a.dst + "\\" + ArcDir(a.date)
	// Если директории не существует
	if _, err := os.Stat(dir); os.IsExist(err) {
		return err
	}
	// Создаем
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Копирует указанный файл в архивную директорию
// относительно Archi.dst
// file = полный путь включая имя файла
func (a *Archi) MakeCopy(file string) error {
	// Создадим папку для архива (если её нет)
	if err := a.MakeDir(); err != nil {
		return err
	}

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Printf("Archi: Невозможно открыть файл: %v\n", err)
		return err
	}

	dstfilename := a.dst + "\\" + ArcDir(time.Now()) + "\\" + filepath.Base(file)

	if _, err = os.Stat(dstfilename); os.IsExist(err) {
		// Файл существует и не будет перезаписан
		return err
	}

	df, err := os.Create(dstfilename)
	defer df.Close()

	if _, err = io.Copy(df, f); err != nil {
		return err
	}
	if err := df.Sync(); err != nil {
		return err
	}

	return err
}

// Полня копия содержимого папки SRC в DST
func (a *Archi) FullCopy() error {
	// Создадим папку для архива (если её нет)
	if err := a.MakeDir(); err != nil {
		return err
	}
	if err := MakeCopyAll(a.ArcDirDstNow(), a.src); err != nil {
		return err
	}
	return nil
}

// Копирует указанный файл в архивную директорию
// src,dst = полный путь включая имя файла
func MakeCopy(dst string, src string) error {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil {
		log.Printf("Archi: Невозможно открыть файл: %v\n", err)
		return err
	}

	if _, err = os.Stat(dst); os.IsExist(err) {
		// Файл существует и не будет перезаписан
		return err
	}

	df, err := os.Create(dst)
	defer df.Close()

	if _, err = io.Copy(df, f); err != nil {
		return err
	}
	if err := df.Sync(); err != nil {
		return err
	}

	return err
}

// Копирует все файлы в архивную директорию dst
func MakeCopyAll(dst string, src string) error {
	mask := []string{"*.*", "*"}
	// Готовим список всех файлов имеющихся в директории src
	files := findFiles(src, mask)
	// Копируем каждый файл в папку назначения dst
	for f := range files {
		if err := MakeCopy(dst+"\\"+filepath.Base(f), f); err != nil {
			return err
		}
	}
	return nil
}

// Разархивирует архивные файлы при помощи unrar.exe
// destpath = путь к папке для разархивации
func unArc(file string, dstpath string) error {
	if !ifPathExist(file) {
		return errors.New("UnArc: Файл \"" + file + "\" не существует.")
	}
	if !ifPathExist(dstpath) {
		return errors.New("UnArc: путь \"" + dstpath + "\" не существует.")
	}

	archiver := "unrar.exe"
	commandString := fmt.Sprintf(archiver+` e -ep1 %s %s`, file, dstpath)
	//commandString := fmt.Sprintf(archiver + ` e %s %s`, file, dstpath)

	if strings.Contains(strings.ToUpper(filepath.Ext(file)), ".ARJ") {
		archiver = "arj.exe"
		commandString = fmt.Sprintf(archiver+` e -p1 %s %s`, file, dstpath)
		//commandString = fmt.Sprintf(archiver + ` e %s %s`, file, dstpath)
	}

	if strings.Contains(strings.ToUpper(filepath.Ext(file)), ".ZIP") {
		archiver = "unzip.exe"
		commandString = fmt.Sprintf(archiver+` -j %s -d %s`, file, dstpath)
		//commandString = fmt.Sprintf(archiver + ` -x %s -d %s`, file, dstpath)
	}

	commandSlice := strings.Fields(commandString)
	fmt.Println(commandString)
	c := exec.Command(commandSlice[0], commandSlice[1:]...)
	if err := c.Run(); err != nil {
		fmt.Printf("Запуск %s для файла %s: %v\n", archiver, file, err)
		if strings.Contains(strings.ToUpper(err.Error()), "EXIT STATUS 2") {
			return err
		}
	}
	return nil
}

// Выполняет разархивацию файла во временный каталог
// и возвращает список файлов из временного каталога
// mask = маска по которой был найден файл (для monsvc)
func prepUnArc(file, dir string) (unarc map[string]string, err error) {
	if len(file) == 0 {
		return unarc, errors.New("PrepUnArc: Файлы для обработки не найдены.")
	}

	// Распакуем файл во временный каталог
	if err := unArc(file, dir); err != nil {
		return unarc, err
	}

	// Соберем все найденные файлы
	unarc = findFiles(dir, []string{"*.*", "*"})

	if len(unarc) == 0 {
		return unarc, errors.New("PrepUnArc: Архив пустой.")
	}

	return unarc, nil
}
