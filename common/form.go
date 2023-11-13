package common

type IDUriForm struct {
	ID int `uri:"id" binding:"min=1"`
}

type PagerForm struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=20"`
}
