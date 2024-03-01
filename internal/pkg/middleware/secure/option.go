package secure

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/internal/pkg/code"
	"strings"
)

type Logic int

const (
	AND Logic = iota
	OR
)

const (
	DefaultSeprator = " "
)

type PermissionParserFunc func(str string) []string

type Option func(opts *Options)

type Options struct {
	Logic            Logic
	PermissionParser PermissionParserFunc
	Separator        string
	Forbidden        app.HandlerFunc
}

func DefaultOption() *Options {
	return &Options{
		Logic:            AND,
		PermissionParser: PermissionParserWithSeparator(DefaultSeprator),
		Separator:        DefaultSeprator,
		Forbidden: func(ctx context.Context, c *app.RequestContext) {
			c.AbortWithStatus(code.FORBIDDEN)
		},
	}
}

func PermissionParserWithSeparator(sep string) PermissionParserFunc {
	return func(str string) []string {
		return strings.Split(str, sep)
	}
}

func WithLogic(logic Logic) Option {
	return func(opts *Options) {
		opts.Logic = logic
	}
}

func WithPermissionParser(f PermissionParserFunc) Option {
	return func(opts *Options) {
		opts.PermissionParser = f
	}
}

func WithSeparator(sep string) Option {
	return func(opts *Options) {
		opts.Separator = sep
	}
}

func WithForbiddon(handle app.HandlerFunc) Option {
	return func(opts *Options) {
		opts.Forbidden = handle
	}
}
