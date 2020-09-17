package whitelist

import (
	"strings"

	ytrelay "github.com/mirror-media/yt-relay"
	"github.com/mirror-media/yt-relay/config"
)

type API struct {
	Whitelist config.Whitelists
}

// func (api API) ValidateParameters(options Options) bool {
// 	for _, param := range params {
// 		if effective, present := api.Whitelist.APIParameters[param]; !present && effective {
// 			return false
// 		}
// 	}
// 	return true
// }

func (api API) ValidateChannelID(options ytrelay.Options) bool {
	if effective, present := api.Whitelist.ChannelIDs[options.ChannelID]; present && effective {
		return true
	}
	return false
}

func (api API) ValidatePlaylistIDs(options ytrelay.Options) bool {
	ids := make([]string, 1)
	if options.PlaylistID != "" {
		ids = append(ids, options.PlaylistID)
	}

	// IDs may contain multiple ids seperated by comma
	if options.IDs != "" {
		ids = append(ids, strings.Split(options.IDs, ",")...)
	}
	for _, id := range ids {
		if effective, present := api.Whitelist.PlaylistIDs[id]; present && effective {
			return true
		}
	}
	return false
}
