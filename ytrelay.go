package ytrelay

// These are the response structure from YouTube search API, but we only relay the response with interface{}
// They are for references only
type ID struct {
	Kind       string  `json:"kind"`
	VideoID    *string `json:"videoId,omitempty"`
	PlaylistID *string `json:"playlistId,omitempty"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

type Thumbnails struct {
	Default Thumbnail `json:"default"`
	Medium  Thumbnail `json:"medium"`
	High    Thumbnail `json:"high"`
}

type Snippet struct {
	PublishedAt          string     `json:"publishedAt"`
	ChannelID            string     `json:"channelId"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Thumbnails           Thumbnails `json:"thumbnails"`
	ChannelTitle         string     `json:"channelTitle"`
	LiveBroadcastContent string     `json:"liveBroadcastContent"`
	PublishTime          string     `json:"publishTime"`
}

type VideoResource struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	ID      ID      `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type PageInfo struct {
	TotalResults   int64 `json:"totalResults"`
	ResultsPerPage int64 `json:"resultsPerPage"`
}

type VideoList struct {
	Kind           string          `json:"kind"`
	Etag           string          `json:"etag"`
	NextPageToken  string          `json:"nextPageToken"`
	RegionCode     string          `json:"regionCode"`
	PageInfo       PageInfo        `json:"pageInfo"`
	VideoResources []VideoResource `json:"items"`
}

type Options struct {
	ChannelID  string `form:"channelId"`  // For YouTube
	EventType  string `form:"eventType"`
	IDs        string `form:"id"`         // For YouTube
	MaxResults int64  `form:"maxResults"` // For YouTube
	Order      string `form:"order"`      // For YouTube
	PageToken  string `form:"pageToken"`  // For YouTube
	Part       string `form:"part"`       // For YouTube
	PlaylistID string `form:"playlistId"` // For YouTube
	Query      string `form:"q"`          // For YouTube
	SafeSearch string `form:"safeSearch"` // For YouTube
	Type       string `form:"type"`       // For YouTube
}

type VideoRelay interface {
	Search(options Options) (resp interface{}, err error)
	ListByVideoIDs(options Options) (resp interface{}, err error)
	ListPlaylistVideos(options Options) (resp interface{}, err error)
}

type APIWhitelist interface {
	// ValidateParameters(options Options) bool
	ValidateChannelID(channelID string) bool
	ValidatePlaylistIDs(playlistID string) bool
}
