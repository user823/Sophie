package v1

type SysFile struct {
	Name string `json:"name" gorm:"column:name"`
	Url  string `json:"url" gorm:"column:url"`
}
