package video

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/playlist"
	"github.com/ystv/web-api/services/creator/types/series"
	"github.com/ystv/web-api/services/creator/types/video"
	"github.com/ystv/web-api/services/encoder"
)

// Store encapsulates our dependencies
type Store struct {
	db   *sqlx.DB
	cdn  *s3.S3
	enc  encoder.Repo
	conf *creator.Config
}

func getSeason(t time.Time) string {
	m := int(t.Month())
	switch {
	case m >= 9 && m <= 12:
		return "aut"
	case m >= 1 && m <= 6:
		return "spr"
	default:
		return "sum"
	}
}

func FileDBToFile(fileDB video.FileDB) video.File {
	var size *int64
	if fileDB.Size.Valid {
		size = &fileDB.Size.Int64
	}

	return video.File{
		URI:          fileDB.URI,
		EncodeFormat: fileDB.EncodeFormat,
		Status:       fileDB.Status,
		Size:         size,
		MimeType:     fileDB.MimeType,
	}
}

func PresetDBToPreset(presetDB video.PresetDB) video.Preset {
	var presetID *int64
	var presetName *string
	if presetDB.PresetID.Valid {
		presetID = &presetDB.PresetID.Int64
	}
	if presetDB.PresetName.Valid {
		presetName = &presetDB.PresetName.String
	}
	return video.Preset{
		PresetID:   presetID,
		PresetName: presetName,
	}
}

func MetaDBToMeta(metaDB video.MetaDB) video.Meta {
	var updatedByID, deletedByID, presetID *int64
	var updatedByNick, deletedByNick, presetName *string
	var updatedAt, deletedAt *time.Time
	if metaDB.UpdatedByID.Valid {
		updatedByID = &metaDB.UpdatedByID.Int64
	}
	if metaDB.DeletedByID.Valid {
		deletedByID = &metaDB.DeletedByID.Int64
	}
	if metaDB.UpdatedAt.Valid {
		updatedByNick = &metaDB.UpdatedByNick.String
	}
	if metaDB.DeletedByNick.Valid {
		deletedByNick = &metaDB.DeletedByNick.String
	}
	if metaDB.UpdatedAt.Valid {
		updatedAt = &metaDB.UpdatedAt.Time
	}
	if metaDB.DeletedAt.Valid {
		deletedAt = &metaDB.DeletedAt.Time
	}
	if metaDB.PresetID.Valid {
		presetID = &metaDB.PresetID.Int64
	}
	if metaDB.PresetName.Valid {
		presetName = &metaDB.PresetName.String
	}

	return video.Meta{
		ID:          metaDB.ID,
		SeriesID:    metaDB.SeriesID,
		Name:        metaDB.Name,
		URL:         metaDB.URL,
		Description: metaDB.Description,
		Thumbnail:   metaDB.Thumbnail,
		Duration:    metaDB.Duration,
		Views:       metaDB.Views,
		Tags:        metaDB.Tags,
		Status:      metaDB.Status,
		Preset: video.Preset{
			PresetID:   presetID,
			PresetName: presetName,
		},
		BroadcastDate: metaDB.BroadcastDate,
		CreatedAt:     metaDB.CreatedAt,
		CreatedByID:   metaDB.CreatedByID,
		CreatedByNick: metaDB.CreatedByNick,
		UpdatedAt:     updatedAt,
		UpdatedByID:   updatedByID,
		UpdatedByNick: updatedByNick,
		DeletedAt:     deletedAt,
		DeletedByID:   deletedByID,
		DeletedByNick: deletedByNick,
	}
}

func ItemDBToItem(itemDB video.ItemDB) video.Item {
	files := make([]video.File, 0)
	for _, file := range itemDB.Files {
		files = append(files, FileDBToFile(file))
	}

	var updatedByID, deletedByID, presetID *int64
	var updatedByNick, deletedByNick, presetName *string
	var updatedAt, deletedAt *time.Time
	if itemDB.UpdatedByID.Valid {
		updatedByID = &itemDB.UpdatedByID.Int64
	}
	if itemDB.DeletedByID.Valid {
		deletedByID = &itemDB.DeletedByID.Int64
	}
	if itemDB.UpdatedAt.Valid {
		updatedByNick = &itemDB.UpdatedByNick.String
	}
	if itemDB.DeletedByNick.Valid {
		deletedByNick = &itemDB.DeletedByNick.String
	}
	if itemDB.UpdatedAt.Valid {
		updatedAt = &itemDB.UpdatedAt.Time
	}
	if itemDB.DeletedAt.Valid {
		deletedAt = &itemDB.DeletedAt.Time
	}
	if itemDB.PresetID.Valid {
		presetID = &itemDB.PresetID.Int64
	}
	if itemDB.PresetName.Valid {
		presetName = &itemDB.PresetName.String
	}

	return video.Item{
		Meta: video.Meta{
			ID:          itemDB.ID,
			SeriesID:    itemDB.SeriesID,
			Name:        itemDB.Name,
			URL:         itemDB.URL,
			Description: itemDB.Description,
			Thumbnail:   itemDB.Thumbnail,
			Duration:    itemDB.Duration,
			Views:       itemDB.Views,
			Tags:        itemDB.Tags,
			Status:      itemDB.Status,
			Preset: video.Preset{
				PresetID:   presetID,
				PresetName: presetName,
			},
			BroadcastDate: itemDB.BroadcastDate,
			CreatedAt:     itemDB.CreatedAt,
			CreatedByID:   itemDB.CreatedByID,
			CreatedByNick: itemDB.CreatedByNick,
			UpdatedAt:     updatedAt,
			UpdatedByID:   updatedByID,
			UpdatedByNick: updatedByNick,
			DeletedAt:     deletedAt,
			DeletedByID:   deletedByID,
			DeletedByNick: deletedByNick,
		},
		Files: files,
	}
}

func PlaylistDBToPlaylist(playlistDB playlist.PlaylistDB) playlist.Playlist {
	metas := make([]video.Meta, 0)
	for _, v := range playlistDB.Videos {
		metas = append(metas, MetaDBToMeta(v))
	}

	return playlist.Playlist{
		Meta: playlist.Meta{
			ID:          playlistDB.ID,
			Name:        playlistDB.Name,
			Description: playlistDB.Description,
			Thumbnail:   playlistDB.Thumbnail,
			Status:      playlistDB.Status,
			CreatedAt:   playlistDB.CreatedAt,
			CreatedBy:   playlistDB.CreatedBy,
			UpdatedAt:   playlistDB.UpdatedAt,
			UpdatedBy:   playlistDB.UpdatedBy,
		},
		Videos: metas,
	}
}

func SeriesDBToSeries(seriesDB series.SeriesDB) series.Series {
	metas := make([]video.Meta, 0)
	for _, v := range seriesDB.ChildVideos {
		metas = append(metas, MetaDBToMeta(v))
	}

	return series.Series{
		Meta: series.Meta{
			SeriesID:    seriesDB.SeriesID,
			URL:         seriesDB.URL,
			SeriesName:  seriesDB.SeriesName,
			Description: seriesDB.Description,
			Thumbnail:   seriesDB.Thumbnail,
			Depth:       seriesDB.Depth,
		},
		ImmediateChildSeries: seriesDB.ImmediateChildSeries,
		ChildVideos:          metas,
	}
}
