package color

var Enabled = detectTerminal()

func Color(text string, styles ...Style) string {
	if !Enabled {
		return text
	}
	result := ""
	for _, s := range styles {
		result += string(s)
	}
	result += text + string(Reset)
	return result
}
