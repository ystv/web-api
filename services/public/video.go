package public

import (
	"log"
	"strconv"

	_ "github.com/lib/pq" // for DB, although likely not needed
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

type (
	// VideoItem represents public info about video item.
	VideoItem struct {
		VideoMeta
		Files []VideoFile `json:"files"`
	}
	// VideoFile represents each file that a video item has stored.
	VideoFile struct {
		URI      string `json:"uri"`
		MimeType string `db:"mime_type" json:"mimeType"`
		Mode     string `db:"mode" json:"mode"`
		Width    int    `db:"width" json:"width"`
		Height   int    `db:"height" json:"height"`
	}
	// VideoMeta represents basic information about the videoitem used for listing.
	VideoMeta struct {
		VideoID       int         `db:"video_id" json:"id"`
		SeriesID      int         `db:"series_id" json:"seriesID"`
		Name          string      `db:"name" json:"name"`
		URL           string      `db:"url" json:"url"`
		Description   null.String `db:"description" json:"description"`
		Thumbnail     null.String `db:"thumbnail" json:"thumbnail"`
		BroadcastDate string      `db:"broadcast_date" json:"broadcastDate"`
		Views         int         `db:"views" json:"views"`
		Duration      null.Int    `db:"duration" json:"duration"`
	}
)

// VideoList returns all video metadata
func VideoList(offset int, page int) (*[]VideoMeta, error) {
	v := []VideoMeta{}
	// TODO Change pagination method
	err := utils.DB.Select(&v,
		`SELECT video_id, series_id, name, url, description, thumbnail,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE status = 'public'
		ORDER BY broadcast_date DESC
		OFFSET $1 LIMIT $2;`, page, offset)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// VideoFind returns a VideoItem, including the files, based on a given VideoItem ID.
func VideoFind(id int) (*VideoItem, error) {
	v := VideoItem{}
	err := utils.DB.Get(&v,
		`SELECT video_id, series_id, name, url, description, thumbnail,
	views, EXTRACT(EPOCH FROM duration)::int AS duration,
	trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date
	FROM video.items
	WHERE video_id = $1
	AND status = 'public'
	LIMIT 1;`, id)
	if err != nil {
		return nil, err
	}
	err = utils.DB.Select(&v.Files,
		`SELECT uri, mime_type, mode, width, height
	FROM video.files
	INNER JOIN video.encode_formats ON id = encode_format
	WHERE status = 'public'
	AND video_id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func isInt(number string) bool {
	if _, err := strconv.Atoi(number); err == nil {
		return true
	}
	return false
}

// VideoOfSeries returns all the videos belonging to a series
func VideoOfSeries(SeriesID int) ([]VideoMeta, error) {
	v := []VideoMeta{}
	err := utils.DB.Select(&v,
		`SELECT video_id, series_id, name, url, description, thumbnail,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE series_id = $1 AND status = 'public'
		ORDER BY series_position ASC;`, SeriesID)
	if err != nil {
		log.Printf("Failed to select VideoOfSeries: %+v", err)
	}
	return v, err
}

// type Credit struct {
// 	Person string `db:"person_name" json:"person"`
// 	Position string `db:"position_name" json:"position"`
// }

// func CreditsOfVideo(videoID int) ([]Credit, error) {
// 	var c []Credit{}
// 	err := utils.DB.Select(&c,
// 		`SELECT CONCAT(person.first_name, ' ', person.last_name) person, pos.name position
// 		FROM `, videoID)
// }

// VideoBreadcrumb returns the absolute path from a VideoID
func VideoBreadcrumb(VideoID int) ([]Breadcrumb, error) {
	var vB Breadcrumb // Video breadcrumb
	err := utils.DB.Get(&vB,
		`SELECT video_id as id, series_id, COALESCE(name, url) as name, url
		FROM video.items
		WHERE video_id = $1`, VideoID)
	if err != nil {
		log.Printf("VideoBreadcrumb failed: %+v", err)
		return nil, err
	}
	sB, err := SeriesBreadcrumb(vB.SeriesID)
	if err != nil {
		return nil, err
	}
	sB = append(sB, vB)

	return sB, err
}
