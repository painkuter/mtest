package controller

import (
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

//type PostRequest map[string]interface{}
type PostRequest struct {
	Message string `form:"message"`
}

func Handler(req PostRequest,r render.Render) {
	r.JSON(200, map[string]interface{}{"field": req.Message})
}
