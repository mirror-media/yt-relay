package relay

import (
	"context"
	"fmt"
	"reflect"
	"strings"

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

// Search supports the following parameters: part, channelId, q, maxResults, pageToken, order, safeSearch
func (s *ServiceV3) Search(options ytrelay.Options) (resp interface{}, err error) {
	yt := s.youtubeService
	call := yt.Search.List(strings.Split(options.Part, ","))
	if !isZero(options.ChannelID) {
		call.ChannelId(options.ChannelID)
	}
	if !isZero(options.Query) {
		call.Q(options.Query)
	}
	if !isZero(options.MaxResults) {
		call.MaxResults(options.MaxResults)
	}
	if !isZero(options.PageToken) {
		call.PageToken(options.PageToken)
	}
	if !isZero(options.Order) {
		call.Order(options.Order)
	}
	if !isZero(options.SafeSearch) {
		call.SafeSearch(options.SafeSearch)
	}

	return call.Do()
}

func (s *ServiceV3) ListByVideoIDs(options ytrelay.Options) (resp interface{}, err error) {
	return nil, nil
}

// ListPlaylistVideos supports the following parameters: part, id, playlistId, maxResults, pageToken
func (s *ServiceV3) ListPlaylistVideos(options ytrelay.Options) (resp interface{}, err error) {
	yt := s.youtubeService
	call := yt.PlaylistItems.List(strings.Split(options.Part, ","))
	if !isZero(options.IDs) {
		call.Id(strings.Split(options.IDs, ",")...)
	}
	if !isZero(options.PlaylistID) {
		call.PlaylistId(options.PlaylistID)
	}
	if !isZero(options.PageToken) {
		call.PageToken(options.PageToken)
	}
	if !isZero(options.MaxResults) {
		call.MaxResults(options.MaxResults)
	}
	return call.Do()
}

func isZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
