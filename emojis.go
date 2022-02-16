package log

// (R,Y,G) (3,4,5) (Error,Warning,Okay) can
// be summary items for execution checkpoints.
// Grn, Ylw, Red are calm B/G indicator lights .
const (
	// NOTIFICATION / SUMMARY
	EmojiPanic   = "❌"  // 2 X
	EmojiError   = "❌"  // 3 R
	EmojiWarning = "🟨"  // 4 Y
	EmojOkay     = "🟩"  // 5 G
	EmojiInfo    = "ℹ️" // 6 I
	// TRANSIENT
	EmojiProgress = "▫️" // 7
	EmojiDbg      = "❓"  // misspelled cos 8 != RFC5424 "7"
	/* STATE INDICATORS
	Red = "🔴"
	Ylw = "🟡"
	Grn = "🟢"
	*/
)

/*  This    RFC5424
0   -      Emergency (system is unusable)
1   -      Alert (take action ASAP)
2  Panic   Critical
3  Error   Error
4  Warning Warning
5  Okay    Notice (normal but significant condition)
6  Info    Informational
7  Debug   Debug
*/

func EmojiOfLevel(L Level) string {
	switch L {
	case 0, 1, 2:
		return "💀❌💀"
	case 3:
		return "❌"
	case 4:
		return "🟨"
	case 5:
		return "🟩"
	case 6:
		return "💬"
	case 7:
		return "〰️"
	case 8:
		return "❓"
	}
	return "?!?!"
}

/*
Stockpile of useful emojis:
⭕ ✅ ❌ ❎
🔴 🟠 🟡 🟢 🔵 🟣 🟤 ⚫ ⚪
🟥 🟧 🟨 🟩 🟦 🟪 🟫 ⬛ ⬜ ◾ ◽
🔶 🔷 🔸 🔹 🔺 🔻 💠 🔘 🔳 🔲
*/
