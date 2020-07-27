package public

import (
	"log"

	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

type (
	// Series provides basic information about a series
	// this is useful when you want to know the current series and
	// see it's immediate children.
	Series struct {
		SeriesMeta
		ImmediateChildSeries []SeriesMeta `json:"childSeries"`
		ChildVideos          []VideoMeta  `json:"videos"`
	}
	// SeriesMeta is used as a children object for a series
	SeriesMeta struct {
		SeriesID    int         `json:"seriesID" db:"series_id"`
		SeriesName  null.String `json:"seriesName" db:"name"`
		Description null.String `json:"description" db:"description"`
		Thumbnail   null.String `json:"thumbnail" db:"thumbnail"`
		Depth       int         `json:"-" db:"depth"`
	}
)

// SeriesAndChildren provides the immediate children of children and videos
func SeriesAndChildren(SeriesID int) (Series, error) {
	s := Series{}
	s, err := SeriesInfo(SeriesID)
	s.ImmediateChildSeries, err = SeriesImmediateChildrenSeries(SeriesID)
	s.ChildVideos, err = VideoOfSeries(SeriesID)
	if err != nil {
		log.Printf("SeriesInfo failed: %+v", err)
	}
	return s, err
}

// SeriesInfo provides basic information for only the selected series
func SeriesInfo(SeriesID int) (Series, error) {
	s := Series{}
	err := utils.DB.Get(&s,
		`SELECT series_id, name, description, thumbnail
		FROM video.series
		WHERE series_id = $1`, SeriesID)
	if err != nil {
		log.Printf("SeriesInfo failed: %+v", err)
	}
	return s, err
}

// SeriesAllChildrenSeriesNoDepth returns all series a chosen series
// without depth
func SeriesAllChildrenSeriesNoDepth(SeriesID int) ([]SeriesMeta, error) {
	s := []SeriesMeta{}
	err := utils.DB.Select(&s,
		`SELECT
			child.series_id,
			child.url,
			child.series_left
		FROM
			video.series child,
			video.series parent
		WHERE
			child.series_left between parent.series_left AND parent.series_right
			AND parent.series_id != child.series_id
			AND parent.series_id = $1
		ORDER BY child.series_left ASC;
		`, SeriesID)
	if err != nil {
		log.Print(err)
	}
	return s, err
}

// SeriesImmediateChildrenSeries returns series directly below the chosen series
func SeriesImmediateChildrenSeries(SeriesID int) ([]SeriesMeta, error) {
	s := []SeriesMeta{}
	err := utils.DB.Select(&s,
		`SELECT * from (
			SELECT 
						node.series_id, node.name, node.description, node.thumbnail,
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

// SeriesAll returns all series in the DB including their depth
func SeriesAll() ([]SeriesMeta, error) {
	s := []SeriesMeta{}
	err := utils.DB.Select(&s,
		`SELECT
			child.series_id, child.name, child.description, child.thumbnail,
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

// SeriesAllBelow returns all series below a certain series
func SeriesAllBelow(SeriesID int) ([]SeriesMeta, error) {
	s := []SeriesMeta{}
	err := utils.DB.Select(&s,
		`SELECT 
			node.series_id, node.name, node.description, node.thumbnail,
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
