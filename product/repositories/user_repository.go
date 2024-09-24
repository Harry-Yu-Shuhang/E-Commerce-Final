package repositories

import (
	"database/sql"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"

	"github.com/kataras/iris/v12/x/errors"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userID int64, err error)
}

func NewUserManagerRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table, db}
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "Select * from " + u.table + " where userName =?"
	rows, err := u.mysqlConn.Query(sql, userName)
	if err != nil {
		return &datamodels.User{}, err
	}
	defer rows.Close()
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRow(rows)
	// fmt.Println("result是:")
	// fmt.Println(result) //到这里password还是对的？不对！怎么变成passWord了，mysql数据库的列名称错了！
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在!")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	// fmt.Println("user是:")
	// fmt.Println(user) //改动
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT " + u.table + " SET nickName=?,userName=?,password=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return userId, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return userId, err
	}
	userId, err = result.LastInsertId()
	return
}

func (u *UserManagerRepository) SelectByID(userID int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where ID=" + strconv.FormatInt(userID, 10)
	row, err := u.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.User{}, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在!")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
