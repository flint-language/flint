package color

func RedText(text string) string       { return Color(text, Red) }
func GreenText(text string) string     { return Color(text, Green) }
func YellowText(text string) string    { return Color(text, Yellow) }
func BlueText(text string) string      { return Color(text, Blue) }
func CyanText(text string) string      { return Color(text, Cyan) }
func MagentaText(text string) string   { return Color(text, Magenta) }
func BoldText(text string) string      { return Color(text, Bold) }
func UnderlineText(text string) string { return Color(text, Underline) }
func BoldRed(text string) string       { return Color(text, Bold, Red) }
func BoldGreen(text string) string     { return Color(text, Bold, Green) }
