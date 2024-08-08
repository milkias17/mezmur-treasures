package main

import (
	"log"
	"net/http"
	"os"
	"time"

	// "time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
	"github.com/milkias17/mezmur-treasures/internal/db"
	"github.com/milkias17/mezmur-treasures/internal/tgbot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TG_API_KEY")

	if token == "" {
		log.Fatal("TG_API_KEY not set")
	}

	log.Printf("Using token: %s", token)
	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout, // Customise the default request timeout here
				APIURL:  gotgbot.DefaultAPIURL,  // As well as the Default API URL here (in case of using local bot API servers)
			},
		},
	})

	if err != nil {
		log.Fatalf("Failed to create new bot: %s", err.Error())
	}

	log.Printf("Created Bot")

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Printf("An error occurred while handling an update: %s\n", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	client := tgbot.NewClient()

	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{
			handlers.NewCommand("start", client.GreetingHandler),
			handlers.NewMessage(client.LocalizedEqual("List Artists"), client.ListArtists),
			handlers.NewMessage(client.LocalizedEqual("Search Lyrics"), client.HandleSearchLyrics),
			handlers.NewMessage(client.LocalizedEqual("Search Artist"), client.GetArtist),
			handlers.NewMessage(client.LocalizedEqual("Change Language"), client.ChooseLanguage),
		},
		map[string][]ext.Handler{
			tgbot.ARTIST: {
				handlers.NewMessage(client.LocalizedEqual("List Artists"), client.ListArtists),
				handlers.NewMessage(
					client.LocalizedEqual("Search Lyrics"),
					client.HandleSearchLyrics,
				),
				handlers.NewMessage(message.Text, client.ListAlbums),
				handlers.NewCallback(callbackquery.Prefix("artist:"), client.ListAlbums),
				handlers.NewCallback(callbackquery.Prefix("page_artist:"), client.PaginateArtists),
				handlers.NewCallback(callbackquery.Prefix("album:"), client.ListTracks),
				handlers.NewCallback(callbackquery.Prefix("track:"), client.GetLyrics),
			},
			tgbot.SEARCH_LYRICS: {
				handlers.NewMessage(message.Text, client.ListLyricsSearchResults),
				handlers.NewCallback(callbackquery.Prefix("track:"), client.GetLyrics),
			},
			tgbot.LANGUAGE: {
				handlers.NewMessage(message.Text, client.SetLanguage),
			},
		},
		&handlers.ConversationOpts{
			Exits: []ext.Handler{
				handlers.NewMessage(client.LocalizedEqual("🏠Home🏠"), client.GreetingHandler)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))

	if os.Getenv("DEV") == "true" {
		err = updater.StartPolling(bot, &ext.PollingOpts{
			DropPendingUpdates: false,
			GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
				Timeout: 9,
				RequestOpts: &gotgbot.RequestOpts{
					Timeout: time.Second * 10,
				},
			},
		})

		if err != nil {
			log.Fatalf("Failed to start polling: %s\n", err.Error())
		}
	} else {
		webhookDomain := os.Getenv("WEBHOOK_DOMAIN")

		if webhookDomain == "" {
			log.Fatalln("WEBHOOK_DOMAIN not set")
		}
		webhookOpts := ext.WebhookOpts{
			ListenAddr:  os.Getenv("WEBHOOK_LISTEN_ADDR"),
			SecretToken: os.Getenv("WEBHOOK_SECRET"),
			KeyFile:     os.Getenv("WEBHOOK_KEY_FILE"),
			CertFile:    os.Getenv("WEBHOOK_CERT_FILE"),
		}
		err = updater.StartWebhook(bot, "mezmur-treasures/", webhookOpts)
		if err != nil {
			log.Fatalf("Failed to start webhook: %s\n", err.Error())
		}

		err = updater.SetAllBotWebhooks(webhookDomain, &gotgbot.SetWebhookOpts{
			MaxConnections:     100,
			DropPendingUpdates: false,
			SecretToken:        os.Getenv("WEBHOOK_SECRET"),
		})

		if err != nil {
			log.Fatalf("Failed to set webhooks: %s\n", err.Error())
		}
	}

	log.Printf("%s has been started....\n", bot.User.Username)

	updater.Idle()

	db.GetDB()
}
