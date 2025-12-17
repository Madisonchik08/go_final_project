package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate calculates the next date according to repeat rules.
// now     - current reference time
// dstart  - base date in format 20060102
// repeat  - repetition rule (d, y, w, m formats)
// NextDate вычисляет следующую дату согласно правилам повторения.
// now     - текущее опорное время
// dstart  - базовая дата в формате 20060102
// repeat  - правило повторения (форматы d, y, w, m)
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	date = normalize(date)
	ref := normalize(now)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid day repeat format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		next := advanceByDays(date, ref, interval)
		return next.Format("20060102"), nil
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid yearly repeat format")
		}
		next := advanceByYear(date, ref)
		return next.Format("20060102"), nil
	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid weekly repeat format")
		}
		days, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		next := findNextMatching(date, ref, func(t time.Time) bool {
			wd := int(t.Weekday())
			if wd == 0 {
				wd = 7
			}
			return days[wd]
		})
		return next.Format("20060102"), nil
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid monthly repeat format")
		}
		days, err := parseMonthDays(parts[1])
		if err != nil {
			return "", err
		}
		months := [13]bool{}
		if len(parts) == 3 {
			months, err = parseMonths(parts[2])
			if err != nil {
				return "", err
			}
		} else {
			for i := 1; i <= 12; i++ {
				months[i] = true
			}
		}

		next := findNextMatching(date, ref, func(t time.Time) bool {
			if !months[int(t.Month())] {
				return false
			}
			return matchMonthDay(t, days)
		})
		return next.Format("20060102"), nil
	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}

// normalize removes time component, keeping only date.
// normalize удаляет компонент времени, оставляя только дату.
func normalize(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// afterNow checks if date is after now.
// afterNow проверяет, что дата позже текущей.
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// advanceByDays advances date by specified interval in days.
// advanceByDays продвигает дату на указанный интервал в днях.
func advanceByDays(date, now time.Time, interval int) time.Time {
	d := date.AddDate(0, 0, interval)
	for !afterNow(d, now) {
		d = d.AddDate(0, 0, interval)
	}
	return d
}

// advanceByYear advances date by one year.
// advanceByYear продвигает дату на один год.
func advanceByYear(date, now time.Time) time.Time {
	d := addOneYear(date)
	for !afterNow(d, now) {
		d = addOneYear(d)
	}
	return d
}

// addOneYear adds one year to the date, handling leap year edge cases.
// addOneYear добавляет один год к дате, обрабатывая крайние случаи високосных годов.
func addOneYear(t time.Time) time.Time {
	d := t.AddDate(1, 0, 0)
	// Adjust Feb 29 to Mar 1 on non-leap years.
	// Корректируем 29 февраля на 1 марта в невисокосные годы.
	if t.Month() == time.February && t.Day() == 29 && d.Month() == time.February && d.Day() == 28 {
		return d.AddDate(0, 0, 1)
	}
	return d
}

// findNextMatching finds next date that matches the condition.
// findNextMatching находит следующую дату, соответствующую условию.
func findNextMatching(date, now time.Time, match func(time.Time) bool) time.Time {
	d := date.AddDate(0, 0, 1)
	for {
		if afterNow(d, now) && match(d) {
			return d
		}
		d = d.AddDate(0, 0, 1)
	}
}

// parseWeekdays parses weekday numbers from string (1-7, comma-separated).
// parseWeekdays парсит номера дней недели из строки (1-7, разделенные запятыми).
func parseWeekdays(s string) ([8]bool, error) {
	var days [8]bool
	items := strings.Split(s, ",")
	for _, it := range items {
		if it == "" {
			return days, fmt.Errorf("invalid weekday")
		}
		val, err := strconv.Atoi(it)
		if err != nil || val < 1 || val > 7 {
			return days, fmt.Errorf("invalid weekday")
		}
		days[val] = true
	}
	return days, nil
}

// parseMonthDays parses day numbers from string (1-31, -1, -2, comma-separated).
// parseMonthDays парсит номера дней месяца из строки (1-31, -1, -2, разделенные запятыми).
func parseMonthDays(s string) (map[int]bool, error) {
	m := make(map[int]bool)
	items := strings.Split(s, ",")
	for _, it := range items {
		if it == "" {
			return nil, fmt.Errorf("invalid month day")
		}
		val, err := strconv.Atoi(it)
		if err != nil {
			return nil, fmt.Errorf("invalid month day")
		}
		if val == -1 || val == -2 || (val >= 1 && val <= 31) {
			m[val] = true
			continue
		}
		return nil, fmt.Errorf("invalid month day")
	}
	return m, nil
}

// parseMonths parses month numbers from string (1-12, comma-separated).
// parseMonths парсит номера месяцев из строки (1-12, разделенные запятыми).
func parseMonths(s string) ([13]bool, error) {
	var months [13]bool
	items := strings.Split(s, ",")
	for _, it := range items {
		if it == "" {
			return months, fmt.Errorf("invalid month")
		}
		val, err := strconv.Atoi(it)
		if err != nil || val < 1 || val > 12 {
			return months, fmt.Errorf("invalid month")
		}
		months[val] = true
	}
	return months, nil
}

// matchMonthDay checks if date matches any of the specified month days.
// matchMonthDay проверяет, соответствует ли дата одному из указанных дней месяца.
func matchMonthDay(t time.Time, days map[int]bool) bool {
	day := t.Day()
	last := lastDayOfMonth(t)
	if days[day] {
		return true
	}
	if days[-1] && day == last {
		return true
	}
	if days[-2] && day == last-1 {
		return true
	}
	return false
}

// lastDayOfMonth returns the last day number of the month.
// lastDayOfMonth возвращает номер последнего дня месяца.
func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
