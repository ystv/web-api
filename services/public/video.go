package public

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

type (
	// VideoItem represents public info about video item.
	VideoItem struct {
		VideoID       int         `db:"video_id" json:"id"`
		SeriesID      int         `db:"series_id" json:"series_id"`
		Name          string      `db:"name" json:"name"`
		URL           string      `db:"url" json:"url"`
		Description   null.String `db:"description" json:"description"`
		Thumbnail     null.String `db:"thumbnail" json:"thumbnail"`
		BroadcastDate string      `db:"broadcast_date" json:"broadcastDate"`
		Views         int         `db:"views" json:"views"`
		Duration      null.Int    `db:"duration" json:"duration"`
		Files         []VideoFile `json:"files"`
	}
	// VideoFile represents each file that a video item has stored.
	VideoFile struct {
		URI      string `json:"uri"`
		MimeType string `db:"mime_type" json:"mimeType"`
		Width    int    `db:"width" json:"width"`
		Height   int    `db:"height" json:"height"`
	}
	// VideoMeta represents basic information about the videoitem used for listing.
	VideoMeta struct {
		VideoID       int         `db:"video_id" json:"id"`
		SeriesID      int         `db:"series_id" json:"series_id"`
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
		`SELECT uri, mime_type, width, height
	FROM video.files
	INNER JOIN video.encode_formats ON id = encode_format
	WHERE status = 'public'
	AND video_id = $1`, id)
	log.Print(err)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func URLToVideo(url *url.URL) (*VideoItem, error) {
	log.Print(url.Path)
	// Splitting URL
	videoPath := strings.Split(url.Path, "/")
	// If length is 1 and it's a number might be a video
	if len(videoPath) == 1 && isInt(videoPath[0]) {
		videoID, _ := strconv.Atoi(videoPath[0])
		video, err := VideoFind(videoID)
		return video, err
	}
	// from, err := utils.DB.Prepare(``)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func isInt(number string) bool {
	if _, err := strconv.Atoi(number); err == nil {
		return true
	}
	return false
}

func urlToBr(videoPath []string) {
	for depth := 0; depth < len(videoPath); depth++ {
		// If this is not the final crumb, it is parent
		if depth <= len(videoPath) {

		}
	}

}

// VideoOfSeries returns all the vidos belonging to a series
func VideoOfSeries(SeriesID int) ([]VideoMeta, error) {
	v := []VideoMeta{}
	err := utils.DB.Select(&v,
		`SELECT video_id, name, url, description, thumbnail,
		trim(both '"' from to_json(broadcast_date)::text) AS broadcast_date,
		views, EXTRACT(EPOCH FROM duration)::int AS duration
		FROM video.items
		WHERE series_id = $1 AND status = 'public'`, SeriesID)
	if err != nil {
		log.Printf("Failed to select VideoOfSeries: %+v", err)
	}
	return v, err
}
