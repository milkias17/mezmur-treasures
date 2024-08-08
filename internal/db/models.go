package db

type Artist struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	AmharicName string `db:"amharic_name"`
}

type Album struct {
	ID       string `db:"id"`
	Title    string `db:"title"`
	AmharicTitle  string `db:"amharic_title"`
	Volume   *string `db:"volume"`
	ArtistID string `db:"artist_id"`
}

type Track struct {
	ID      string `db:"id"`
	Title   string `db:"title"`
	Amharic string `db:"amharic_title"`
	Lyrics  string `db:"lyrics"`
	AlbumID string `db:"album_id"`
}

