package tgbot

import (
	"encoding/json"
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Translation struct {
	Text      string            `json:"text"`
	Languages map[string]string `json:"languages"`
}

var loaded bool = false

var translations map[string]map[string]string

func parseTranslations() {
	jsonFile, err := os.Open("assets/translations.json")
	if err != nil {
		log.Fatalf("Error reading translations.json: %v", err)
	}
	log.Println("Parsed translations.json")
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&translations)
	if err != nil {
		log.Fatalf("Error parsing translations.json: %v", err)
	}
	loaded = true
}

func GetTranslation(text string, ctx *ext.Context, c *Client) string {
	if !loaded {
		parseTranslations()
	}

	lang, err := c.getUserLang(ctx)
	if err != nil {
		log.Printf("Error getting language: %s", err)
		return text
	}

	langTranslations, ok := translations[lang]
	if !ok {
		return text
	}

	translatedText, ok := langTranslations[text]
	if !ok {
		return text
	}
	return translatedText
}

func GetTranslationByLang(text string, lang string) string {
	if !loaded {
		parseTranslations()
	}

	langTranslations, ok := translations[lang]
	if !ok {
		return text
	}

	translatedText, ok := langTranslations[text]
	if !ok {
		return text
	}
	return translatedText
}
