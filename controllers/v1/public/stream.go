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
	}
	return c.JSON(http.StatusOK, s)
}
