package tgbot

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/milkias17/mezmur-treasures/internal/db"
)

const (
	HOME          string = "home"
	ARTIST        string = "artist"
	SEARCH_LYRICS string = "search"
	LANGUAGE      string = "language"
)

func makeGetTranslation(ctx *ext.Context, c *Client) func(string) string {
	return func(text string) string {
		return GetTranslation(text, ctx, c)
	}
}

func getHomeKeyboard(page string, getTranslation func(string) string) *gotgbot.ReplyKeyboardMarkup {
	var buttons [][]gotgbot.KeyboardButton

	buttons = append(buttons, []gotgbot.KeyboardButton{})

	if page == "home" {
		buttons[0] = append(buttons[0], gotgbot.KeyboardButton{
			Text: getTranslation("List Artists"),
		})
		buttons[0] = append(buttons[0], gotgbot.KeyboardButton{
			Text: getTranslation("Search Lyrics"),
		})
		buttons = append(buttons, []gotgbot.KeyboardButton{
			{
				Text: getTranslation("Search Artist"),
			},
			{
				Text: getTranslation("Change Language"),
			},
		})
	} else if page == "search" {
		buttons[0] = append(buttons[0], gotgbot.KeyboardButton{
			Text: getTranslation("🏠Home🏠"),
		})
	} else if page == "language" {
		buttons[0] = append(buttons[0], gotgbot.KeyboardButton{
			Text: "English",
		})
		buttons = append(buttons, []gotgbot.KeyboardButton{
			{
				Text: "አማርኛ",
			},
		})
	}

	return &gotgbot.ReplyKeyboardMarkup{
		Keyboard: buttons,
	}
}

const PAGE_SIZE = 10

func (c *Client) GetArtist(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.Message.Reply(b, GetTranslation("Enter Name of Artist", ctx, c), nil)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return handlers.NextConversationState(ARTIST)
}

func (c *Client) ListArtists(b *gotgbot.Bot, ctx *ext.Context) error {
	page, ok := c.getUserData(ctx, "page")
	if !ok {
		c.setUserData(ctx, "page", 1)
		page = 1
	}
	artists := c.artists

	btns := make([][]gotgbot.InlineKeyboardButton, 0)
	userLang, err := c.getUserLang(ctx)

	if err != nil {
		log.Printf("Error getting language: %s", err)
		userLang = "en"
	}

	start := (PAGE_SIZE + 1) * (page.(int) - 1)
	end := start + PAGE_SIZE + 1
	if end > len(artists) {
		end = len(artists)
	}
	for i, artist := range artists[start:end] {
		var artistName string
		if userLang == "am" {
			artistName = artist.AmharicName
		} else {
			artistName = artist.Name
		}
		btns = append(btns, []gotgbot.InlineKeyboardButton{{
			Text:         fmt.Sprintf("%d. %s", i+1, artistName),
			CallbackData: fmt.Sprintf("artist:%s", artist.ID),
		}})
	}

	btns = append(btns, []gotgbot.InlineKeyboardButton{})

	if page.(int) > 1 {
		btns[len(btns)-1] = append(btns[len(btns)-1], gotgbot.InlineKeyboardButton{
			Text:         GetTranslation("Previous", ctx, c),
			CallbackData: fmt.Sprintf("page_artist:back"),
		})
	}

	if page.(int)+1*PAGE_SIZE < len(artists) {
		btns[len(btns)-1] = append(btns[len(btns)-1], gotgbot.InlineKeyboardButton{
			Text:         GetTranslation("Next", ctx, c),
			CallbackData: fmt.Sprintf("page_artist:next"),
		})
	}

	if len(btns[len(btns)-1]) == 0 {
		btns = btns[:len(btns)-1]
	}

	opts := &gotgbot.SendMessageOpts{
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: btns,
		},
	}

	err = sendMessage(b, ctx, GetTranslation("Artists", ctx, c), nil, opts)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return handlers.NextConversationState(ARTIST)
}

func (c *Client) PaginateArtists(b *gotgbot.Bot, ctx *ext.Context) error {
	page, ok := c.getUserData(ctx, "page")
	if !ok {
		c.setUserData(ctx, "page", 1)
		page = 1
	}
	artists := c.artists

	operation := ctx.CallbackQuery.Data
	operation = strings.TrimPrefix(operation, "page_artist:")

	switch operation {
	case "next":
		c.setUserData(ctx, "page", page.(int)+1)
		page = page.(int) + 1
	case "back":
		c.setUserData(ctx, "page", page.(int)-1)
		page = page.(int) - 1
	}

	btns := make([][]gotgbot.InlineKeyboardButton, 0)
	start := (PAGE_SIZE + 1) * (page.(int) - 1)
	end := start + PAGE_SIZE + 1
	if end > len(artists) {
		end = len(artists)
	}
	for i, artist := range artists[start:end] {
		btns = append(btns, []gotgbot.InlineKeyboardButton{{
			Text:         fmt.Sprintf("%d. %s", (PAGE_SIZE+1)*(page.(int)-1)+i+1, artist.Name),
			CallbackData: fmt.Sprintf("artist:%s", artist.ID),
		}})
	}

	btns = append(btns, []gotgbot.InlineKeyboardButton{})

	if page.(int) > 1 {
		btns[len(btns)-1] = append(btns[len(btns)-1], gotgbot.InlineKeyboardButton{
			Text:         GetTranslation("Previous", ctx, c),
			CallbackData: fmt.Sprintf("page_artist:back"),
		})
	}

	if (page.(int)+1)*PAGE_SIZE < len(artists) {
		btns[len(btns)-1] = append(btns[len(btns)-1], gotgbot.InlineKeyboardButton{
			Text:         GetTranslation("Next", ctx, c),
			CallbackData: fmt.Sprintf("page_artist:next"),
		})
	}

	if len(btns[len(btns)-1]) == 0 {
		btns = btns[:len(btns)-1]
	}

	opts := &gotgbot.EditMessageTextOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: btns,
		},
	}

	_, _, err := ctx.CallbackQuery.Message.EditText(b, GetTranslation("Artists", ctx, c), opts)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return nil
}

func (c *Client) ListAlbums(b *gotgbot.Bot, ctx *ext.Context) error {
	artistId := ""
	if ctx.CallbackQuery == nil {
		artistName := ctx.Message.Text
		tmp, err := db.GetArtistIdByName(artistName, c.db)

		if err != nil {
			log.Printf("Error getting artist: %s", err)
			ctx.Message.Reply(b, GetTranslation("Artist not found", ctx, c), nil)
			return err
		}
		artistId = tmp
	} else {
		ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text: GetTranslation("Searching for albums", ctx, c),
		})

		artistId = ctx.CallbackQuery.Data
		artistId = strings.TrimPrefix(artistId, "artist:")
	}

	albums, err := db.GetAllAlbums(artistId, nil, nil, c.db)

	if err != nil {
		log.Printf("Error getting all albums: %s", err)
		return err
	}

	userLang, err := c.getUserLang(ctx)

	if err != nil {
		log.Printf("Error getting language: %s", err)
		userLang = "en"
	}

	btns := make([][]gotgbot.InlineKeyboardButton, 0)
	for i, album := range albums {
		var albumTitle string
		if userLang == "am" {
			albumTitle = album.AmharicTitle
		} else {
			albumTitle = album.Title
		}
		btns = append(btns, []gotgbot.InlineKeyboardButton{{
			Text:         fmt.Sprintf("%d. %s", i+1, albumTitle),
			CallbackData: fmt.Sprintf("album:%s", album.ID),
		}})
	}

	if ctx.CallbackQuery != nil {
		opts := &gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: btns,
			},
		}

		_, _, err = ctx.CallbackQuery.Message.EditText(b, GetTranslation("Albums", ctx, c), opts)
	} else {
		_, err = b.SendMessage(ctx.EffectiveMessage.Chat.Id, GetTranslation("Albums", ctx, c), &gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: btns,
			},
		})
	}

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return nil
}

func (c *Client) ListTracks(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: GetTranslation("Searching for tracks", ctx, c),
	})

	albumId := ctx.CallbackQuery.Data
	albumId = strings.TrimPrefix(albumId, "album:")

	tracks, err := db.GetAllTracks(albumId, nil, nil, c.db)

	if err != nil {
		log.Printf("Error getting all tracks: %s", err)
		return err
	}

	if len(tracks) == 0 {
		ctx.CallbackQuery.Message.EditText(b, GetTranslation("No Tracks Found", ctx, c), nil)
		return nil
	}

	userLang, err := c.getUserLang(ctx)

	if err != nil {
		log.Printf("Error getting language: %s", err)
		userLang = "en"
	}

	btns := make([][]gotgbot.InlineKeyboardButton, 0)
	for i, track := range tracks {
		var trackTitle string
		if userLang == "am" {
			trackTitle = track.Amharic
		} else {
			trackTitle = track.Title
		}
		btns = append(btns, []gotgbot.InlineKeyboardButton{{
			Text:         fmt.Sprintf("%d. %s", i+1, trackTitle),
			CallbackData: fmt.Sprintf("track:%s", track.ID),
		}})
	}

	artistId, err := db.GetArtistIdByAlbum(albumId, c.db)
	if err != nil {
		log.Printf("Error getting artist id: %s", err)
		return err
	}

	btns = append(btns, []gotgbot.InlineKeyboardButton{
		{
			Text:         GetTranslation("Back", ctx, c),
			CallbackData: fmt.Sprintf("artist:%s", artistId),
		},
	})

	opts := &gotgbot.EditMessageTextOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: btns,
		},
	}

	_, _, err = ctx.CallbackQuery.Message.EditText(b, GetTranslation("Tracks", ctx, c), opts)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return nil
}

func (c *Client) GetLyrics(b *gotgbot.Bot, ctx *ext.Context) error {
	trackId := ctx.CallbackQuery.Data
	trackId = strings.TrimPrefix(trackId, "track:")

	ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: GetTranslation("Grabbing Lyrics", ctx, c),
	})

	track, err := db.GetTrackByID(trackId, c.db)

	if err != nil {
		log.Printf("Error getting track: %s", err)
		return err
	}

	if track.Track.Lyrics == "" {
		ctx.CallbackQuery.Message.EditText(b, GetTranslation("No Lyrics Found", ctx, c), nil)
		return nil
	}

	userLang, err := c.getUserLang(ctx)

	if err != nil {
		log.Printf("Error getting language: %s", err)
		userLang = "en"
	}

	var artistName, albumTitle, trackTitle string
	if userLang == "am" {
		artistName = track.ArtistAmharicName
		albumTitle = track.AlbumAmharicTitle
		trackTitle = track.Track.Amharic
	} else {
		artistName = track.ArtistName
		albumTitle = track.AlbumTitle
		trackTitle = track.Track.Title
	}

	text := fmt.Sprintf(
		"<b>%s: </b>%s\n<b>%s: </b>%s\n<b>%s: </b>%s\n\n",
		GetTranslationByLang("Artist", userLang),
		artistName,
		GetTranslationByLang("Album", userLang),
		albumTitle,
		GetTranslationByLang("Track", userLang),
		trackTitle,
	)
	text += track.Track.Lyrics
	err = sendMessage(b, ctx, text, nil, &gotgbot.SendMessageOpts{
		ParseMode:   gotgbot.ParseModeHTML,
		ReplyMarkup: getHomeKeyboard("home", makeGetTranslation(ctx, c)),
	})

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	ctx.CallbackQuery.Message.Delete(b, nil)

	return handlers.EndConversation()
}

func (c *Client) HandleSearchLyrics(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		b,
		GetTranslation("Enter phrase to search", ctx, c),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: getHomeKeyboard("search", makeGetTranslation(ctx, c)),
		},
	)
	log.Printf("Replied to message: %s", err)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return handlers.NextConversationState(SEARCH_LYRICS)
}

func (c *Client) ListLyricsSearchResults(b *gotgbot.Bot, ctx *ext.Context) error {
	phrase := ctx.Message.Text

	lyricsChoices, err := db.GetLyricsFromPhrase(phrase, c.db)

	if err != nil {
		log.Printf("Error getting lyrics: %s", err)
		return err
	}

	userLang, err := c.getUserLang(ctx)

	if err != nil {
		log.Printf("Error getting language: %s", err)
		userLang = "en"
	}

	btns := make([][]gotgbot.InlineKeyboardButton, 0)
	for i, choice := range lyricsChoices {
		var artistName, trackTitle string
		if userLang == "am" {
			artistName = choice.ArtistAmharicName
			trackTitle = choice.TrackAmharicTitle
		} else {
			artistName = choice.ArtistName
			trackTitle = choice.TrackTitle
		}
		btns = append(btns, []gotgbot.InlineKeyboardButton{{
			Text:         fmt.Sprintf("%d. %s - %s", i+1, artistName, trackTitle),
			CallbackData: fmt.Sprintf("track:%s", choice.TrackID),
		}})
	}

	if len(lyricsChoices) == 0 {
		ctx.Message.Reply(b, GetTranslation("No Lyrics Found", ctx, c), nil)
		return nil
	}

	_, err = b.SendMessage(
		ctx.EffectiveMessage.Chat.Id,
		GetTranslation("Search Results", ctx, c),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: btns,
			},
		},
	)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return nil
}

func (c *Client) GreetingHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	text := GetTranslation("Welcome to Mezmur Treasures\n\n", ctx, c)
	text += GetTranslation("Use the provided keyboard below👇", ctx, c)

	_, err := ctx.EffectiveMessage.Reply(b, GetTranslation(text, ctx, c), &gotgbot.SendMessageOpts{
		ReplyMarkup: getHomeKeyboard("home", makeGetTranslation(ctx, c)),
	})

	return err
}

func (c *Client) ChooseLanguage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		b,
		GetTranslation("Choose Language", ctx, c),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: getHomeKeyboard("language", makeGetTranslation(ctx, c)),
		},
	)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return handlers.NextConversationState(LANGUAGE)
}

func (c *Client) SetLanguage(b *gotgbot.Bot, ctx *ext.Context) error {
	language := ctx.Message.Text

	if language != "English" && language != "አማርኛ" {
		ctx.Message.Reply(b, GetTranslation("Invalid Language", ctx, c), nil)
		return nil
	}

	var tmp string
	if language == "English" {
		tmp = "en"
	} else {
		tmp = "am"
	}

	err := c.setUserLang(ctx, tmp)

	if err != nil {
		log.Printf("Error setting language: %s", err)
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(
		b,
		GetTranslation("Language has been set", ctx, c),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: getHomeKeyboard("home", makeGetTranslation(ctx, c)),
		},
	)

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return err
	}

	return handlers.EndConversation()
}
