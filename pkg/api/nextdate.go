package api

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	tNow, err := time.Parse(DateFormat, now)
	if err != nil {
		tNow = time.Now()
	}
	resDate, err := NextDate(tNow, date, repeat)
	if err != nil {
		slog.Error(err.Error(), "func", nextDayHandler)
	}
	msg := fmt.Sprintf("%s\n", resDate)
	io.WriteString(w, msg)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if dstart == "" {
		slog.Error("NextDate: wrong format", "date", dstart)
		return "", nil
	}
	if repeat == "" {
		slog.Error("NextDate: Wrong format", "repeat", repeat)
		return "", nil
	}
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		slog.Error(err.Error())
		return "", nil
	}
	// Разбор repeat
	sRepeat := strings.Split(repeat, " ")
	if len(sRepeat) >= 1 {
		switch sRepeat[0] {
		// раз в год
		case "y":
			for {
				date = date.AddDate(1, 0, 0)
				if afterNow(date, now) {
					break
				}
			}
			return date.Format(DateFormat), nil
		// через N дней
		case "d":
			if len(sRepeat) == 1 {
				slog.Error("NextDate: Wrong number of days:", "repeat", repeat)
				return "", nil
			}
			days, err := strconv.Atoi(sRepeat[1])
			if err != nil {
				slog.Error("NextDate: Not integer")
				return "", nil
			}
			if days <= 0 || days > 31 {
				slog.Error("NextDate: Wrong day")
				return "", nil
			}

			for {
				date = date.AddDate(0, 0, days)
				if afterNow(date, now) {
					break
				}
			}
			return date.Format(DateFormat), nil
		// заданы дни и месяцы
		case "m":
			var resDate, nextDate time.Time
			switch len(sRepeat) {
			case 1:
				slog.Error("NextDate: Wrong number of month days:", "repeat", repeat)
				return "", nil
			//  Заданы только дни
			case 2:
				arrDay := strings.Split(sRepeat[1], ",")

				for _, day := range arrDay {
					nDay, err := strconv.Atoi(day)
					if err != nil {
						slog.Error("NextDate: Not integer", "day at repeat", day)
						return "", nil
					}
					if nDay == 0 || nDay > 31 || nDay < -2 {
						slog.Error("NextDate: Wrong day", "day at repeat", nDay)
						return "", nil
					}
					nextDate = date
					for {
						if nDay > 0 {
							dayMonth := int(time.Date(nextDate.Year(), nextDate.Month()+1, 1, 0, 0, 0, 0, nextDate.Location()).Sub(time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, nextDate.Location())).Hours() / 24)
							if nDay > dayMonth {
								nextDate = nextDate.AddDate(0, 0, dayMonth)
							}
							nextDate = time.Date(nextDate.Year(), nextDate.Month(), nDay, 0, 0, 0, 0, nextDate.Location())
							if afterNow(nextDate, now) {
								break
							}
							nextDate = nextDate.AddDate(0, 1, 0)
						} else {
							nextDate = time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, nextDate.Location())
							nextDate = nextDate.AddDate(0, 1, nDay)
							if afterNow(nextDate, now) {
								break
							}
						}
					}
					if resDate.IsZero() {
						resDate = nextDate
					} else {
						if resDate.After(nextDate) {
							resDate = nextDate
						}
					}
				}
				return resDate.Format(DateFormat), nil

			//  Заданы дни и месяцы
			case 3:
				arrDay := strings.Split(sRepeat[1], ",")
				arrMonth := strings.Split(sRepeat[2], ",")
				for _, day := range arrDay {
					for _, month := range arrMonth {
						nDay, err := strconv.Atoi(day)
						if err != nil {
							slog.Error("NextDate: Not integer")
							return "", nil
						}
						nMonth, err := strconv.Atoi(month)
						if err != nil {
							slog.Error("NextDate: Not integer")
							return "", nil
						}
						if nDay == 0 || nMonth <= 0 || nDay > 31 || nMonth > 12 || nDay < -2 {
							slog.Error("NextDate: Wrong value day or month")
							return "", nil
						}
						nextDate = date
						for {
							if nDay > 0 {
								nextDate = time.Date(nextDate.Year(), time.Month(nMonth), nDay, 0, 0, 0, 0, nextDate.Location())
								if afterNow(nextDate, now) {
									break
								}
								nextDate = nextDate.AddDate(1, 0, 0)
							} else {
								nextDate = time.Date(nextDate.Year(), time.Month(nMonth), 1, 0, 0, 0, 0, nextDate.Location())
								nextDate = nextDate.AddDate(0, 1, nDay)
								if afterNow(nextDate, now) {
									break
								}
								nextDate = nextDate.AddDate(1, 0, 0)
							}
						}
						if resDate.IsZero() {
							resDate = nextDate
						} else {
							if resDate.After(nextDate) {
								resDate = nextDate
							}
						}
					}
				}
				return resDate.Format(DateFormat), nil

			default:
				slog.Error("NextDate: Wrong number of month days:", "repeat", repeat)
				return "", nil
			}

		// заданы дни недели
		case "w":
			if len(sRepeat) == 1 {
				slog.Error("NextDate: Wrong number of weekdays:", "repeat", repeat)
				return "", nil
			}
			var resDate, nextDate time.Time
			arrDay := strings.Split(sRepeat[1], ",")

			for _, day := range arrDay {
				nDay, err := strconv.Atoi(day)
				if err != nil {
					slog.Error("NextDate: Not integer")
					return "", nil
				}
				if nDay <= 0 || nDay > 7 {
					slog.Error("NextDate: Wrong weekday")
					return "", nil
				}
				nextDate = time.Date(date.Year(), date.Month(), date.Day()+((7+nDay-int(date.Weekday()))%7), 0, 0, 0, 0, date.Location())
				for {
					if afterNow(nextDate, now) {
						break
					}
					nextDate = nextDate.AddDate(0, 0, 7)
				}
				if resDate.IsZero() {
					resDate = nextDate
				} else {
					if resDate.After(nextDate) {
						resDate = nextDate
					}
				}
			}
			return resDate.Format(DateFormat), nil

		default:
			slog.Error("NextDate: Wrong format Repeat")
			return "", nil
		}
	} else {
		slog.Error("NextDate: Wrong format Repeat")
		return "", nil
	}
}

func afterNow(date, now time.Time) bool {
	return date.After(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
}

func checkRepeat(repeat string) (bool, string) {
	if repeat == "" {
		return true, ""
	}
	sRepeat := strings.Split(repeat, " ")
	lenRepeat := len(sRepeat)
	if lenRepeat >= 1 {
		switch sRepeat[0] {
		case "y":
			if lenRepeat > 1 {
				return false, "wrong format Repeat"
			} else {
				return true, ""
			}
		case "d":
			if lenRepeat != 2 {
				return false, "wrong format Repeat"
			} else {
				return true, ""
			}
		case "m":
			if lenRepeat == 1 || lenRepeat > 3 {
				return false, "wrong format Repeat"
			} else {
				return true, ""
			}
		case "w":
			if lenRepeat != 2 {
				return false, "wrong format Repeat"
			} else {
				return true, ""
			}
		default:
			return false, "wrong format Repeat"
		}
	}
	return true, ""
}
