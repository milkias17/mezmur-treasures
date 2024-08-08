package tgbot

import (
	"database/sql"
	"log"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
	"github.com/milkias17/mezmur-treasures/internal/db"
)

type Client struct {
	rwMux    sync.RWMutex
	userData map[int64]map[string]any
	db       *sql.DB
	artists  []db.Artist
}

func NewClient() *Client {
	artists, err := db.GetAllArtists(nil, nil)
	if err != nil {
		log.Fatalf("Error getting artists: %s", err)
	}
	return &Client{
		userData: map[int64]map[string]any{},
		db:       db.GetDB(),
		artists:  artists,
	}
}

func (c *Client) getUserData(ctx *ext.Context, key string) (any, bool) {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()

	if c.userData == nil {
		return nil, false
	}

	userData, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		return nil, false
	}

	v, ok := userData[key]
	return v, ok
}

func (c *Client) setUserData(ctx *ext.Context, key string, val any) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()

	if c.userData == nil {
		c.userData = map[int64]map[string]any{}
	}

	_, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		c.userData[ctx.EffectiveUser.Id] = map[string]any{}
	}
	c.userData[ctx.EffectiveUser.Id][key] = val
}

func (c *Client) getUserLang(ctx *ext.Context) (string, error) {
	lang, ok := c.getUserData(ctx, "lang")

	if !ok {
		lang, err := db.GetLanguage(ctx.EffectiveUser.Id, c.db)

		if err != nil {
			log.Printf("Error getting language: %s", err)
			return "", err
		}

		c.setUserData(ctx, "lang", lang)

		return lang, nil
	} else {
		return lang.(string), nil
	}
}

func (c *Client) setUserLang(ctx *ext.Context, lang string) error {

	err := db.SetLanguage(lang, ctx.EffectiveChat.Id, c.db)

	if err != nil {
		log.Printf("Error setting language: %s", err)
		return err
	}

	c.setUserData(ctx, "lang", lang)

	return nil
}

func (c *Client) LocalizedEqual(eq string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		chatId := msg.Chat.Id
		lang, err := db.GetLanguage(chatId, c.db)

		if err != nil {
			log.Printf("Error getting language: %s", err)
			lang = "en"
		}

		if lang == "am" {
			return msg.GetText() == GetTranslationByLang(eq, "am")
		}

		return msg.GetText() == eq
	}
}
