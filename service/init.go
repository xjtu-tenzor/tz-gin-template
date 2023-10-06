package service

import "encoding/gob"

func init() {
	gob.Register(UserSession{})
}
