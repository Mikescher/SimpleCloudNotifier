package apihighlight

type ErrHighlight int

//goland:noinspection GoSnakeCaseUsage
const (
	NONE            ErrHighlight = -1
	USER_ID         ErrHighlight = 101
	USER_KEY        ErrHighlight = 102
	TITLE           ErrHighlight = 103
	CONTENT         ErrHighlight = 104
	PRIORITY        ErrHighlight = 105
	CHANNEL         ErrHighlight = 106
	SENDER_NAME     ErrHighlight = 107
	USER_MESSAGE_ID ErrHighlight = 108
)
