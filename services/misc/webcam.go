package misc

type (
	// WebcamStatistic represents general watch statistics
	WebcamStatistic struct {
		TopWatchersAllTime []WebcamUser
		TopWatchersMonth   []WebcamUser
		TopWatchersWeek    []WebcamUser
		LongestWatcher     WebcamUser
	}
	// WebcamUser represents a person and a duration of time in seconds.
	WebcamUser struct {
		Name      string `json:"name"`
		WatchTime int    `json:"watchTime"`
	}
)

// WebcamCurrentViewers returns who is currently watching the webcamms
// and how long their session currently is.
func WebcamCurrentViewers() ([]WebcamUser, error) {
	return []WebcamUser{
		{
			Name:      "Rhys",
			WatchTime: 600,
		},
		{
			Name:      "Cow",
			WatchTime: 1600,
		},
	}, nil
}

// WebcamStatistics returns users watch statistics on webcams.
func WebcamStatistics() (WebcamStatistic, error) {
	return WebcamStatistic{
		TopWatchersAllTime: []WebcamUser{
			{
				Name:      "Rhys",
				WatchTime: 200,
			},
		},
		TopWatchersMonth: []WebcamUser{
			{
				Name:      "Rhys",
				WatchTime: 200,
			},
		},
	}, nil
}
