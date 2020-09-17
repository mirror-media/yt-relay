package relay

import (
	"context"
	"fmt"

	ytrelay "github.com/mirror-media/yt-relay"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// ServiceV3 implements the VideoRelay interface and provides api for searching videos with youtube sdk v3
type ServiceV3 struct {
	youtubeService *youtube.Service
}

func New(apiKey string) (*ServiceV3, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apikey is empty for youtube service")
	}
	s, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	return &ServiceV3{
		youtubeService: s,
	}, err
}

func (s *ServiceV3) Search(keyword string, options ytrelay.Options) (interface{}, error) {
	return nil, nil
}

func (s *ServiceV3) ListByVideoIDs(videoIDs []string, options ytrelay.Options) (interface{}, error) {
	return nil, nil
}

func (s *ServiceV3) ListPlaylistVideos(playlistID string, options ytrelay.Options) (interface{}, error) {
	return nil, nil
}
