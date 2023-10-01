package service

import (
	"fmt"
	"strconv"
	"time"
)

type Hello struct {
}

func (h *Hello) Hello(msg string) (string, error) {
	msg_int, err := strconv.Atoi(msg)
	if err != nil {
		err = ErrNew(err, SysErr)
		return "", err
	}
	return fmt.Sprintf("hello %v times", msg_int), nil
}

func (h *Hello) HelloTime(date time.Time) string {
	return fmt.Sprintf("tomorrow is %v", date.Format("2006-01-02"))
}
