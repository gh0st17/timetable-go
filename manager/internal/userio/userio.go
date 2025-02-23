// Пакет userio предоставляет набор функции для
// обработки пользавательского ввода-вывода
//
// Основные функции:
//   - PrintLines: Печать списка групп с нумерацией
//   - GetUserSelection: Обработка пользовательского
//     ввода номера группы из списка
package userio

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gh0st17/timetable-go/params"
)

// Печать списка групп с нумерацией
func PrintGroupLines(lines []string, p *params.Params, printOnly bool) {
	if printOnly {
		for _, line := range lines {
			fmt.Println(line)
		}
	} else {
		fmt.Printf("Группы %d курса института №%d\n\n", p.Course, p.Dep)
		for i, line := range lines {
			fmt.Printf("%d. %s\n", i+1, line)
		}
	}
}

// Обработка пользовательского ввода номера группы из списка
func GetUserSelection(lines []string) uint64 {
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
		if err != nil || result < 1 || result > uint64(len(lines)) {
			fmt.Println("Неверный ввод. Попробуйте снова.")
			continue
		}

		fmt.Println()
		return result - 1
	}
}
