package captcha

const DefaultCaptchaType = "string"

const (
	imgWidth  = 150
	imgheight = 50
)

type CaptchaGenerator interface {
	GetType() string
	Generate() (id string, img string, ans string)
}

func NewCaptchaGenerator(captchaType string) CaptchaGenerator {
	switch captchaType {
	case "string":
		return newStringCaptcha()
	case "math":
		return newMathCaptcha()
	default:
		return newStringCaptcha()
	}
}
