package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	ytrelay "github.com/mirror-media/yt-relay"
	"github.com/mirror-media/yt-relay/api"
	"github.com/mirror-media/yt-relay/relay"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/youtube/v3"
)

const (
	ErrorEmptyPart = "part cannot be empty"
	ErrorEmptyID   = "id cannot be empty"
)

var json jsoniter.API

func init() {
	json = jsoniter.ConfigFastest
}

func Set(r *gin.Engine, relayService ytrelay.VideoRelay, whitelist ytrelay.APIWhitelist) error {

	ytRouter := r.Group("/youtube/v3")

	// search videos. ChannelID is required
	ytRouter.GET("/search", func(c *gin.Context) {

		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		queries, err := parseQueries(c)
		if err != nil {
			apiLogger.Error(err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Check the mandatory parameters
		if queries.Part == "" {
			apiLogger.Error(ErrorEmptyPart)
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: ErrorEmptyPart})
			return
		}

		// Check whitelist
		if !whitelist.ValidateChannelID(queries) {
			err = fmt.Errorf("channelId(%s) is invalid", queries.ChannelID)
			apiLogger.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: ErrorEmptyID})
			return
		}

		resp, err := relayService.Search(queries)
		if err != nil {
			apiLogger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResp{Error: ErrorEmptyID})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	// list video by video id
	// IDs of videos is required
	ytRouter.GET("/video", func(c *gin.Context) {

		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		queries, err := parseQueries(c)
		if err != nil {
			apiLogger.Error(err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Check the mandatory parameters
		if queries.Part == "" {
			apiLogger.Error(ErrorEmptyPart)
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: ErrorEmptyPart})
			return
		}
		if queries.IDs == "" {
			apiLogger.Error(ErrorEmptyID)
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: ErrorEmptyID})
			return
		}

		resp, err := relayService.ListByVideoIDs(queries)
		if err != nil {
			apiLogger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResp{Error: ErrorEmptyID})
			return
		}

		// verify channel id for YouTube
		_, isYouTube := relayService.(*relay.YouTubeServiceV3)
		if isYouTube {
			if err = validateYouTubeVideoListResponse(whitelist, resp); err != nil {
				err = errors.Wrap(err, "some video's channel id is invalid")
				apiLogger.Error(err)
				c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, resp)
	})

	// list video by playlistID
	ytRouter.GET("/playlistItems", func(c *gin.Context) {

		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		queries, err := parseQueries(c)
		if err != nil {
			apiLogger.Error(err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Check the mandatory parameters
		if queries.Part == "" {
			apiLogger.Error(ErrorEmptyPart)
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: ErrorEmptyPart})
			return
		}

		// Check whitelist
		if !whitelist.ValidatePlaylistIDs(queries) {
			err = fmt.Errorf("playlist(%s) or id(%s) is invalid", queries.PlaylistID, queries.IDs)
			apiLogger.Error(err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		resp, err := relayService.ListByVideoIDs(queries)
		if err != nil {
			apiLogger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResp{Error: ErrorEmptyID})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	return nil
}

func parseQueries(c *gin.Context) (ytrelay.Options, error) {
	var queries ytrelay.Options
	err := c.BindQuery(&queries)

	return queries, err
}

func validateYouTubeVideoListResponse(whitelist ytrelay.APIWhitelist, resp interface{}) (err error) {
	for _, item := range resp.(*youtube.VideoListResponse).Items {
		if !whitelist.ValidateChannelID(ytrelay.Options{ChannelID: item.Snippet.ChannelId}) {
			err = fmt.Errorf("channelId(%s) is invalid", item.Snippet.ChannelId)
			return err
		}
	}
	return nil
}
