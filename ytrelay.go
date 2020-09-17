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
	ID         string `form:"id"`         // For YouTube
	MaxResults uint   `form:"maxResults"` // For YouTube
	PageToken  string `form:"pageToken"`  // For YouTube
	// NextPageToken string `form:"nextPageToken"` // For YouTube
	// PrevPageToken string `form:"prevPageToken"` // For YouTube
	Part       string `form:"part"`       // For YouTube
	SafeSearch string `form:"safeSearch"` // For YouTube
	Type       string `form:"type"`
	Order      string `form:"order"` // For YouTube
	Query      string `form:"q"`     // For YouTube
}

type VideoRelay interface {
	Search(keyword string, options Options) (interface{}, error)
	ListByVideoIDs(videoIDs []string, options Options) (interface{}, error)
	ListPlaylistVideos(playlistID string, options Options) (interface{}, error)
}

type APIWhitelist interface {
	ValidateParameters(params []string) bool
	ValidateChannelID(channelID string) bool
	ValidatePlaylistID(playlistID string) bool
}
