package tgbot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func getMessageTexts(text string, prefix *string) []string {
	if prefix == nil {
		prefix = new(string)
	}
	messageTexts := make([]string, 0)

	textCopy := text

	for len(textCopy) > 0 {
		if len(textCopy) <= 4096-len(*prefix) {
			messageTexts = append(messageTexts, *prefix+textCopy)
			textCopy = ""
		} else {
			messageTexts = append(messageTexts, *prefix+textCopy[:4096-len(*prefix)])
			textCopy = textCopy[4096-len(*prefix):]
		}
		if len(messageTexts) == 1 {
			prefix = new(string)
		}
	}

	return messageTexts
}

func sendMessage(b *gotgbot.Bot, ctx *ext.Context, fullText string, prefix *string, opts *gotgbot.SendMessageOpts) error {
	texts := getMessageTexts(fullText, prefix)

	for _, text := range texts {
		_, err := ctx.EffectiveChat.SendMessage(b, text, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
