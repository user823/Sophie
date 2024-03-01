package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type stringCaptcha struct {
	driverString *base64Captcha.DriverString
}

func newStringCaptcha() *stringCaptcha {
	strLength := 4
	showLineOptions := base64Captcha.OptionShowSlimeLine
	source := "0123456789abcdefghjkmnpqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	colour := &color.RGBA{240, 240, 246, 246}
	fonts := []string{"wqy-microhei.ttc"}
	return &stringCaptcha{
		driverString: base64Captcha.NewDriverString(imgheight, imgWidth, 4, showLineOptions, strLength, source, colour, base64Captcha.DefaultEmbeddedFonts, fonts),
	}
}

func (c *stringCaptcha) GetType() string {
	return "string"
}

func (c *stringCaptcha) Generate() (string, string, string) {
	uuid, content, ans := c.driverString.GenerateIdQuestionAnswer()
	item, _ := c.driverString.DrawCaptcha(content)
	return uuid, item.EncodeB64string(), ans
}
