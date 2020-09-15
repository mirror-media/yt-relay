package relay

import (
	"context"

	ytrelay "github.com/mirror-media/yt-relay"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	youtubeService *youtube.Service
}

func New(apiKey string) (*Service, error) {
	s, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	return &Service{
		youtubeService: s,
	}, err
}

func (s *Service) Search(keyword string) (ytrelay.VideoList, error) {
	return ytrelay.VideoList{}, nil
}
func (s *Service) ListByVideoID(videoID string) (interface{}, error) {
	return nil, nil
}
func (s *Service) ListPlaylistVideos(playlist string, nextPageToken string) (interface{}, error) {
	return nil, nil
}
func (s *Service) ListChannelVideos(channel string, nextPageToken string) (interface{}, error) {
	return nil, nil
}
