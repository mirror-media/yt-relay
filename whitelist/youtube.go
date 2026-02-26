package whitelist

import (
	"github.com/mirror-media/yt-relay/config"
)

// YouTubeAPI implements the Whitelist interface
type YouTubeAPI struct {
	Whitelist config.Whitelists
}

func (api YouTubeAPI) ValidateChannelID(channelID string) bool {
	effective, present := api.Whitelist.ChannelIDs[channelID]
	return present && effective
}

func (api YouTubeAPI) ValidatePlaylistIDs(playlistID string) bool {
	effective, present := api.Whitelist.PlaylistIDs[playlistID]
	return present && effective
}
