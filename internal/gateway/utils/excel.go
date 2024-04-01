package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils/excelutil"
	"github.com/xuri/excelize/v2"
)

func ExportExcel(c *app.RequestContext, sheet string, object any) {
	file := excelize.NewFile()
	excelutil.WriteXlsx(file, sheet, object)
	// 删除默认表
	if err := file.DeleteSheet("Sheet1"); err != nil {
		log.Errorf("Delete default sheet error: %s", err.Error())
	}
	c.SetContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := file.Write(c); err != nil {
		log.Errorf("Write excel file data error: %s", err.Error())
	}
}
