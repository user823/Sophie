package ds

import (
	"bytes"
	goflag "flag"
	"fmt"
	flag "github.com/spf13/pflag"
	"io"
	"strings"
	"text/tabwriter"
)

// 将多个FlagSet合并，每个FlagSet 打印一组信息
type FlagGroup struct {
	elems map[string]*flag.FlagSet
	order []string
}

func NewFlagGroup() *FlagGroup {
	return &FlagGroup{
		elems: map[string]*flag.FlagSet{},
		order: []string{},
	}
}

func (f *FlagGroup) FlagSet(name string) *flag.FlagSet {
	if _, ok := f.elems[name]; !ok {
		f.elems[name] = flag.NewFlagSet(name, flag.ExitOnError)
		f.order = append(f.order, name)
	}
	return f.elems[name]
}

func (f *FlagGroup) Merge(fg *FlagGroup) {
	for _, k := range fg.order {
		if _, ok := f.elems[k]; !ok {
			f.elems[k] = fg.elems[k]
			f.order = append(f.order, k)
		}
	}
}

func (f *FlagGroup) AddGlobalFlags(flag *flag.Flag) {
	//f.FlagSet("global").BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
	f.FlagSet("global").AddFlag(flag)
}

// 分组打印flag的帮助信息
func (f *FlagGroup) PrintFlags(w io.Writer, maxWidth int) {
	for _, name := range f.order {
		fs := f.elems[name]
		fmt.Fprintf(w, "\n%s flags:\n\n", strings.ToUpper(name[:1])+name[1:])
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			var info bytes.Buffer
			if f.Shorthand != "" {
				info.WriteString(fmt.Sprintf("\t-%s", f.Shorthand))
			}
			info.WriteString(fmt.Sprintf("\t--%s\t", f.Name))
			info.WriteString(fmt.Sprintf("%s\n", LimitWidth(f.Usage, maxWidth-info.Len())))
			tw.Write(info.Bytes())
			tw.Flush()
		})
	}
}

func (f *FlagGroup) FlagSets() (res []*flag.FlagSet) {
	for _, name := range f.order {
		res = append(res, f.elems[name])
	}
	return
}

func LimitWidth(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	return s[:maxWidth-3] + "..."
}

/*
	FlagSet 相关
*/
// 初始化flagSet的默认行为
func InitFlags(flags *flag.FlagSet) {
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	flags.AddGoFlagSet(goflag.CommandLine)
	flags.SortFlags = false
}

// WordSepNormalizeFunc changes all flags that contain "_" separators.
func WordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
	if strings.Contains(name, "_") {
		return flag.NormalizedName(strings.ReplaceAll(name, "-", "_"))
	}
	return flag.NormalizedName(name)
}

// 打印flagset中的flag 信息
func PrintFlags(fs *flag.FlagSet) {
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Printf("FLAG: --%s=%q\n", f.Name, f.Value)
	})
}
