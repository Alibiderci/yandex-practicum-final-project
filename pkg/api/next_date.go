package api

import (
	"strings"
	"strconv"
	"time"
	"fmt"
	"errors"
	"net/http"
)

// layout определяет стандартный формат даты YYYYMMDD, используемый во всем приложении.
const layout = "20060102"

// afterNow сравнивает две даты, игнорируя время суток.
// Возвращает true, если `date` строго позже (на следующий день или далее), чем `now`.
func afterNow(date, now time.Time) bool {
	// Truncate обрезает время до начала дня (00:00:00).
	dateOnly := date.Truncate(24 * time.Hour)
	nowOnly := now.Truncate(24 * time.Hour)

	return dateOnly.After(nowOnly)
}

// NextDate вычисляет следующую дату выполнения задачи на основе правила повторения.
// now - текущая дата, от которой идет поиск.
// dstart - дата, с которой начался отсчет повторений.
// repeat - строка с правилом повторения.
// Возвращаемая дата всегда будет строго больше `now`.
func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("параметр repeat не должен быть пустым")
	}

	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка во время парсинга даты старта задания")
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	// Правило "d": Повторение через N дней 
	case "d":
		// Проверяем, что правило указано полностью (например, "d 7", а не просто "d").
		if len(parts) < 2 {
			return "", errors.New("Не указан интервал в днях")
		}

		// Конвертируем второй элемент ("7") в число.
		days, err := strconv.Atoi(parts[1]) 
		if err != nil {
			return "", fmt.Errorf("Ошибка при конвертировании: %w", err)
		}

		// Проверяем, что число дней находится в допустимом диапазоне (1-400).
		if days <= 0 || days > 400 {
			return "", errors.New("Превышен максимально допустимый интервал")
		}
		
		// В цикле добавляем указанное количество дней к последней дате.
		// Продолжаем до тех пор, пока не найдем первую дату, которая будет строго после 'now'.
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break // Нашли подходящую дату, выходим из цикла.
			}
		}

		return date.Format(layout), nil

	// Правило "y": Ежегодное повторение 
	case "y":
		// Логика аналогична правилу "d", но мы всегда добавляем ровно один год.
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

		return date.Format(layout), nil

	// Правило "w": Повторение по дням недели 
	case "w":
		// Проверяем, что правило указано полностью (например, "w 2,4,6", а не просто "w").
		if len(parts) < 2 {
			return "", errors.New("Не указаны дни недели")
		}

		// Используем массив bool как "множество" (set) для быстрой проверки.
		// Индекс 0 не используется, т.к. дни недели нумеруются с 1 до 7.
		var validWeekdays [8]bool
		daysOfWeekStr := strings.Split(parts[1], ",")
		for _, s := range daysOfWeekStr {
			day, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("Неверный формат недели: %w", err)
			}

			// Валидация: день недели должен быть от 1 (Пн) до 7 (Вс).
			if day < 1 || day > 7 {
				return "", errors.New("Недопустимое значение для дня недели")
			}
			validWeekdays[day] = true // Отмечаем день как допустимый.
		}

		// Итерируемся по одному дню вперед от последней даты
		// и проверяем, подходит ли каждый новый день под наши правила.
		for {
			date = date.AddDate(0, 0, 1)

			// Приводим день недели из формата Go (Воскресенье=0) к нашему формату (Понедельник=1..Воскресенье=7).
			currentWeekday := int(date.Weekday())
			if currentWeekday == 0 {  // Если это Воскресенье (0 в Go)...
				currentWeekday = 7  // ...то для нас это 7.
			}

			// Проверяем, является ли текущий день недели одним из допустимых.
			if validWeekdays[currentWeekday] {
				// Если день недели подходит, проверяем, находится ли эта дата в будущем.
				if afterNow(date, now) {
					// Если да, то мы нашли первую подходящую будущую дату. Возвращаем ее.
					return date.Format(layout), nil	
				}
			}
		}

	// Правило "m": Повторение по дням и месяцам
	case "m":
		// Проверяем, что в правиле указаны хотя бы дни месяца (например, "m 10,20", а не просто "m").
		if len(parts) < 2 {
			return "", errors.New("Не указаны дни месяца")
		}

		// Готовим структуры для валидации дней месяца
		var validDays [32]bool  // Для дней 1-31.
		var isLastDay, isPenultimateDay bool // Флаги для особых правил -1 и -2.

		daysStr := strings.Split(parts[1], ",")
		for _, s := range daysStr {
			day, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("Неверный формат дня месяца: %w", err)
			}

			switch {
			case day == -1: // Правило "последний день месяца".
				isLastDay = true
			case day == -2: // Правило "предпоследний день месяца".
				isPenultimateDay = true
			case day >= 1 && day <= 31:
				validDays[day] = true
			default:
				return "", errors.New("Недопустимый день месяца")
			}
		}
	    // Готовим структуру для валидации месяцев
		var validMonths [13]bool

		// Флаг, который показывает, были ли указаны конкретные месяцы.
		// Если нет, то правило применяется ко всем месяцам.
		allMonthsSpecified := true

		// Если указаны конкретные месяца
		if len(parts) > 2 {
			allMonthsSpecified = false // Пользователь указал месяцы, флаг в false.
			monthsStr := strings.Split(parts[2], ",")
			for _, s := range monthsStr {
				month, err := strconv.Atoi(s)
				if err != nil {
					return "", fmt.Errorf("Неверный формат месяца: %w", err)
				}
				if month < 1 || month > 12 {
					return "", errors.New("Недопустимый месяц")
				}
				validMonths[month] = true
			}
		} else {
			// Если месяцы не указаны, считаем все 12 месяцев допустимыми.
			// Хоть и есть флаг allMonthsSpecified, не будет лишним иметь такой защитный механизм
			for i := 1; i <= 12; i++ {
				validMonths[i] = true
			}
		}

		// Итерируемся по одному дню, как и в правиле "w".
		for {
			date = date.AddDate(0, 0, 1)

			// Если указаны конкретные месяцы и текущий месяц не в их числе,
			// то сразу переходим к следующему дню, не делая лишних проверок
			if !allMonthsSpecified && !validMonths[int(date.Month())] {
				continue
			}

			// Проверяем, соответствует ли текущий день правилам.
			dayOfMonth := date.Day()
			// День является последним, если следующий день - 1-е число.
			isLast := date.AddDate(0, 0, 1).Day() == 1
			// День является предпоследним, если послезавтра - 1-е число.
			isPenultimate := date.AddDate(0, 0, 2).Day() == 1

			// День подходит, если:
			// - его номер есть в списке (validDays[dayOfMonth]), ИЛИ
			// - правило "последний день" активно и сегодня последний день (isLastDay && isLast), ИЛИ
			// - правило "предпоследний день" активно и сегодня предпоследний день (isPenultimateDay && isPenultimate).
			isDayMatch := validDays[dayOfMonth] || (isLastDay && isLast) || (isPenultimateDay && isPenultimate)

			// Если день подошел, проверяем, что он в будущем.
			if isDayMatch {
				if afterNow(date, now) {
					return date.Format(layout), nil
				}
			}
		}
	
	default: 
		return "", errors.New("Неподдерживаемый формат")
	}
}

// nextDateHandler обрабатывает GET-запросы к эндпоинту /api/nextdate.
// Его основная задача — быть "оберткой" для сложной функции NextDate,
// предоставляя к ней доступ через HTTP. Он получает все необходимые параметры
// из URL-запроса (query parameters), вызывает NextDate для вычисления
// и возвращает результат в виде простого текста.
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "неподдерживаемый метод"})
		return
	}

	var now time.Time
	var err error

	// Получаем GET-параметр 'now' из URL. r.FormValue удобно работает и для GET, и для POST.
	nowStr := r.FormValue("now")
	if nowStr == "" {
		// Если параметр 'now' не был предоставлен, в качестве точки отсчета
		// используется текущее время на сервере. Это поведение по умолчанию.
		now = time.Now()
	} else {
		// Если параметр 'now' есть, пытаемся его распарсить из строки.
		// Константа `layout` ("20060102") гарантирует правильный формат.
		now, err = time.Parse(layout, nowStr)	
		if err != nil {
			// Если парсинг не удался, значит клиент прислал дату в неверном формате.
			// Возвращаем ошибку 400 Bad Request.
			http.Error(w, "неправильный формат текущего времени", http.StatusBadRequest)
			return
		}
	}

	// Получаем обязательные параметры 'date' (дата старта задачи) и 'repeat' (правило повторения).
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Вызываем главную функцию nextDate для выполнения всех сложных вычислений.
	nextDate, err := nextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}
