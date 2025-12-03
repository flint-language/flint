package color

type Style string

const (
	Reset     Style = "\033[0m"
	Bold      Style = "\033[1m"
	Underline Style = "\033[4m"

	Black   Style = "\033[30m"
	Red     Style = "\033[31m"
	Green   Style = "\033[32m"
	Yellow  Style = "\033[33m"
	Blue    Style = "\033[34m"
	Magenta Style = "\033[35m"
	Cyan    Style = "\033[36m"
	White   Style = "\033[37m"
)
