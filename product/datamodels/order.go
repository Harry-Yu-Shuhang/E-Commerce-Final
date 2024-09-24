package datamodels

type Order struct {
	ID          int64 `json:"id" sql:"ID" imooc:"ID"` //在sql中映射成ID字段
	UserID      int64 `json:"user_id" sql:"userID" imooc:"UserID"`
	ProductID   int64 `json:"product_id" sql:"productID" imooc:"ProductID"`
	OrderStatus int   `json:"order_status" sql:"orderStatus" imooc:"OrderStatus"`
}

const (
	OrderWait    = iota //从0开始计数，后面的值不定义则递增，步长1
	OrderSuccess        //1
	OrderFail           //2
)
