package test

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/pkg/db/sql"
	"github.com/user823/Sophie/pkg/utils"
	"gorm.io/gorm"
	"testing"
	"time"
)

var (
	db *gorm.DB
	e  error
)

func init() {
	db, e = sql.NewDB("mysql", &sql.MysqlConfig{
		Host:                  "127.0.0.1:3306",
		Username:              "sophie",
		Password:              "123456",
		Database:              "sophie",
		MaxIdleConnections:    10,
		MaxOpenConnections:    10,
		MaxConnectionLifeTime: 3600 * time.Second,
		LogLevel:              0,
		Logger:                nil,
	})
	if e != nil {
		panic(e)
	}
}

type TestTable struct {
	Id        int       `gorm:"column:id"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"golumn:password"`
	CreatedAt time.Time `gorm:"column:create_time"`
}

type TestBook struct {
	Id            int         `gorm:"column:id"`
	Title         string      `gorm:"column:title"`
	Isbn          string      `gorm:"column:isbn"`
	Author        string      `gorm:"column:author"`
	PublisherName string      `gorm:"column:publisher_name"`
	Table         TestTable   `gorm:"foreignKey:Author;references:Username"`
	Tables        []TestTable `gorm:"many2many:table_books;foreignKey:Id;joinForeignKey:UserId;references:Id;joinReferences:RoleId"`
}

type TableBook struct {
	UserId int `gorm:"column:user_id"`
	RoleId int `gorm:"column:role_id"`
}

func (t *TestTable) TableName() string {
	return "test"
}

func (t *TestTable) String() string {
	data, _ := jsoniter.Marshal(t)
	return utils.B2s(data)
}

func (t *TableBook) TableName() string {
	return "table_books"
}

func (t *TestBook) TableName() string {
	return "test_book"
}

func (t *TestTable) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, t)
}

func TestQuery(t *testing.T) {
	var result TestTable
	err := db.Where("test.create_time <= ? ", time.Now()).First(&result).Error
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(result)
}

func TestMarshal(t *testing.T) {
	table := &TestTable{Username: "1", Password: "2"}
	str := table.String()
	m := &TestTable{}
	m.Unmarshal(str)
	fmt.Printf("%v\n", m)
}

type result struct {
	Id     int       `gorm:"column:id"`
	Title  string    `gorm:"column:title"`
	Author string    `gorm:"column:author"`
	Table  TestTable `gorm:"embedded;foreignKey:Author;references:Username"`
}

func (r *result) TableName() string {
	return "test_book"
}

func TestJoins(t *testing.T) {

	// embedded 测试
	var results []result
	query := db.Table("test_book a").Joins("left join test b on b.username = a.author").Select("" +
		"a.id, a.title, b.id, b.username, b.password, b.create_time")
	err := query.Find(&results).Error
	if err != nil {
		t.Error(err)
	}
	for i := range results {
		t.Logf("%v", results[i])
	}
	t.Log("---------")

	// 外键测试 （不能同时支持embedded 和 外键)
	var results1 []result
	query1 := db.Table("test_book").Preload("Table")
	err = query1.Find(&results1).Error
	if err != nil {
		t.Error(err)
	}
	for i := range results1 {
		t.Logf("%v", results1[i])
	}
}

func TestForeignKey(t *testing.T) {
	var result []TestBook
	db.Table("test_book").Preload("Table").Find(&result)
	for i := range result {
		t.Logf("%v", result[i].Tables)
	}
}

func TestMany2Many(t *testing.T) {
	var result []TestBook

	// 预加载 TestBook 中的 Tables 字段
	err := db.Preload("Tables").Find(&result).Error
	if err != nil {
		t.Error(err)
	}

	for i := range result {
		// 打印每个 TestBook 对象中的 Tables 字段
		t.Logf("%v", result[i].Tables)
	}
}

func TestSysUserAutoIngregate(t *testing.T) {
	var users []v1.SysUser
	err := db.Preload("Dept").Preload("Roles").Find(&users).Error
	if err != nil {
		t.Errorf(err.Error())
	}

	// Dept 自动填入
	for i := range users {
		t.Logf("%v", users[i].Dept)
	}
	t.Log("-------")
	// Roles 自动填入
	for i := range users {
		t.Logf("%v", users[i].Roles)
	}
}

func TestSysUserQuery(t *testing.T) {
	var result v1.SysUser
	db.Preload("Dept").Preload("Roles").Where("sys_user.user_name = ? and sys_user.del_flag = 0", "amber").First(&result)
	t.Logf("%v", result)
}

func TestSysRoleQuery(t *testing.T) {
	var result v1.SysRole
	err := db.Table("sys_role r").Joins("" +
		"left join sys_user_role ur on ur.role_id = r.role_id").Joins("" +
		"left join sys_user u on u.user_id = ur.user_id").Joins("" +
		"left join sys_dept d on u.dept_id = d.dept_id").Distinct("r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, " +
		"r.menu_check_strictly, r.dept_check_strictly, r.status, r.del_flag, r.create_time, r.remark").Find(&result).Error
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestGormWhere(t *testing.T) {
	var result TestTable
	err := db.Table("test").Where("id = 1").Where(" id > 0 OR id < 3").First(&result).Error
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestTxRollback(t *testing.T) {
	tx := db.Begin()
	tx.Create(&v1.SysUser{UserId: 5, Username: "test"})
	tx.Create(&v1.SysRole{RoleId: 5, RoleName: "test", RoleKey: "test"})
	tx.Rollback()
}

func TestTxCommit(t *testing.T) {
	tx := db.Begin()
	tx.Create(&v1.SysUser{UserId: 5, Username: "test"})
	tx.Create(&v1.SysRole{RoleId: 5, RoleName: "test", RoleKey: "test"})
	tx.Commit()
}

func TestDBSub(t *testing.T) {
	t.Run("test-connection", TestQuery)
	t.Run("test-marshal", TestMarshal)
	t.Run("test-scans", TestJoins)
	t.Run("test-foreignKey", TestForeignKey)
	t.Run("test-many2many", TestMany2Many)

	// SysUser tests
	t.Run("test-autoingregate", TestSysUserAutoIngregate)
	t.Run("test-sysuserquery", TestSysUserQuery)
	t.Run("test-sysrolequery", TestSysRoleQuery)
	t.Run("test-gormwhere", TestGormWhere)

	t.Run("test-TxRollBack", TestTxRollback)
	t.Run("test-TxCommit", TestTxCommit)
}
