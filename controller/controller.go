package controller

type Controller struct {
	Hello
	Std StdController
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}
