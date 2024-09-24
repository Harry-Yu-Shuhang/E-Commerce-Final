package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置模板目标
	app.HandleDir("/public", iris.Dir("./frontend/web/public"))
	//3.访问生成好的html静态文件
	app.HandleDir("/html", iris.Dir("./frontend/web/htmlProductShow"))

	app.Run(
		iris.Addr("0.0.0.0:80"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
