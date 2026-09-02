package log

import LU "github.com/fbaube/logutils"

func (t *ConsoleTarget) StartLogDetailsBlock(sCatg string, E *Entry) {
	t.Process(E)
	di := t.DetailsInfo
	di.DoingDetails = true
	di.MinLogLevel = LU.LevelOkay
	di.Category = sCatg
	di.Subcategory = ""
}

func (t *ConsoleTarget) CloseLogDetailsBlock(s string) {

}

func (t *ConsoleTarget) LogTextQuote(E *Entry, s string) {

}
