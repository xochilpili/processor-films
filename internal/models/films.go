package models

type BitTorrent struct {
	AddedOn                  int64   `json:"added_on"`
	AmountLeft               int64   `json:"amount_left"`
	AutoTmm                  bool    `json:"auto_tmm"`
	Availability             float64 `json:"availability"`
	Category                 string  `json:"category"`
	Completed                int64   `json:"completed"`
	CompletionOn             int64   `json:"completion_on"`
	ContentPath              string  `json:"content_path"`
	DlLimit                  int64   `json:"dl_limit"`
	Dlspeed                  int64   `json:"dlspeed"`
	DownloadPath             string  `json:"download_path"`
	Downloaded               int64   `json:"downloaded"`
	DownloadedSession        int64   `json:"downloaded_session"`
	Eta                      int64   `json:"eta"`
	FLPiecePrio              bool    `json:"f_l_piece_prio"`
	ForceStart               bool    `json:"force_start"`
	Hash                     string  `json:"hash"`
	InactiveSeedingTimeLimit int64   `json:"inactive_seeding_time_limit"`
	InfohashV1               string  `json:"infohash_v1"`
	InfohashV2               string  `json:"infohash_v2"`
	LastActivity             int64   `json:"last_activity"`
	MagnetUri                string  `json:"magnet_uri"`
	MaxInactiveSeedingTime   int64   `json:"max_inactive_seeding_time"`
	MaxRatio                 float64 `json:"max_ratio"`
	MaxSeedingTime           int64   `json:"max_seeding_time"`
	Name                     string  `json:"name"`
	NumComplete              int64   `json:"num_complete"`
	NumIncomplete            int64   `json:"num_incomplete"`
	NumLeechs                int64   `json:"num_leechs"`
	NumSeeds                 int64   `json:"num_seeds"`
	Priority                 int64   `json:"priority"`
	Progress                 float64 `json:"progress"`
	Ratio                    float64 `json:"ratio"`
	RatioLimit               float64 `json:"ratio_limit"`
	SavePath                 string  `json:"save_path"`
	SeedingTime              int64   `json:"seeding_time"`
	SeedingTimeLimit         int64   `json:"seeding_time_limit"`
	SeenComplete             int64   `json:"seen_complete"`
	SeqDl                    bool    `json:"seq_dl"`
	Size                     int64   `json:"size"`
	State                    string  `json:"state"`
	SuperSeeding             bool    `json:"super_seeding"`
	Tags                     string  `json:"tags"`
	TimeActive               int64   `json:"time_active"`
	TotalSize                int64   `json:"total_size"`
	Tracker                  string  `json:"tracker"`
	TrackersCount            int64   `json:"trackers_count"`
	UpLimit                  int64   `json:"up_limit"`
	Uploaded                 int64   `json:"uploaded"`
	UploadedSession          int64   `json:"uploaded_session"`
	Upspeed                  int64   `json:"upspeed"`
}

type FilmItem struct {
	Id       int
	Provider string
	Title    string
	Year     int
	Genres   []string
}

type Torrent struct {
	Provider      string    `json:"provider"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	OriginalTitle string    `json:"original_title"`
	Year          int       `json:"year"`
	Group         string    `json:"group"`
	Torrents      []Torrent `json:"torrents"`
	Season        int       `json:"season,omitempty"`
	Episode       int       `json:"episode,omitempty"`
	Resolution    string    `json:"resolution"`
	Codec         string    `json:"codec,omitempty"`
	Quality       string    `json:"quality"`
	Seeds         int       `json:"seeds"`
	Peers         int       `json:"peers"`
	Size          string    `json:"size"`
	Magnet        string    `json:"magnet"`
}

type Subtitle struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Group       []string `json:"group"`
	Quality     []string `json:"quality"`
	Resolution  []string `json:"resolution"`
	Duration    []string `json:"duration"`
}

type MetadataFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`
}

type TorrentMetadata struct {
	Data struct {
		Announce  []string       `json:"announce"`
		Files     []MetadataFile `json:"files"`
		InfoHash  string         `json:"infoHash"`
		MagnetUri string         `json:"magnetURI"`
		Name      string         `json:"name"`
		Peers     int            `json:"peers"`
		Seeds     int            `json:"seeds"`
	} `json:"data"`
}

type GenericResponse[T any] struct {
	Message string `json:"message"`
	Total   int    `json:"total"`
	Data    []T    `json:"data"`
}

type OperationType int

const (
	FESTIVALS OperationType = iota + 1 // EnumIndex + 1
	POPULAR
)

func (p OperationType) String() string {
	return [...]string{"films_festivals", "films_popular"}[p-1]
}

func (p OperationType) EnumIndex() int {
	return int(p)
}
