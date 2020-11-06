package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// NewSlackMessage は、SlackMessageのコンストラクタです。
func NewSlackMessage() SlackMessage {
	return SlackMessage{}
}

// SendMessage は、メッセージブロックをJson形式へ変換し、Slackへ投稿します。
func (s *SlackMessage) SendMessage(webhookURL string) error {
	payload, err := json.Marshal(s)
	if err != nil {
		return err
	}

	resp, err := http.PostForm(
		webhookURL,
		url.Values{"payload": {string(payload)}},
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to access slack api. result=%v, statusCode=%v", result, resp.StatusCode)
	}

	return nil
}

// AddHeaderBlock は、HeaderBlockを利用したメッセージブロックを追加します。
func (s *SlackMessage) AddHeaderBlock(text string) {
	s.Blocks = append(s.Blocks, Block{
		Type: "header",
		Text: &Text{
			Type: "plain_text",
			Text: text,
		},
	})
}

// AddSectionBlock は、SectionBlockを利用したメッセージブロックを追加します。
func (s *SlackMessage) AddSectionBlock(text string) {
	s.Blocks = append(s.Blocks, Block{
		Type: "section",
		Text: &Text{
			Type: "mrkdwn",
			Text: text,
		},
	})
}

// AddContextImageBlock は、ContentImageBlockを利用したメッセージブロックを追加します。
// ただし、jsonのelementsフィールドがAddActionsButtonBlockと競合してしまうため、現状コメントアウトしている。
// func (s *SlackMessage) AddContextImageBlock(text, imageURL, altText string) {
// 	s.Blocks = append(s.Blocks, Block{
// 		Type: "context",
// 		ContextElements: &[]ContextElement{
// 			ContextElement{
// 				Type:     "image",
// 				ImageURL: &imageURL,
// 				AltText:  &altText,
// 			},
// 			ContextElement{
// 				Type: "mrkdwn",
// 				Text: &text,
// 			},
// 		},
// 	})
// }

// AddActionsButtonBlock は、ActionsButtonBlockを利用したメッセージブロックを追加します。
func (s *SlackMessage) AddActionsButtonBlock(buttonText, linkURL string) {
	s.Blocks = append(s.Blocks, Block{
		Type: "actions",
		ActionsElements: &[]ActionsElement{
			ActionsElement{
				Type: "button",
				Text: &Text{
					Type: "plain_text",
					Text: buttonText,
				},
				URL: &linkURL,
			},
		},
	})
}

// AddDividerBlock は、DividerBlockを利用したメッセージブロックを追加します。
func (s *SlackMessage) AddDividerBlock() {
	s.Blocks = append(s.Blocks, Block{
		Type: "divider",
		Text: nil,
	})
}

// SlackMessage は、Slackへ投稿するメッセージブロックを表す構造体です。
type SlackMessage struct {
	Blocks []Block `json:"blocks"`
}

// Block は、メッセージブロックのBlocksフィールドに対応するデータを表す構造体です。
type Block struct {
	Type string `json:"type"`
	Text *Text  `json:"text,omitempty"`
	// ActionsElementsとJsonのキーが競合してしまうため、現状コメントアウトしている。
	// ContextElements *[]ContextElement `json:"elements,omitempty"`
	ActionsElements *[]ActionsElement `json:"elements,omitempty"`
}

// Text は、メッセージブロックのTextフィールドに対応するデータを表す構造体です。
type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ContextElement は、メッセージブロックのElementsフィールドに対応するデータを表す構造体です。
type ContextElement struct {
	Type     string  `json:"type"`
	Text     *string `json:"text,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
	AltText  *string `json:"alt_text,omitempty"`
}

// ActionsElement は、メッセージブロックのElementsフィールドに対応するデータを表す構造体です。
type ActionsElement struct {
	Type string  `json:"type"`
	Text *Text   `json:"text,omitempty"`
	URL  *string `json:"url,omitempty"`
}
