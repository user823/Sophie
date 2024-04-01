package excelutil

import (
	"github.com/xuri/excelize/v2"
	"testing"
)

type UserTest struct {
	Username string `xlsx:"n:用户名;w:30;"`
	Percent  int    `xlsx:"n:份额;s:%"`
	Sex      string `xlsx:"n:性别;exp:0=男,1=女"`
	Role     Role   `xlsx:"inline"`
}

type UserTestPtr struct {
	Username string `xlsx:"n:用户名;w:30;"`
	Percent  int    `xlsx:"n:份额;s:%"`
	Sex      string `xlsx:"n:性别;exp:0=男,1=女"`
	Role     *Role  `xlsx:"inline"`
}

type Role struct {
	Rolename string `xlsx:"n:角色"`
}

func TestExcel(t *testing.T) {
	file := excelize.NewFile()
	users := []UserTest{
		{Username: "test1", Percent: 70, Sex: "0", Role: Role{Rolename: "管理员"}},
		{Username: "test2", Percent: 30, Sex: "1", Role: Role{Rolename: "测试员"}},
	}
	WriteXlsx(file, "users", users)
	file.DeleteSheet("Sheet1")
	// 删除初始表单
	file.SaveAs("testStruct.xlsx")
}

// 测试结构体指针
func TestExcelPtr(t *testing.T) {
	file := excelize.NewFile()
	users := []*UserTestPtr{
		{Username: "test1", Percent: 70, Sex: "0", Role: &Role{Rolename: "管理员"}},
		{Username: "test2", Percent: 30, Sex: "1", Role: &Role{Rolename: "测试员"}},
	}
	WriteXlsx(file, "users", users)
	file.DeleteSheet("Sheet1")
	// 删除初始表单
	file.SaveAs("testPtr.xlsx")
}
