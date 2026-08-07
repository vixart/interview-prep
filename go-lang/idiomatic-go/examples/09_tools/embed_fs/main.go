// //go:embed каталога дает embed.FS — виртуальную файловую систему внутри бинарника,
// по ней можно ходить через io/fs (WalkDir) и читать файлы (ReadFile).
// Запуск: go run .  (список тем) или go run . info.txt
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// Встраиваем целый каталог. Для скрытых файлов (. и _) пишут all:help.
//
//go:embed help
var helpInfo embed.FS // файловая система внутри бинарника, реализует io/fs

// embed.FS — файловая система внутри бинарника

func main() {
	if len(os.Args) == 1 {
		printHelpFiles()
		os.Exit(0)
	}
	data, err := helpInfo.ReadFile("help/" + os.Args[1])
	// путь всегда с именем встроенного каталога и через прямой слеш
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func printHelpFiles() {
	fmt.Println("contents:")
	fs.WalkDir(helpInfo, "help",
		// embed.FS реализует io/fs — работают все обычные обходы
		func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				_, fileName, _ := strings.Cut(path, "/")
				fmt.Println(fileName)
			}
			return nil
		})
}
