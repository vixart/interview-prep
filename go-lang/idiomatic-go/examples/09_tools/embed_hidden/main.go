// Три способа встроить каталог и что происходит со скрытыми файлами
// (имена на . и _): parent_dir — без скрытых, parent_dir/* — скрытые только
// в самом каталоге, all:parent_dir — скрытые на всех уровнях.
package main

import (
	"embed"
	"fmt"
)

//go:embed parent_dir
var noHidden embed.FS

// скрытые файлы НЕ попадут

//go:embed parent_dir/*
var parentHiddenOnly embed.FS

// звездочка добавит скрытые только в самом parent_dir

//go:embed all:parent_dir
var allHidden embed.FS

// all: — скрытые на всех уровнях вложенности

func main() {
	checkForHidden("noHidden", noHidden)
	checkForHidden("parentHiddenOnly", parentHiddenOnly)
	checkForHidden("allHidden", allHidden)
}

func checkForHidden(name string, dir embed.FS) {
	fmt.Println(name)
	allFileNames := []string{"parent_dir/.hidden", "parent_dir/child_dir/.hidden"}
	for _, v := range allFileNames {
		_, err := dir.Open(v)
		if err == nil {
			fmt.Println(v, "found")
		}
	}
	fmt.Println()
}
