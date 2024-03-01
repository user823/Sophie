package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type PostStore interface {
	// 查询岗位数据集合
	SelectPostList(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) ([]*v1.SysPost, error)
	// 查询所有岗位
	SelectPostAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysPost, error)
	// 通过岗位id查询岗位信息
	SelectPostById(ctx context.Context, postid int64, opts *api.GetOptions) (*v1.SysPost, error)
	// 根据用户id获取岗位选择框列表
	SelectPostListByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]int64, error)
	// 查询用户所属岗位组
	SelectPostsByUserName(ctx context.Context, name string, opts *api.GetOptions) ([]*v1.SysPost, error)
	// 删除岗位信息
	DeletePostById(ctx context.Context, postid int64, opts *api.DeleteOptions) error
	// 批量删除岗位信息
	DeletePostByIds(ctx context.Context, postids []int64, opts *api.DeleteOptions) error
	// 修改岗位信息
	UpdatePost(ctx context.Context, post *v1.SysPost, opts *api.UpdateOptions) error
	// 新增岗位信息
	InsertPost(ctx context.Context, post *v1.SysPost, opts *api.CreateOptions) error
	// 检验岗位名称
	CheckPostNameUnique(ctx context.Context, name string, opts *api.GetOptions) *v1.SysPost
	// 校验岗位编码
	CheckPostCodeUnique(ctx context.Context, postCode string, opts *api.GetOptions) *v1.SysPost
}
