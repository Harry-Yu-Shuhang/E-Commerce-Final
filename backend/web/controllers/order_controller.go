package controllers

import (
	"imooc-product/common"
	"imooc-product/datamodels"
	"imooc-product/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}

}

// 修改订单
func (o *OrderController) PostUpdate() {
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(o.Ctx.Request().Form, order); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	err := o.OrderService.UpdateOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}

func (o *OrderController) GetAdd() mvc.View {
	return mvc.View{
		Name: "order/add.html",
	}
}

func (o *OrderController) PostAdd() { //添加
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(o.Ctx.Request().Form, order); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	_, err := o.OrderService.InsertOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}

func (o *OrderController) GetManager() mvc.View { //订单管理界面，负责修改
	idString := o.Ctx.URLParam("id")              //获取id
	id, err := strconv.ParseInt(idString, 10, 16) //由10进制转为16进制
	// fmt.Println("-----------id是---------------")
	// fmt.Println(id) //1
	if err != nil {
		o.Ctx.Application().Logger().Debug(err) //打印error
	}
	order, err := o.OrderService.GetOrderByID(id)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "order/manager.html",
		Data: iris.Map{
			"order": order,
		},
	}
}

func (o *OrderController) GetDelete() { //删除
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	isOk := o.OrderService.DeleteOrderByID(id)
	if isOk {
		o.Ctx.Application().Logger().Debug("删除订单成功，ID为：" + idString)
	} else {
		o.Ctx.Application().Logger().Debug("删除订单失败，ID为：" + idString)
	}
	o.Ctx.Redirect("/order/all")
}
