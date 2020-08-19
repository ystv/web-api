package public

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/public"
	"gopkg.in/guregu/null.v4"
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
				Description:   null.NewString("ComedySoc does the most funny comedy!", true),
				Thumbnail:     null.NewString("https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg", true),
				BroadcastDate: time.Now().String(),
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
				Description:   null.NewString("Join us for a lovely after-party!", true),
				Thumbnail:     null.NewString("https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg", true),
				BroadcastDate: time.Now().String(),
				Views:         58,
			},
			Status: "scheduled",
		},
		{
			VideoMeta: public.VideoMeta{
				VideoID:       4550,
				SeriesID:      227,
				Name:          "SwimSoc",
				Description:   null.NewString("Get them in the pool!", true),
				URL:           "swim",
				Thumbnail:     null.NewString("https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg", true),
				BroadcastDate: time.Now().String(),
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
		Description:   null.NewString("ComedySoc does the most funny comedy!", true),
		Thumbnail:     null.NewString("https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg", true),
		BroadcastDate: time.Now().String(),
		Views:         58,
	},
		Status: "live"}
	return c.JSON(http.StatusOK, s)
}

// StreamHome handles returning stream information for the homepage
func StreamHome(c echo.Context) error {
	s := StreamMeta{VideoMeta: public.VideoMeta{
		VideoID:       4550,
		SeriesID:      227,
		Name:          "Comedy Night Live",
		URL:           "cnl",
		Description:   null.NewString("ComedySoc does the most funny comedy!", true),
		Thumbnail:     null.NewString("https://ystv.co.uk/static/images/videos/thumbnails/02331.jpg", true),
		BroadcastDate: time.Now().String(),
		Views:         58,
	},
		Status: "live"}
	return c.JSON(http.StatusOK, s)
}
