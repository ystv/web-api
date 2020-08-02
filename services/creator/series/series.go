package series

import (
	"database/sql"
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
		Depth       int         `json:"depth" db:"depth"`
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
			child.lft BETWEEN parent.lft AND parent.rgt
		GROUP BY child.series_id
		ORDER BY child.lft ASC;`)
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

// FromPath will return a series from a given path
func FromPath(path string) (Series, error) {
	var s Series
	err := utils.DB.Get(&s.SeriesID, `SELECT series_id FROM video.series_paths WHERE path = $1`, path)
	if err != nil {
		// We ignore ErrNoRows since it's not a log worthy error and the path function will generate this eror when used
		if err != sql.ErrNoRows {
			log.Printf("FromPath failed: %+v", err)
		}
		return s, err
	}
	s, err = View(s.SeriesID)
	return s, err
}

type rec struct {
	Name     string
	Depth    int
	Children []rec
}

var GlobSeries []Meta

func Init() {
	GlobSeries, _ = All()
}

func recursive(depth int, indexglb int) ([]rec, int) {
	all := []rec{}
	for index, item := range GlobSeries[indexglb:] {
		temp := rec{}
		if item.Depth < depth {
			return all, index
		}
		temp.Name = item.URL
		temp.Depth = item.Depth
		if index == len(GlobSeries) {
			all = append(all, temp)
			return all, index
		}
		if GlobSeries[index+1].Depth < depth {
			all = append(all, temp)
			return all, index
		}
		if GlobSeries[index+1].Depth > depth {
			temp.Children, index = recursive(item.Depth, index+1)
		}
		all = append(all, temp)
	}
	return all, len(GlobSeries)
}

func ToJSON() []rec {
	Init()
	final, index := recursive(0, 0)
	log.Print(final)
	log.Print(index)
	return final
}
