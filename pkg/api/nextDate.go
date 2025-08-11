package api

import (
	"strings"
	"strconv"
	"time"
	"fmt"
	"errors"
	"net/http"
)

const layout = "20060102"

func afterNow(date, now time.Time) bool {
	dateOnly := date.Truncate(24 * time.Hour)
	nowOnly := now.Truncate(24 * time.Hour)

	return dateOnly.After(nowOnly)
}

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
	case "d":
		if len(parts) < 2 {
			return "", errors.New("Не указан интервал в днях")
		}
		days, err := strconv.Atoi(parts[1]) 
		if err != nil {
			return "", fmt.Errorf("Ошибка при конвертировании: %w", err)
		}
		if days <= 0 || days > 400 {
			return "", errors.New("Превышен максимально допустимый интервал")
		}
		
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

		return date.Format(layout), nil

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

		return date.Format(layout), nil

	case "w":
		if len(parts) < 2 {
			return "", errors.New("Не указаны дни недели")
		}

		var validWeekdays [8]bool
		daysOfWeekStr := strings.Split(parts[1], ",")
		for _, s := range daysOfWeekStr {
			day, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("Неверный формат недели: %w", err)
			}
			if day < 1 || day > 7 {
				return "", errors.New("Недопустимое значение для дня недели")
			}
			validWeekdays[day] = true
		}

		for {
			date = date.AddDate(0, 0, 1)

			currentWeekday := int(date.Weekday())
			if currentWeekday == 0 {
				currentWeekday = 7
			}

			if validWeekdays[currentWeekday] {
				if afterNow(date, now) {
					return date.Format(layout), nil	
				}
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("Не указаны дни месяца")
		}

		var validDays [32]bool
		var isLastDay, isPenultimateDay bool

		daysStr := strings.Split(parts[1], ",")
		for _, s := range daysStr {
			day, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("Неверный формат дня месяца: %w", err)
			}

			switch {
			case day == -1:
				isLastDay = true
			case day == -2:
				isPenultimateDay = true
			case day >= 1 && day <= 31:
				validDays[day] = true
			default:
				return "", errors.New("Недопустимый день месяца")
			}
		}

		var validMonths [13]bool
		allMonthsSpecified := true

		if len(parts) > 2 {
			allMonthsSpecified = false
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
			for i := 1; i <= 12; i++ {
				validMonths[i] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			if !allMonthsSpecified && !validMonths[int(date.Month())] {
				continue
			}

			dayOfMonth := date.Day()
			isLast := date.AddDate(0, 0, 1).Day() == 1
			isPenultimate := date.AddDate(0, 0, 2).Day() == 1

			isDayMatch := validDays[dayOfMonth] || (isLastDay && isLast) || (isPenultimateDay && isPenultimate)

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



func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	var err error

	nowStr := r.FormValue("now")
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(layout, nowStr)	
		if err != nil {
			http.Error(w, "неправильный формат текущего времени", http.StatusBadRequest)
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := nextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}
