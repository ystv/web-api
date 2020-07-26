package public

import (
	"log"

	"github.com/ystv/web-api/utils"
)

type (
	Series struct {
		SeriesID int
	}
)

func SeriesChildren(id int) {
	s := Series{}
	err := utils.DB.Select(&s,
		`select
			child.series_id,
			child.url,
			child.series_left
		from
			video.series child,
			video.series parent
		where
			child.series_left between parent.series_left and parent.series_right
			and parent.series_id != child.series_id
			and parent.series_id = 628
		order by child.series_left asc;
		`)
	if err != nil {
		log.Print(err)
	}
}
