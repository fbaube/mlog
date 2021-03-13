package log

func (t *ConsoleTarget) StartLogDetailsBlock(sCatg string, E *Entry) {
	t.Process(E)
	di := t.DetailsInfo
	di.DoingDetails = true
	di.MinLogLevel = LevelOkay
	di.Category = sCatg
	di.Subcategory = ""
}

func (t *ConsoleTarget) CloseLogDetailsBlock(s string) {

}

func (t *ConsoleTarget) LogTextQuote(E *Entry, s string) {

}
