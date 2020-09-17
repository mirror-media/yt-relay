package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	ytrelay "github.com/mirror-media/yt-relay"
	log "github.com/sirupsen/logrus"
)

const (
	ErrorEmptyPart = "part cannot be empty"
	ErrorEmptyID   = "id cannot be empty"
)

var json jsoniter.API

func init() {
	json = jsoniter.ConfigFastest
}

func Set(r *gin.Engine, relay ytrelay.VideoRelay, whitelist ytrelay.APIWhitelist) error {

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
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s", ErrorEmptyPart))
			return
		}

		// Check whitelist
		if !whitelist.ValidateChannelID(queries) {
			apiLogger.Error(fmt.Sprintf("channelId(%s) is invalid", queries.ChannelID))
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		resp, err := relay.Search(queries)
		if err != nil {
			apiLogger.Error(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
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
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s", ErrorEmptyPart))
			return
		}
		if queries.IDs == "" {
			apiLogger.Error(ErrorEmptyID)
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s", ErrorEmptyID))
			return
		}

		resp, err := relay.ListByVideoIDs(queries)
		if err != nil {
			apiLogger.Error(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// TODO verify channel id after implement relay service
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
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s", ErrorEmptyPart))
			return
		}

		// Check whitelist
		if !whitelist.ValidatePlaylistIDs(queries) {
			apiLogger.Error(fmt.Sprintf("channelId(%s) is invalid", queries.ChannelID))
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		resp, err := relay.ListByVideoIDs(queries)
		if err != nil {
			apiLogger.Error(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
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
