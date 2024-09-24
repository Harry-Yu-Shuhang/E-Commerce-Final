package main

import (
	"context"
	"imooc-product/common"
	"imooc-product/frontend/middleware"
	"imooc-product/frontend/web/controllers"
	"imooc-product/repositories"
	"imooc-product/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板目标
	app.HandleDir("/public", iris.Dir("./frontend/web/public"))
	//访问生成好的html静态文件
	app.HandleDir("/html", iris.Dir("./frontend/web/htmlProductShow"))
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background()) // 创建了一个上下文，用于在应用中传递状态和取消信号。
	defer cancel()

	sess := sessions.New(sessions.Config{
		Cookie:  "helloworld",
		Expires: 60 * time.Minute,
	})

	user := repositories.NewUserManagerRepository("user", db) //user数据库
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx, sess.Start)
	userPro.Handle(new(controllers.UserController))

	//注册product控制器
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("`order`", db) //order数据库需要加单引号，不然会跟系统数据库冲突
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct) //权限设置，只有登陆以后才能点开商品详情
	pro.Register(productService, orderService, sess.Start)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
