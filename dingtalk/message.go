package dingtalk

// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
// Msg for shorten message
const (
	msgText       = "text"       // text type
	msgLink       = "link"       // link type
	msgMarkDown   = "markdown"   // markdown type
	msgActionCard = "actionCard" // ActionCard type
)

type atReceiver struct {
	AtMobiles []string `json:"atMobiles,omitemty"`
	IsAtAll   bool     `json:"isAtAll,omitemty"`
}

type textMsg struct {
	MsgType string      `json:"msgtype"`
	Text    textContent `json:"text"`
	At      atReceiver  `json:"at"`
}

type textContent struct {
	Content string `json:"content"`
}

type linkMsg struct {
	MsgType string      `json:"msgtype"`
	Link    linkContent `json:"link"`
}

type linkContent struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	MsgURL string `json:"messageUrl"`
	PicURL string `json:"picUrl,omitempty"`
}

type markdownMsg struct {
	MsgType  string          `json:"msgtype"`
	Markdown markdownContent `json:"markdown"`
	At       atReceiver      `json:"at"`
}

type markdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type actionCardMsg struct {
	MsgType    string            `json:"msgtype"`
	ActionCard actionCardContent `json:"actionCard"`
}

type actionCardContent struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
	BtnOrientation string `json:"btnOrientation,omitempty"`
	HideAvatar     string `json:"hideAvatar,omitempty"`
}
