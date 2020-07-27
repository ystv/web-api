package public

import (
	"log"

	"github.com/ystv/web-api/utils"
)

type (
	// Series provides basic information about a series
	// this is useful when you want to know the current series and
	// see it's immediate children.
	Series struct {
		Series               SeriesMeta  `json:"series"`
		ImmediateChildSeries SeriesMeta  `json:"childSeries"`
		ChildVideos          []VideoMeta `json:"videos"`
	}
	// SeriesMeta is used as a children object for a series
	SeriesMeta []struct {
		SeriesID    int    `json:"seriesID" db:"series_id"`
		SeriesName  string `json:"seriesName" db:"series_name"`
		Description string `json:"description" db:"description"`
		Thumbnail   string `json:"thumbnail" db:"thumbnail"`
		Depth       int    `json:"depth" db:"depth"`
	}
)

// SeriesAllChildrenSeriesNoDepth returns all series a chosen series
// without depth
func (s SeriesMeta) SeriesAllChildrenSeriesNoDepth(SeriesID int) error {
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
	return err
}

// SeriesImmediateChildrenSeries returns series directly below the chosen series
func (s SeriesMeta) SeriesImmediateChildrenSeries(SeriesID int) error {
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
	return err
}

// SeriesAll returns all series in the DB including their depth
func (s SeriesMeta) SeriesAll() error {
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
	return err
}

// SeriesAllBelow returns all series below a certain series
func (s SeriesMeta) SeriesAllBelow(SeriesID int) error {
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
	return err
}
