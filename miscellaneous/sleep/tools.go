package sleep

import(
	"errors"
	"fmt"
	"strings"
	"time"
)

// sum up durations
func Sum(durations ...time.Duration) (sum time.Duration) {
	for _, duration := range durations {
		sum += duration
	}; return
}

// parse string, give back matching format, and time
func ParseUntil(unparsed string) (format string, result time.Time, e error) {
	unparsed = strings.ToUpper(unparsed); date := false; if now.IsZero() { now = time.Now() }

	// first parse times
	if until, e := time.ParseInLocation(STAMP1, unparsed, now.Location()); e == nil {
		until = until.AddDate(now.Year(), int(now.Month()) -1, now.Day() -1); result = until; format = STAMP1
	} else if until, e := time.ParseInLocation(STAMP2, unparsed, now.Location()); e == nil {
		until = until.AddDate(now.Year(), int(now.Month()) -1, now.Day() -1); result = until; format = STAMP2
	} else if until, e := time.ParseInLocation(STAMP3, unparsed, now.Location()); e == nil {
		until = until.AddDate(now.Year(), int(now.Month()) -1, now.Day() -1); result = until; format = STAMP3
	} else if until, e := time.ParseInLocation(STAMP4, unparsed, now.Location()); e == nil {
		until = until.AddDate(now.Year(), int(now.Month()) -1, now.Day() -1); result = until; format = STAMP4
	} else if until, e := time.ParseInLocation(STAMP5, unparsed, now.Location()); e == nil {
		until = until.AddDate(now.Year(), int(now.Month()) -1, now.Day() -1); result = until; format = STAMP5
	} else if until, e := time.ParseInLocation(STAMP6, unparsed, now.Location()); e == nil {
		until = until.AddDate(now.Year(), int(now.Month()) -1, now.Day() -1); result = until; format = STAMP6

		/* reserved */

	} else if  until, e := time.ParseInLocation(STAMP100, unparsed, now.Location()); e == nil {
		result = until; format = STAMP100; date = true
	} else if until, e := time.ParseInLocation(STAMP101, unparsed, now.Location()); e == nil {
		result = until; format = STAMP101; date = true
	} else if until, e := time.ParseInLocation(STAMP102, unparsed, now.Location()); e == nil {
		result = until; format = STAMP102; date =true
	} else if until, e := time.ParseInLocation(STAMP103, unparsed, now.Location()); e == nil {
		result = until; format = STAMP103; date =true
	} else if until, e := time.ParseInLocation(STAMP104, unparsed, now.Location()); e == nil {
		result = until; format = STAMP104; date = true
	} else if until, e := time.ParseInLocation(STAMP105, unparsed, now.Location()); e == nil {
		result = until; format = STAMP105; date = true
	} else if until, e := time.ParseInLocation(STAMP106, unparsed, now.Location()); e == nil {
		result = until; format = STAMP106; date = true
	} else if until, e := time.ParseInLocation(STAMP107, unparsed, now.Location()); e == nil {
		result = until; format = STAMP107; date = true
	} else if until, e := time.ParseInLocation(STAMP108, unparsed, now.Location()); e == nil {
		result = until; format = STAMP108; date = true
	} else if until, e := time.ParseInLocation(STAMP109, unparsed, now.Location()); e == nil {
		result = until; format = STAMP109; date = true
	} else if until, e := time.ParseInLocation(STAMP110, unparsed, now.Location()); e == nil {
		result = until; format = STAMP110; date = true
	} else if until, e := time.ParseInLocation(STAMP111, unparsed, now.Location()); e == nil {
		result = until; format = STAMP111; date = true
	} else if until, e := time.ParseInLocation(STAMP112, unparsed, now.Location()); e == nil {
		result = until; format = STAMP112; date = true
	} else if until, e := time.ParseInLocation(STAMP113, unparsed, now.Location()); e == nil {
		result = until; format = STAMP113; date = true
	} else if until, e := time.ParseInLocation(STAMP114, unparsed, now.Location()); e == nil {
		result = until; format = STAMP114; date = true
	} else if until, e := time.ParseInLocation(STAMP115, unparsed, now.Location()); e == nil {
		result = until; format = STAMP115; date = true
	} else if until, e := time.ParseInLocation(STAMP116, unparsed, now.Location()); e == nil {
		result = until; format = STAMP116; date = true
	} else {
		e = errors.New(fmt.Sprintf("cannot parse '%v' as until time!", unparsed))
		return format, result, e
	}

again:
	// now validate
	if !now.Before(result) {
		if date {
			e = FluxCapacitorMalfunction(result)
		} else {
			result = result.AddDate(0, 0, 1); goto again
		}
	}

	return
}

func FluxCapacitorMalfunction(past time.Time) error {
	return errors.New(fmt.Sprintf("time %v is before now! We just passed it. \n" +
					"When? - Just now. - Well, go back to then. - We can't. - Why not?\n" +
					"We already passed it. - When will then be now? - Soon." , past.Format(STAMP104)))
}
