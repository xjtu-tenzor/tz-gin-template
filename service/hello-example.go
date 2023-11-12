package service

import (
	"fmt"
	"strconv"
	"template/common"
	"time"
)

type Hello struct {
}

func (h *Hello) Hello(msg string) (string, error) {
	msg_int, err := strconv.Atoi(msg)
	if err != nil {
		return "", common.ErrNew(err, common.SysErr)
	}
	return fmt.Sprintf("hello %v times", msg_int), nil
}

func (h *Hello) HelloTime(date time.Time) string {
	return fmt.Sprintf("tomorrow is %v, really?", date.Format("2006-01-02"))
}
