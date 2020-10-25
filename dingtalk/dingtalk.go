package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Robot webhook acess_token
type Robot struct {
	webHook string
	secret  string
}

type rebotRespMsg struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// NewRobot return a new Robot
func NewRobot(webHook string) Sender {
	return &Robot{webHook: webHook}
}

// SetSecret set secret for signature while send request
func (r *Robot) SetSecret(secret string) {
	r.secret = secret
}

// SendText send text type message
func (r Robot) SendText(content string, atMobiles []string, isAtAll bool) error {
	return r.send(&textMsg{
		MsgType: msgText,
		Text: textContent{
			Content: content,
		},
		At: atReceiver{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	})
}

// SendLink send link type message
func (r Robot) SendLink(title, text, messageURL, picURL string) error {
	return r.send(&linkMsg{
		MsgType: msgLink,
		Link: linkContent{
			Title:  title,
			Text:   text,
			MsgURL: messageURL,
			PicURL: picURL,
		},
	})
}

// SendMarkdown send markdown type message
func (r Robot) SendMarkdown(title, text string, atMobiles []string, isAtAll bool) error {
	return r.send(&markdownMsg{
		MsgType: msgMarkDown,
		Markdown: markdownContent{
			Title: title,
			Text:  text,
		},
		At: atReceiver{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	})
}

// SendActionCard send actionCard message
func (r Robot) SendActionCard(title, text, singleTitle, singleURL, btnOrientation, hideAvatar string) error {
	return r.send(&actionCardMsg{
		MsgType: msgActionCard,
		ActionCard: actionCardContent{
			Title:          title,
			Text:           text,
			SingleTitle:    singleTitle,
			SingleURL:      singleURL,
			BtnOrientation: btnOrientation,
			HideAvatar:     hideAvatar,
		},
	})
}

func (r Robot) send(msg interface{}) (err error) {
	m, err := json.Marshal(msg)
	fmt.Println(string(m))
	if err != nil {
		return
	}
	webHook := r.webHook

	if len(r.secret) != 0 {
		webHook = webHook + genSignedURL(r.secret)
	}

	resp, err := http.Post(webHook, "application/json", bytes.NewReader(m))

	if err != nil {
		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rrm rebotRespMsg
	err = json.Unmarshal(data, &rrm)
	if err != nil {
		return err
	}
	if rrm.Errcode != 0 {
		return fmt.Errorf("send failed: %v", rrm.Errmsg)
	}

	return

}

func genSignedURL(secret string) string {
	timeStr := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	sign := fmt.Sprintf("%s\n%s", timeStr, secret)
	signData := computeHmacSha256(sign, secret)
	encodeURL := url.QueryEscape(signData)
	return fmt.Sprintf("&timestamp=%s&sign=%s", timeStr, encodeURL)
}

func computeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
