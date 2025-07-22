package ui

type FlashStyle int

const (
	Error FlashStyle = iota
	Warning
	Success
)

var styleName = map[FlashStyle]string{
	Error:   "error",
	Warning: "warning",
	Success: "success",
}

type Flash struct {
	Message string
	Style   FlashStyle
}

func (f Flash) GetStyleString() string {
	return styleName[f.Style]
}

func getFlash(message string, style FlashStyle) *Flash {
	return &Flash{Message: message, Style: style}
}

func GetErrorFlash(message string) *Flash {
	return getFlash(message, Error)
}

func GetWarningFlash(message string) *Flash {
	return getFlash(message, Warning)
}

func GetSuccessFlash(message string) *Flash {
	return getFlash(message, Success)
}
