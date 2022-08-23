package models

type (
	CreateLink struct {
		Link string `json:"link"`
	}
	ResponseBody struct {
		Link  string `json:"link"`
		Short string `json:"short"`
	}
)
