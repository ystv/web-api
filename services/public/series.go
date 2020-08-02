package public

import (
	"log"

	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

// TODO add AND to ensure only public is displayed

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
		URL         string      `json:"url" db:"url"`
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
		`SELECT series_id, url, name, description, thumbnail
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
			child.lft
		FROM
			video.series child,
			video.series parent
		WHERE
			child.lft between parent.lft AND parent.rgt
			AND parent.series_id != child.series_id
			AND parent.series_id = $1
		ORDER BY child.lft ASC;
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
								node.lft between parent.lft and parent.rgt
								and node.series_id = $1
							GROUP BY node.series_id
							ORDER BY node.lft ASC
						) AS sub_tree
					WHERE
						node.lft BETWEEN parent.lft AND parent.rgt
						AND node.lft BETWEEN sub_parent.lft AND sub_parent.rgt
						AND sub_parent.series_id = sub_tree.series_id
					GROUP BY node.series_id, sub_tree.depth
					ORDER BY node.lft asc
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
			child.series_id, child.url, child.name, child.description, child.thumbnail,
			(COUNT(parent.*) -1) AS depth
		FROM
			video.series child,
			video.series parent
		WHERE
			child.lft BETWEEN parent.lft AND parent.rgt
		GROUP BY child.series_id
		ORDER BY child.lft ASC;`)
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
					node.lft between parent.lft and parent.rgt
					and node.series_id = $1
				GROUP BY node.series_id
				ORDER BY node.lft ASC
			) AS sub_tree
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt
			AND node.lft BETWEEN sub_parent.lft AND sub_parent.rgt
			AND sub_parent.series_id = sub_tree.series_id
		GROUP BY node.series_id, sub_tree.depth
		ORDER BY node.lft asc;`, SeriesID)
	if err != nil {
		log.Printf("Failed SeriesAllBelow: %+v", err)
	}
	return s, err
}

// SeriesBreadcrumb will return the breadcrumb from SeriesID to root
func SeriesBreadcrumb(SeriesID int) ([]Breadcrumb, error) {
	s := []Breadcrumb{}
	// TODO Need a bool to indicate if series is in URL
	err := utils.DB.Select(&s,
		`SELECT parent.series_id as id, parent.url as url, COALESCE(parent.name, parent.url) as name
		FROM
			video.series node,
			video.series parent
		WHERE
			node.lft BETWEEN parent.lft AND parent.rgt
			AND node.series_id = $1
		ORDER BY parent.lft;`, SeriesID)
	if err != nil {
		log.Printf("BreadcrumbSeries failed: %+v", err)
	}
	return s, err
}
