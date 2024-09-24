package main

import (
	"context"
	"imooc-product/backend/web/controllers"
	"imooc-product/common"
	"imooc-product/repositories"
	"imooc-product/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板目标
	app.HandleDir("/assets", iris.Dir("./backend/web/assets"))
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

	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productSerivce := services.NewProductService(productRepository)
	productParty := app.Party("/product")              //创建了一个路由组，所有以 /product 开头的请求都将由这个路由组处理。通过域名后面+product来访问这个控制器
	product := mvc.New(productParty)                   //创建了一个新的 MVC 实例，将路由组与 MVC 模式关联。mvc:model,view,controller
	product.Register(ctx, productSerivce)              //注册业务逻辑服务到 MVC 实例。
	product.Handle(new(controllers.ProductController)) //将控制器 ProductController 绑定到 MVC 实例中，以处理请求。

	orderRepository := repositories.NewOrderManagerRepository("`order`", db) //order是mysql内置表，加单引号避免冲突
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		//iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
