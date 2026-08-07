// Функция читает файл — значит тесту нужны данные на диске (см. каталог testdata).
package text

import (
	"os"
	"unicode/utf8"
)

func CountCharacters(fileName string) (int, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	return utf8.RuneCount(data), nil
	// считаем руны, а не байты — файл в testdata содержит не только ASCII
}
