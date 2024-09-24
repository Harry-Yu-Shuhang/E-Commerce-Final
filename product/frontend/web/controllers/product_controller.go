package controllers

import (
	"imooc-product/datamodels"
	"imooc-product/services"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

var (
	htmlOutPath  = "./frontend/web/htmlProductShow/" //生成html保存目录
	templatePath = "./frontend/web/views/template/"  //静态文件模板目录
)

func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID") //搜索框输入http://localhost:8082/product/generate/html?productID=1,则productString=1
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//1.获取模板文件地址
	contentTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))

	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//2.获取html生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	//3.获取模板渲染数据
	product, err := p.ProductService.GetProductByID(int64(productID)) //productID取决于输入的网址最后productID=几，这里输入1
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//4.生成静态html文件
	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
}

// 生成html静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) { //传入上下文
	//1.判断静态文件是否存在
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}
	//2.生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm) //以只写模式打开文件，如果文件不存在，则创建该文件，并赋予文件完全的权限
	//os.O_CREATE: 这是一个文件打开的标志，表示如果文件不存在则创建文件。如果文件已存在，文件不会被截断或修改。
	//| 是按位或操作符，用于将两个标志组合起来，使得文件可以同时符合这两个条件。
	//os.O_WRONLY:以只写模式打开文件。文件打开后只能向文件中写入数据，不能读取文件中的数据。
	//os.ModePerm: 这是指定文件权限的参数，os.ModePerm 是 0777，表示文件将以完全的读、写、执行权限被创建。
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file, &product)
}

func exist(fileName string) bool { //判断文件是否存在
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail() mvc.View { //这里函数名是Detail，这个功能是product/detail界面的后端
	// id := p.Ctx.URLParam("productID")//这样可以传入id
	product, err := p.ProductService.GetProductByID(1) //写死
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View { //这个功能是product/order界面的后端
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid") //userID
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "Purchase failed." //默认抢购失败
	//判断商品数量是否满足需求
	if product.ProductNum > 0 {
		//扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		//创建订单
		userID, err := strconv.Atoi(userString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamodels.Order{
			UserID:      int64(userID),
			ProductID:   int64(productID),
			OrderStatus: datamodels.OrderSuccess, //int类型，不需要转换,value=1,抢到了就说明已经下单成功了
		}
		//新建订单
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "Purchase success." //抢购成功
		}
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}
}
