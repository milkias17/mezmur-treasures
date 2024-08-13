package db

import (
	"database/sql"
	"fmt"
)

func GetArtistIdByAlbum(albumID string, db *sql.DB) (string, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	var artistId string
	err := db.QueryRow(`SELECT artist_id FROM albums WHERE id = ?`, albumID).
		Scan(&artistId)

	if err != nil {
		return "", err
	}

	return artistId, nil
}

func GetArtistIdByName(name string, db *sql.DB) (string, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	var artistId string
	var err error
	var field string
	if isEnglish(name) {
		field = "name"
	} else {
		field = "amharic_name"
	}

	err = db.QueryRow(fmt.Sprintf(`SELECT id FROM artists WHERE LOWER(%s) = LOWER(?)`, field), name).
		Scan(&artistId)
	if artistId == "" {
		err = db.QueryRow(fmt.Sprintf(`SELECT id FROM artists WHERE %s LIKE ? OR name LIKE ?`, field), field, name+" %", "% "+name).
			Scan(&artistId)
	}

	if err != nil {
		return "", err
	}

	return artistId, nil
}

func GetArtistByName(name string, db *sql.DB) (Artist, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	var artist Artist
	err := db.QueryRow(`SELECT * FROM artists WHERE name = ?`, name).
		Scan(&artist.ID, &artist.Name, &artist.AmharicName)

	if err != nil {
		return Artist{}, err
	}

	return artist, nil
}

func GetAllArtists(limit *int, offset *int) ([]Artist, error) {
	db := GetDB()
	defer db.Close()

	var artists []Artist
	var rows *sql.Rows
	var err error

	if limit != nil && offset != nil {
		rows, err = db.Query(
			`SELECT * FROM artists ORDER BY name LIMIT ? OFFSET ?`,
			*limit,
			*offset,
		)
	} else if limit != nil {
		rows, err = db.Query(`SELECT * FROM artists ORDER BY name LIMIT ?`, *limit)
	} else if offset != nil {
		rows, err = db.Query(`SELECT * FROM artists ORDER BY name OFFSET ?`, *offset)
	} else {
		rows, err = db.Query(`SELECT * FROM artists ORDER BY name`)
	}

	if err != nil {
		return nil, fmt.Errorf("Error querying artists: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var artist Artist

		if err := rows.Scan(&artist.ID, &artist.Name, &artist.AmharicName); err != nil {
			return nil, fmt.Errorf("Error scanning artist: %s", err)
		}
		artists = append(artists, artist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over artists: %s", err)
	}

	return artists, nil
}

func GetArtistsByPage(page int, pagesize *int) ([]Artist, error) {
	if pagesize == nil {
		pagesize = new(int)
		*pagesize = 10
	}

	offset := (page - 1) * (*pagesize)
	return GetAllArtists(pagesize, &offset)
}

func GetAllAlbums(artistID string, limit *int, offset *int, db *sql.DB) ([]Album, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	var albums []Album
	var rows *sql.Rows
	var err error

	if limit != nil && offset != nil {
		rows, err = db.Query(
			`SELECT * FROM albums WHERE artist_id = ? ORDER BY title LIMIT ? OFFSET ?`,
			artistID,
			*limit,
			*offset,
		)
	} else if limit != nil {
		rows, err = db.Query(`SELECT * FROM albums WHERE artist_id = ? ORDER BY title LIMIT ?`, artistID, *limit)
	} else if offset != nil {
		rows, err = db.Query(`SELECT * FROM albums WHERE artist_id = ? ORDER BY title OFFSET ?`, artistID, *offset)
	} else {
		rows, err = db.Query(`SELECT * FROM albums WHERE artist_id = ? ORDER BY title`, artistID)
	}

	if err != nil {
		return nil, fmt.Errorf("Error querying albums: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var album Album

		if err := rows.Scan(&album.ID, &album.Title, &album.AmharicTitle, &album.Volume, &album.ArtistID); err != nil {
			return nil, fmt.Errorf("Error scanning album: %s", err)
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over albums: %s", err)
	}

	return albums, nil
}

func GetAlbumsByPage(artistID string, page int, pagesize *int) ([]Album, error) {
	if pagesize == nil {
		pagesize = new(int)
		*pagesize = 10
	}

	offset := (page - 1) * (*pagesize)
	return GetAllAlbums(artistID, pagesize, &offset, nil)
}

func GetAllTracks(albumID string, limit *int, offset *int, db *sql.DB) ([]Track, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	var tracks []Track
	var rows *sql.Rows
	var err error

	if limit != nil && offset != nil {
		rows, err = db.Query(
			`SELECT * FROM tracks WHERE album_id = ? ORDER BY title LIMIT ? OFFSET ?`,
			albumID,
			*limit,
			*offset,
		)
	} else if limit != nil {
		rows, err = db.Query(`SELECT * FROM tracks WHERE album_id = ? ORDER BY title LIMIT ?`, albumID, *limit)
	} else if offset != nil {
		rows, err = db.Query(`SELECT * FROM tracks WHERE album_id = ? ORDER BY title OFFSET ?`, albumID, *offset)
	} else {
		rows, err = db.Query(`SELECT * FROM tracks WHERE album_id = ? ORDER BY title`, albumID)
	}

	if err != nil {
		return nil, fmt.Errorf("Error querying tracks: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var track Track

		if err := rows.Scan(&track.ID, &track.Title, &track.Amharic, &track.Lyrics, &track.AlbumID); err != nil {
			return nil, fmt.Errorf("Error scanning track: %s", err)
		}
		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over tracks: %s", err)
	}

	return tracks, nil
}

func GetTracksByPage(albumID string, page int, pagesize *int) ([]Track, error) {
	if pagesize == nil {
		pagesize = new(int)
		*pagesize = 10
	}

	offset := (page - 1) * (*pagesize)
	return GetAllTracks(albumID, pagesize, &offset, nil)
}

type TrackWithMetadata struct {
	Track             Track
	ArtistName        string
	ArtistAmharicName string
	AlbumTitle        string
	AlbumAmharicTitle string
}

func GetTrackByID(id string, db *sql.DB) (TrackWithMetadata, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	var track TrackWithMetadata
	err := db.QueryRow(
		`SELECT 
		tracks.id,
		tracks.title,
		tracks.amharic_title,
		tracks.lyrics,
		tracks.album_id,
		artists.name,
		artists.amharic_name,
		albums.title,
		albums.amharic_title
		FROM tracks
		JOIN albums ON tracks.album_id = albums.id
		JOIN artists ON albums.artist_id = artists.id
		WHERE tracks.id = ?`,
		id,
	).Scan(&track.Track.ID, &track.Track.Title, &track.Track.Amharic, &track.Track.Lyrics, &track.Track.AlbumID, &track.ArtistName, &track.ArtistAmharicName, &track.AlbumTitle, &track.AlbumAmharicTitle)

	if err != nil {
		return TrackWithMetadata{}, err
	}

	return track, nil
}

type LyricsChoice struct {
	ArtistName        string
	ArtistAmharicName string
	AlbumTitle        string
	AlbumAmharicTitle string
	TrackID           string
	TrackTitle        string
	TrackAmharicTitle string
}

func GetLyricsFromPhrase(phrase string, db *sql.DB) ([]LyricsChoice, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	sql := `SELECT 
		artists.name,
		artists.amharic_name,
		albums.title,
		albums.amharic_title,
		tracks.id,
		tracks.title,
		tracks.amharic_title
		FROM tracks
		JOIN albums ON tracks.album_id = albums.id
		JOIN artists ON albums.artist_id = artists.id
		WHERE tracks.lyrics LIKE ?
		ORDER BY artists.name
		LIMIT 25
	`
	rows, err := db.Query(sql, "%"+phrase+"%")
	if err != nil {
		return nil, err
	}

	lyricsChoices := make([]LyricsChoice, 0)

	for rows.Next() {
		var lyricsChoice LyricsChoice

		err := rows.Scan(
			&lyricsChoice.ArtistName,
			&lyricsChoice.ArtistAmharicName,
			&lyricsChoice.AlbumTitle,
			&lyricsChoice.AlbumAmharicTitle,
			&lyricsChoice.TrackID,
			&lyricsChoice.TrackTitle,
			&lyricsChoice.TrackAmharicTitle,
		)
		if err != nil {
			return nil, err
		}
		lyricsChoices = append(lyricsChoices, lyricsChoice)
	}

	defer rows.Close()

	return lyricsChoices, nil
}

func SetLanguage(language string, chatId int64, db *sql.DB) error {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	_, err := db.Exec(
		`INSERT OR REPLACE INTO languagePreference (chat_id, language) VALUES (?, ?)`,
		chatId,
		language,
	)

	return err
}

func GetLanguage(chatId int64, db *sql.DB) (string, error) {
	if db == nil {
		db = GetDB()
		fmt.Println("Passed nil db")
		defer db.Close()
	}

	var language string
	err := db.QueryRow(`SELECT language FROM languagePreference WHERE chat_id = ?`, chatId).
		Scan(&language)

	if err != nil {
		return "", err
	}

	return language, nil
}

func IncreaseCountUsers(chatId string) {
	db := GetDB()
	defer db.Close()

	var count int
	err := db.QueryRow(`SELECT * FROM counter WHERE chat_id = ?`, chatId).
		Scan(&count)
	if err != nil {
		db.Exec(`INSERT INTO counter (chat_id, count) VALUES (?, 1)`, chatId)
	} else {
		db.Exec(`UPDATE counter SET count = count + 1 WHERE chat_id = ?`, chatId)
	}
}
