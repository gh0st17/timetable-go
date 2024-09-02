package manager

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"timetable/errtype"
)

// Чтение файла и возврат массива строк
func readLines(filePath string) ([]string, error) {
	var (
		file *os.File
		err  error
	)

	if file, err = os.Open(filePath); err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, errtype.RuntimeError(fmt.Errorf("файл групп %s пустой", filePath))
	}

	return lines, nil
}

// Функция для записи массива строк в файл
func writeLines(filePath string, lines *[]string) error {
	var (
		file *os.File
		err  error
	)

	// Открываем файл для записи (перезаписываем файл)
	if file, err = os.Create(filePath); err != nil {
		return err
	}
	defer file.Close()

	// Записываем строки в файл
	for _, line := range *lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Функция для записи строки в файл
func writeString(filePath string, data *string) error {
	var (
		file *os.File
		err  error
	)

	if file, err = os.Create(filePath); err != nil {
		return errtype.RuntimeError(fmt.Errorf("ошибка создания файла %s: %s", filePath, err))
	}
	defer file.Close()

	if _, err = file.WriteString(*data); err != nil {
		return errtype.RuntimeError(fmt.Errorf("ошибка записи в файл %s: %s", filePath, err))
	}

	return nil
}

// Печать строк с нумерацией
func printLines(lines *[]string, p *Params, printOnly bool) {
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
	if info, err := os.Stat(dirPath); err != nil {
		return false
	} else {
		return info.IsDir()
	}
}

// Функция для создания директории
func createDir(dirPath string) error {
	// Права доступа: rwxr-xr-x
	return os.Mkdir(dirPath, 0755)
}

// Функция для удаления всех файлов в папке
func removeAllFilesInDir(dirPath string) error {
	var (
		entries []fs.DirEntry
		err     error
	)

	if entries, err = os.ReadDir(dirPath); err != nil {
		return err
	}

	for _, entry := range entries {
		// Проверяем, что это файл, а не директория
		if !entry.IsDir() {
			filePath := filepath.Join(dirPath, entry.Name())
			if err = os.Remove(filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Функция для получения абсолютного пути запускаемой программы
func getWd() (string, error) {
	if executable, err := os.Getwd(); err != nil {
		return "", err
	} else {
		return filepath.Abs(executable)
	}
}
