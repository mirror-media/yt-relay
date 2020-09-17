package relay

import (
	"context"

	ytrelay "github.com/mirror-media/yt-relay"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	youtubeService *youtube.Service
	ChannelIDs     []string
	PlaylistIDs    []string
}

func New(apiKey string) (*Service, error) {
	s, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	return &Service{
		youtubeService: s,
	}, err
}

func (s *Service) Search(keyword string, options ytrelay.Options) (interface{}, error) {
	return nil, nil
}

func (s *Service) ListByVideoIDs(videoIDs []string, options ytrelay.Options) (interface{}, error) {
	return nil, nil
}

func (s *Service) ListPlaylistVideos(playlistID string, options ytrelay.Options) (interface{}, error) {
	return nil, nil
}
