package public

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/public"
)

// StreamMeta contains basic stream information
type StreamMeta struct {
	public.VideoMeta
	Status string
}

// StreamList handles listing current and upcoming livestreams
func StreamList(c echo.Context) error {
	// We can assume that the first item is the primary stream?
	s := []StreamMeta{
		{
			VideoMeta: public.VideoMeta{
				VideoID:       4550,
				SeriesID:      227,
				Name:          "Comedy Night Live",
				URL:           "cnl",
				Description:   "ComedySoc does the most funny comedy!",
				Thumbnail:     "https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg",
				BroadcastDate: time.Now(),
				Views:         58,
			},
			Status: "live",
		},
		{
			VideoMeta: public.VideoMeta{
				VideoID:       4550,
				SeriesID:      227,
				Name:          "Comedy Night Live: After party",
				URL:           "cnl-after",
				Description:   "Join us for a lovely after-party!",
				Thumbnail:     "https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg",
				BroadcastDate: time.Now(),
				Views:         58,
			},
			Status: "scheduled",
		},
		{
			VideoMeta: public.VideoMeta{
				VideoID:       4550,
				SeriesID:      227,
				Name:          "SwimSoc",
				Description:   "Get them in the pool!",
				URL:           "swim",
				Thumbnail:     "https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg",
				BroadcastDate: time.Now(),
				Views:         235,
			},
			Status: "recent",
		},
	}
	return c.JSON(http.StatusOK, s)
}

// StreamFind finds a stream by ID
func StreamFind(c echo.Context) error {
	s := StreamMeta{VideoMeta: public.VideoMeta{
		VideoID:       4550,
		SeriesID:      227,
		Name:          "Comedy Night Live",
		URL:           "cnl",
		Description:   "ComedySoc does the most funny comedy!",
		Thumbnail:     "https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg",
		BroadcastDate: time.Now(),
		Views:         58,
	},
		Status: "live"}
	return c.JSON(http.StatusOK, s)
}

// StreamHome handles returning stream information for the homepage
// This could be absorbed by stream find but always selecting the first item?
func StreamHome(c echo.Context) error {
	s := StreamMeta{VideoMeta: public.VideoMeta{
		VideoID:       4550,
		SeriesID:      227,
		Name:          "Comedy Night Live",
		URL:           "cnl",
		Description:   "ComedySoc does the most funny comedy!",
		Thumbnail:     "https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg",
		BroadcastDate: time.Now(),
		Views:         58,
	},
		Status: "live"}
	return c.JSON(http.StatusOK, s)
}
