package stats

// TODO move away from these individual type packages

// VideoGlobalStats represents interesting global video statistics
type VideoGlobalStats struct {
	TotalVideos        int `db:"videos" json:"totalVideos"`
	TotalPendingVideos int `db:"pending" json:"totalPendingVideos"`
	TotalVideoHits     int `db:"hits" json:"totalVideoHits"`
	TotalStorageUsed   int `db:"storage" json:"totalStorageUsed"`
}
