package whitelist

import (
	"github.com/mirror-media/yt-relay/config"
)

// YouTubeAPI implements the Whitelist interface
type YouTubeAPI struct {
	Whitelist config.Whitelists
}

func (api YouTubeAPI) ValidateChannelID(channelID string) bool {
	if effective, present := api.Whitelist.ChannelIDs[channelID]; present && effective {
		return true
	}
	return false
}

func (api YouTubeAPI) ValidatePlaylistIDs(playlistID string) bool {
	if effective, present := api.Whitelist.PlaylistIDs[playlistID]; present && effective {
		return true
	}
	return false
}
