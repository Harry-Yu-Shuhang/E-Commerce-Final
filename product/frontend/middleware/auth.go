package middleware

import "github.com/kataras/iris/v12"

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid") //在user_controller.go中设置过cookie的名称为uid
	if uid == "" {
		ctx.Application().Logger().Debug("You have to login first.")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("You have successfully logged in!")
	ctx.Next()
}
