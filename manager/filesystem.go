package manager

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"timetable/errtype"
	"timetable/params"
)

// func load_from_file(filepath string) (doc *html.Node) {
// 	// Чтение файла
// 	file, err := os.ReadFile(filepath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	doc, err = html.Parse(strings.NewReader(string(file)))
// 	if err != nil {
// 		panic(err)
// 	}

// 	return doc
// }

// Чтение файла и возврат массива строк
func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, errtype.RuntimeError(fmt.Errorf("файл групп %s пустой", filePath))
	}

	return lines, nil
}

// Функция для записи массива строк в файл
func writeFile(filePath string, lines *[]string) error {
	// Открываем файл для записи (перезаписываем файл)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем строки в файл
	for _, line := range *lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Печать строк с нумерацией
func printLines(lines *[]string, p *params.Params, printOnly bool) {
	if printOnly {
		for _, line := range *lines {
			fmt.Println(line)
		}
	} else {
		fmt.Printf("Группы %d курса института №%d\n\n", p.Course, p.Dep)
		for i, line := range *lines {
			fmt.Printf("%d. %s\n", i+1, line)
		}
	}
}

// Обработка пользовательского ввода
func getUserSelection(lines *[]string) uint64 {
	var (
		err    error
		input  string
		result uint64
	)

	for {
		fmt.Printf("\nВыберите номер группы в списке: ")
		input, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		input = strings.TrimSpace(input)
		result, err = strconv.ParseUint(input, 10, 64)
		if err != nil || result < 1 || result > uint64(len(*lines)) {
			fmt.Println("Неверный ввод. Попробуйте снова.")
			continue
		}

		fmt.Println()
		return result - 1
	}
}

// Функция для проверки существования файла
func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}

	return true
}

// Функция для проверки существования директории
func dirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return info.IsDir()
}

// Функция для создания директории
func createDir(dirPath string) error {
	err := os.Mkdir(dirPath, 0755) // Права доступа: rwxr-xr-x
	if err != nil {
		return err
	}
	return nil
}

// Функция для получения абсолютного пути запускаемой программы
func getWd() (string, error) {
	executable, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Abs(executable)
}
