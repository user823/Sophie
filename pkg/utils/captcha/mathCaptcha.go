package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type mathCaptcha struct {
	driverMath *base64Captcha.DriverMath
}

func newMathCaptcha() *mathCaptcha {
	showLineOptions := base64Captcha.OptionShowSlimeLine
	colour := &color.RGBA{240, 240, 246, 246}
	fonts := []string{"wqy-microhei.ttc"}
	return &mathCaptcha{
		driverMath: base64Captcha.NewDriverMath(imgheight, imgWidth, 4, showLineOptions, colour, base64Captcha.DefaultEmbeddedFonts, fonts),
	}
}

func (c *mathCaptcha) GetType() string {
	return "math"
}

func (c *mathCaptcha) Generate() (string, string, string) {
	uuid, content, ans := c.driverMath.GenerateIdQuestionAnswer()
	item, _ := c.driverMath.DrawCaptcha(content)
	return uuid, item.EncodeB64string(), ans
}
