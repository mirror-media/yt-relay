package whitelist

import "github.com/mirror-media/yt-relay/config"

type API struct {
	Whitelist config.Whitelists
}

func (api API) ValidateParameters(params []string) bool {
	for _, param := range params {
		if effective, present := api.Whitelist.APIParameters[param]; !present && effective {
			return false
		}
	}
	return true
}

func (api API) ValidateChannelID(channelID string) bool {
	if effective, present := api.Whitelist.ChannelIDs[channelID]; present && effective {
		return true
	}
	return false
}

func (api API) ValidatePlaylistID(playlistID string) bool {
	if effective, present := api.Whitelist.PlaylistIDs[playlistID]; present && effective {
		return true
	}
	return false
}
