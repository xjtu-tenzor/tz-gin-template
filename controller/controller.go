package controller

type Controller struct {
	Hello
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}
