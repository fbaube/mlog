package log

/*
// ConsoleTarget writes filtered log messages to console window.
type ConsoleTarget struct {
	*Filter
	ColorMode bool      // whether to use colors to differentiate log levels
	Writer    io.Writer // the writer to write log messages
	close     chan bool
	DetailsInfo
	Category    string
	Subcategory string
}
type DetailsInfo struct {
	AmInDetails, IsLogDetails, IsTextQuote bool
	MinLogLevel                            Level
}
*/

func (t *ConsoleTarget) StartLogDetailsBlock(sCatg string, E *Entry) {
	t.Process(E)
	t.AmInDetails = true
	t.IsLogDetails = true
	t.IsTextQuote = false
	t.MinLogLevel = LevelOkay
	t.Category = sCatg
	t.Subcategory = ""
}

func CloseLogDetailsBlock(s string) {

}
func StartTextQuoteBlock(E *Entry) {

}
func CloseTextQuoteBlock(s string) {

}
