package file

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	v12 "github.com/user823/Sophie/api/thrift/file/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/log"
	utils2 "github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"io"
)

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

// Upload godoc
// @Summary	文件上传
// @Param file formData file true "文件"
// @Accept multipart/form-data
// @Produce application/json
// @Router /file/upload [POST]
func (f *FileController) Upload(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("file")
	if err != nil {
		core.Fail(c, "上传文件失败，请联系管理员", nil)
		return
	}
	filename := string(c.FormValue("fileName"))
	if filename == "" {
		filename = "blob"
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	// 检查文件后缀
	ext := utils2.GetExtension(file)
	if !strutil.ContainsAny(ext, strutil.IMAGE_FILE...) {
		core.Fail(c, fmt.Sprintf("文件格式不正确，仅支持%v", strutil.IMAGE_FILE), nil)
		return
	}

	freader, err := file.Open()
	if err != nil {
		core.Fail(c, "文件打开失败，请重新上传", nil)
		return
	}
	data, err := io.ReadAll(freader)
	if err != nil {
		core.Fail(c, "文件读取失败，请重新上传", nil)
		return
	}

	log.Infof("文件名: %s", filename+ext)
	resp, err := rpc.Remoting.Upload(ctx, &v12.UploadRequest{
		Path:   filename + ext,
		Data:   data,
		UserId: info.User.GetUserId(),
	})

	if err != nil {
		log.Error("error: %s", err.Error())
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.File)
}
