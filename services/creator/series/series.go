package series

import (
	"log"

	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

type (
	// Series provides basic information about a series
	// this is useful when you want to know the current series and
	// see it's immediate children.
	Series struct {
		Meta
		ImmediateChildSeries []Meta       `json:"childSeries"`
		ChildVideos          []video.Meta `json:"videos"`
	}
	// Meta is used as a children object for a series
	Meta struct {
		SeriesID    int         `json:"seriesID" db:"series_id"`
		URL         string      `json:"url" db:"url"`
		SeriesName  null.String `json:"seriesName" db:"name"`
		Description null.String `json:"description" db:"description"`
		Thumbnail   null.String `json:"thumbnail" db:"thumbnail"`
		Depth       int         `json:"-" db:"depth"`
	}
)

// View provides the immediate children of children and videos
func View(SeriesID int) (Series, error) {
	s := Series{}
	s, err := Info(SeriesID)
	s.ImmediateChildSeries, err = ImmediateChildrenSeries(SeriesID)
	s.ChildVideos, err = video.OfSeries(SeriesID)
	if err != nil {
		log.Printf("SeriesInfo failed: %+v", err)
	}
	return s, err
}

// Info provides basic information for only the selected series
func Info(SeriesID int) (Series, error) {
	s := Series{}
	err := utils.DB.Get(&s,
		`SELECT series_id, url, name, description, thumbnail
		FROM video.series
		WHERE series_id = $1`, SeriesID)
	if err != nil {
		log.Printf("SeriesInfo failed: %+v", err)
	}
	return s, err
}

// ImmediateChildrenSeries returns series directly below the chosen series
func ImmediateChildrenSeries(SeriesID int) ([]Meta, error) {
	s := []Meta{}
	err := utils.DB.Select(&s,
		`SELECT * from (
			SELECT 
						node.series_id, node.url, node.name, node.description, node.thumbnail,
						(COUNT(parent.*) - (sub_tree.depth + 1)) AS depth
					FROM
						video.series AS node,
						video.series AS parent,
						video.series AS sub_parent,
						(
							SELECT node.series_id, (COUNT(parent.*) - 1) AS depth
							FROM
								video.series AS node,
								video.series AS parent
							WHERE
								node.series_left between parent.series_left and parent.series_right
								and node.series_id = $1
							GROUP BY node.series_id
							ORDER BY node.series_left ASC
						) AS sub_tree
					WHERE
						node.series_left BETWEEN parent.series_left AND parent.series_right
						AND node.series_left BETWEEN sub_parent.series_left AND sub_parent.series_right
						AND sub_parent.series_id = sub_tree.series_id
					GROUP BY node.series_id, sub_tree.depth
					ORDER BY node.series_left asc
			) as queries
			where depth = 1;`, SeriesID)
	if err != nil {
		log.Printf("Failed SeriesImmediateChildren: %+v", err)
	}
	return s, err
}

// All returns all series in the DB including their depth
func All() ([]Meta, error) {
	s := []Meta{}
	err := utils.DB.Select(&s,
		`SELECT
			child.series_id, child.url, child.name, child.description, child.thumbnail,
			(COUNT(parent.*) -1) AS depth
		FROM
			video.series child,
			video.series parent
		WHERE
			child.series_left BETWEEN parent.series_left AND parent.series_right
		GROUP BY child.series_id
		ORDER BY child.series_left ASC;`)
	if err != nil {
		log.Printf("Failed SeriesAll: %+v", err)
	}
	return s, err
}

// AllBelow returns all series below a certain series including depth
func AllBelow(SeriesID int) ([]Meta, error) {
	s := []Meta{}
	err := utils.DB.Select(&s,
		`SELECT 
			node.series_id, node.url node.name, node.description, node.thumbnail,
			(COUNT(parent.*) - (sub_tree.depth + 1)) AS depth
		FROM
			video.series AS node,
			video.series AS parent,
			video.series AS sub_parent,
			(
				SELECT node.series_id, (COUNT(parent.*) - 1) AS depth
				FROM
					video.series AS node,
					video.series AS parent
				WHERE
					node.series_left between parent.series_left and parent.series_right
					and node.series_id = $1
				GROUP BY node.series_id
				ORDER BY node.series_left ASC
			) AS sub_tree
		WHERE
			node.series_left BETWEEN parent.series_left AND parent.series_right
			AND node.series_left BETWEEN sub_parent.series_left AND sub_parent.series_right
			AND sub_parent.series_id = sub_tree.series_id
		GROUP BY node.series_id, sub_tree.depth
		ORDER BY node.series_left asc;`, SeriesID)
	if err != nil {
		log.Printf("Failed SeriesAllBelow: %+v", err)
	}
	return s, err
}

func PathToSeries(path string) (Series, error) {
	var SeriesID int
	err := utils.DB.Get(SeriesID, `SELECT series_id FROM video.series_paths, video_series WHERE path = $1`, path)
	if err != nil {
		log.Printf("PathToSeries failed: %+v", err)
	}
	s, err := View(SeriesID)
	return s, err
}
