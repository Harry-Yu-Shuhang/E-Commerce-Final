package repositories

import (
	"database/sql"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error) //类型如data[0]= map[string]string{"name":  "Bob","email": "bob@example.com",}
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository { //实现接口的构造函数
	return &OrderManagerRepository{
		table:     table,
		mysqlConn: sql,
	}
}

type OrderManagerRepository struct { //为了实现接口而创建的结构体
	table     string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) Conn() error { //连接数据库
	if o.mysqlConn == nil { //判断是否连接,断了则尝试连接
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "`order`" //order与mysql内置表冲突，需加单引号避免冲突。
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int64, err error) { //插入订单
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT " + o.table + " set userID=?, productID=?, orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return productID, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(order.UserID, order.ProductID, order.OrderStatus)
	if err != nil {
		return productID, err
	}
	// productID, err = result.LastInsertId()
	// return productID, err
	return result.LastInsertId() //这样写更简单
}

func (o *OrderManagerRepository) Delete(productID int64) (isOk bool) { //删除订单
	if err := o.Conn(); err != nil {
		return //返回false直接失败
	}
	sql := "delete from " + o.table + " where ID=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(productID) //填入sql语句的变量ID
	if err != nil {
		return
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) error { //更新订单
	if err := o.Conn(); err != nil {
		return err
	}
	sql := "update " + o.table + " set userID=?, productID=?, orderStatus=? where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.UserID, order.ProductID, order.OrderStatus)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderManagerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) { //通过ID查询订单
	if err = o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sql := "select * from " + o.table + " where ID=" + strconv.FormatInt(orderID, 10)
	row, err := o.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.Order{}, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order, err error) { //查询所有订单
	if err = o.Conn(); err != nil {
		return nil, err
	}
	sql := "select * from " + o.table
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) { //查询所有订单及用户信息
	if err = o.Conn(); err != nil {
		return nil, err
	}
	sql := "select o.ID, p.productName, o.orderStatus from imooc.order as o left join product as p on o.productID = p.ID" //imooc是数据库名称，imooc.order是order表
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	OrderMap = common.GetResultRows(rows)
	return
}
