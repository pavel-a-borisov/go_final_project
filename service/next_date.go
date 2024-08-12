package service

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	FormatDate      = "20060102"
	ShortFormatDate = "02.01.2006"
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	//Если правило повторения не указано, возвращается ошибка.
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	// Пробуем преобразовать строку dateStr в тип time.Time согласно формату FormatDate.
	// Если формат неверный, возвращаем ошибку.
	date, err := time.Parse(FormatDate, dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	// Разбиваем строку repeat на части по пробелам.
	// Признак (d, y, w, m) сохраняем в переменную rule.
	parts := strings.Fields(repeat)
	rule := parts[0]

	var resultDate time.Time

	switch rule {
	case "d":
		//Для значение rule "d" (повторение по дням) проверяем, что указано ровно два элемента.
		if len(parts) != 2 {
			return "", errors.New("неверный формат повторения для 'd'")
		}

		// Преобразуем строку parts[1] в целое число days.
		days, err := strconv.Atoi(parts[1])

		// Если это число не корректно (меньше или равно 0 или превышает 400), возвращаем ошибку.
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неверное кол-во дней")
		}

		resultDate = date.AddDate(0, 0, days)

		for resultDate.Before(now) {
			resultDate = resultDate.AddDate(0, 0, days)
		}

	case "y":
		// Для правила "y" (повторение по годам) проверяем, что указан только один элемент.
		if len(parts) != 1 {
			return "", errors.New("неверный формат повторения для 'y'")
		}

		// Высчитываем новую дату, добавляя один год к дате date и повторяем это, пока результат не будет не раньше текущего момента.
		resultDate = date.AddDate(1, 0, 0)
		for resultDate.Before(now) {
			resultDate = resultDate.AddDate(1, 0, 0)
		}

	case "w":
		// Проверили, что для значения rule "w" (повторение по дням недели), указано ровно два элемента.
		if len(parts) != 2 {
			return "", errors.New("неверный формат повторения для 'w'")
		}

		// Разбиваем второй элемент на элементы по символу ",".
		days := strings.Split(parts[1], ",")
		weekdays := make([]int, len(days))
		for i, day := range days {
			weekday, err := strconv.Atoi(day)
			// Если это число не корректно (меньше 1 или превышает 7), возвращаем ошибку.
			if err != nil || weekday < 1 || weekday > 7 {
				return "", fmt.Errorf("недопустимый день недели: %v", day)
			}
			// Формируем массив дней недели
			weekdays[i] = weekday
		}

		if date.After(now) {
			resultDate = date
		} else {
			resultDate = now
		}

		// Создаем map, чтобы хранить смещение для каждого указанного дня недели (weekday).
		// Ключом будет день недели, а значением — количество дней до этого дня начиная от текущего now.
		closestDays := make(map[int]int)
		nowWeekday := int(resultDate.Weekday())
		for _, weekday := range weekdays {
			// Вычисляем смещение от текущего дня недели.
			// К этой разнице добавляем 7 и берём остаток от деления на 7 для получения положительного смещения в пределах одной недели.
			weekdayOffset := (weekday - nowWeekday + 7) % 7
			// Если смещение равно 0 (то есть текущий день недели совпадает с weekday), корректируем его до 7,
			// чтобы задача переносилась на следующее такое же число дня недели.
			if weekdayOffset == 0 {
				weekdayOffset = 7
			}
			closestDays[weekday] = weekdayOffset
		}

		// Устанавливаем начальное минимальное смещение на 7 дней (максимально возможное значение внутри недели),
		// чтобы гарантировать, что любое действительное смещение будет меньшим.
		minOffset := 7
		for _, offset := range closestDays {
			// Проходим по всем значениям карты closestDays (то есть по всем смещениям для указанных дней недели).
			// Если текущее смещение offset меньше чем minOffset, обновляем minOffset до значения offset.
			// Таким образом, завершая этот цикл, мы находим минимальное смещение среди всех указанных дней недели.
			if offset < minOffset {
				minOffset = offset
			}
		}

		// Устанавливаем resultDate на ближайший день недели
		resultDate = resultDate.AddDate(0, 0, minOffset)

	case "m":
		// Проверили, что для значения rule "m" (повторение по дням месяца ), указано ровно два элемента.
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("неверный формат повторения для 'm'")
		}

		// Записываем дни месяца в daysStr.
		daysStr := parts[1]

		// Если есть месяцы, мы также записываем их в monthsStr и преобразуем в массив months
		var monthsStr string
		if len(parts) >= 3 {
			monthsStr = parts[2]
		}

		days := strings.Split(daysStr, ",")
		var months []int
		if monthsStr != "" {
			for _, month := range strings.Split(monthsStr, ",") {
				m, err := strconv.Atoi(month)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("недопустимое число месяца: %v", m)
				}
				months = append(months, m)
			}
		}

		// проверяем значение дней месяца, что они попадают в заданный диапазон
		validDay := func(day int) bool {
			return (day >= 1 && day <= 31) || day == -1 || day == -2
		}

		// Преобразуем дни месяца в целые числа и проверяем их допустимость.
		intDays := []int{}
		for _, dayStr := range days {
			day, err := strconv.Atoi(dayStr)
			if err != nil || !validDay(day) {
				return "", fmt.Errorf("недопустимый день месяца: %v", day)
			}
			intDays = append(intDays, day)
		}

		// Сортируем массив дней в порядке возрастания
		sort.Slice(intDays, func(i, j int) bool {
			if (intDays[i] > 0 && intDays[j] > 0) || (intDays[i] < 0 && intDays[j] < 0) {
				// Если оба числа положительные или оба отрицательные,
				// то сортируем их по возрастанию.
				return intDays[i] < intDays[j]
			}
			// Положительные числа должны быть перед отрицательными.
			return intDays[i] > intDays[j]
		})

		// фиксируем значение now
		if date.After(now) {
			resultDate = date
		} else {
			resultDate = now
		}

		for {
			for _, day := range intDays {

				// Определяем текущий месяц.
				month := int(now.Month())
				if monthsStr != "" {
					// Проверяем, входит ли текущий месяц в список указанных месяцев.
					found := false
					for _, m := range months {
						if m == month {
							found = true
							break
						}
					}
					// Если текущий месяц не входит в список, продолжаем цикл.
					if !found {
						continue
					}
				}

				// Определяем следующую дату nextDate в зависимости от значения day.
				var nextDate time.Time
				if day > 0 {
					// Если day положительный, устанавливаем nextDate на этот день текущего месяца.
					nextDate = time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, now.Location())
					// Проверяем, что нужный день есть в месяце
					if nextDate.Day() != day {
						nextDate = time.Date(now.Year(), now.Month()+1, day, 0, 0, 0, 0, now.Location())
					}
				} else {
					// Если day отрицательный, вычисляем последний день текущего месяца и корректируем day.
					lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
					nextDate = time.Date(now.Year(), now.Month(), lastDay+day+1, 0, 0, 0, 0, now.Location())
				}

				// Если следующая дата позже текущей, возвращаем её.
				if nextDate.After(resultDate) {
					return nextDate.Format(FormatDate), nil
				}
			}
			// Переключаемся на следующий месяц.
			now = now.AddDate(0, 1, 0)
		}

	default:
		// Если правило не поддерживается (не "" , "d" или "y"), возвращаем ошибку.
		return "", errors.New("неподдерживаемый формат повторения задачи")
	}

	return resultDate.Format(FormatDate), nil
}
