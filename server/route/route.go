package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ytrelay "github.com/mirror-media/yt-relay"
	"github.com/mirror-media/yt-relay/api"
	"github.com/mirror-media/yt-relay/cache"
	"github.com/mirror-media/yt-relay/config"
	"github.com/mirror-media/yt-relay/middleware"
	"github.com/mirror-media/yt-relay/relay"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/youtube/v3"
)

const (
	ErrorEmptyPart = "part cannot be empty"
	ErrorEmptyID   = "id cannot be empty"
)

func getCacheTTL(cacheConf config.Cache, api string) (ttl time.Duration, isDisabled bool) {

	isDisabled = cacheConf.DisabledAPIs[api]

	seconds, ok := cacheConf.OverwriteTTL[api]
	if ok {
		ttl = time.Duration(seconds) * time.Second
	} else {
		ttl = time.Duration(cacheConf.TTL) * time.Second
	}
	return ttl, isDisabled
}

func saveOKCache(isEnabled bool, cacheConf config.Cache, cacheProvider cache.Rediser, apiLogger *log.Entry, appName string, request http.Request, resp interface{}) {

	if cacheConf.IsEnabled {
		ttl, isCacheDisabledForAPI := getCacheTTL(cacheConf, request.RequestURI)
		if !isCacheDisabledForAPI {
			saveCache(cacheConf, cacheProvider, apiLogger, appName, request, http.StatusOK, resp, ttl)
		} else {
			apiLogger.Infof("cache is disabled for %s", request.URL.String())
		}
	}
}
func saveErrCache(isEnabled bool, cacheConf config.Cache, cacheProvider cache.Rediser, apiLogger *log.Entry, appName string, request http.Request, httpResponseCode uint, resp interface{}) {

	if cacheConf.IsEnabled {
		_, isCacheDisabledForAPI := getCacheTTL(cacheConf, request.RequestURI)
		if !isCacheDisabledForAPI {
			ttl := time.Duration(cacheConf.ErrorTTL) * time.Second
			saveCache(cacheConf, cacheProvider, apiLogger, appName, request, http.StatusOK, resp, ttl)
		} else {
			apiLogger.Infof("cache is disabled for %s", request.URL.String())
		}
	}
}

func saveCache(cacheConf config.Cache, cacheProvider cache.Rediser, apiLogger *log.Entry, appName string, request http.Request, respCode int, resp interface{}, ttl time.Duration) {
	s, err := json.Marshal(resp)
	if err != nil {
		apiLogger.Errorf("Cannot marshal resp for %s: %s", request.URL.String(), err)
		return
	}
	s, err = json.Marshal(cache.HTTP{
		StatusCode: respCode,
		Response:   s,
	})
	if err != nil {
		apiLogger.Errorf("Cannot marshal http resp cache for %s: %s", request.URL.String(), err)
		return
	}
	key, err := cache.GetCacheKey(appName, request.URL.String())
	if err != nil {
		apiLogger.Errorf("GetCacheKey for %s encounter error:%v", request.URL.String(), err)
	}
	err = cacheProvider.Set(request.Context(), key, string(s), ttl).Err()
	if err != nil {
		apiLogger.Errorf("setting cache encountered error for %s: %v ", request.URL.String(), err)
		return
	} else {
		apiLogger.Infof("cache for %s is set for ttl(%d) at %v", request.URL.String(), ttl, time.Now().UTC())
	}
}

// Set sets the routing for the gin engine
// TODO move whitelist to YouTube relay service
func Set(r *gin.Engine, appName string, relayService ytrelay.VideoRelay, whitelist ytrelay.APIWhitelist, cacheConf config.Cache, cacheProvider cache.Rediser) error {

	// health check api
	// As more resources and component are used, they should be checked in the api
	r.GET("/health", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	ytRouter := r.Group("/youtube/v3")

	if cacheConf.IsEnabled {
		ytRouter.Use(middleware.Cache(appName, cacheConf, cacheProvider))
	}

	// search videos. ChannelID is required
	ytRouter.GET("/search", func(c *gin.Context) {

		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		queries, err := parseQueries(c)
		if err != nil {
			apiLogger.Error(err)
			resp := api.ErrorResp{Error: err.Error()}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		// Check the mandatory parameters
		if queries.Part == "" {
			apiLogger.Error(ErrorEmptyPart)
			resp := api.ErrorResp{Error: ErrorEmptyPart}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		// Check whitelist
		if !whitelist.ValidateChannelID(queries.ChannelID) {
			err = fmt.Errorf("channelId(%s) is invalid", queries.ChannelID)
			apiLogger.Error(err)
			resp := api.ErrorResp{Error: err.Error()}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		resp, err := relayService.Search(queries)
		if err != nil {
			apiLogger.Error(err)
			resp := api.ErrorResp{Error: err.Error()}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			// FIXME internal error
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: err.Error()})
			return
		}
		saveOKCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, resp)
		c.JSON(http.StatusOK, resp)
	})

	// list video by video id
	// IDs of videos is required
	ytRouter.GET("/videos", func(c *gin.Context) {

		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		queries, err := parseQueries(c)
		if err != nil {
			apiLogger.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResp{Error: err.Error()})
			return
		}

		// Check the mandatory parameters
		if queries.Part == "" {
			apiLogger.Error(ErrorEmptyPart)
			resp := api.ErrorResp{Error: ErrorEmptyPart}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
		if queries.IDs == "" {
			apiLogger.Error(ErrorEmptyID)
			resp := api.ErrorResp{Error: ErrorEmptyID}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		resp, err := relayService.ListByVideoIDs(queries)
		if err != nil {
			apiLogger.Error(err)
			resp := api.ErrorResp{Error: err.Error()}
			// FIXME internal error
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		// verify channel id for YouTube
		_, isYouTube := relayService.(*relay.YouTubeServiceV3)
		if isYouTube {
			if err = validateYouTubeVideoListResponse(whitelist, resp); err != nil {
				err = errors.Wrap(err, "some video's channel id is invalid")
				apiLogger.Error(err)
				resp := api.ErrorResp{Error: err.Error()}
				saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
				c.AbortWithStatusJSON(http.StatusBadRequest, resp)
				return
			}
		}

		saveOKCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, resp)
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
			resp := api.ErrorResp{Error: err.Error()}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		// Check the mandatory parameters
		if queries.Part == "" {
			apiLogger.Error(ErrorEmptyPart)
			resp := api.ErrorResp{Error: ErrorEmptyPart}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		// Check whitelist
		if !whitelist.ValidatePlaylistIDs(queries.PlaylistID) {
			err = fmt.Errorf("playlistId(%s) is invalid", queries.PlaylistID)
			apiLogger.Error(err)
			resp := api.ErrorResp{Error: err.Error()}
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		resp, err := relayService.ListPlaylistVideos(queries)
		if err != nil {
			apiLogger.Error(err)
			resp := api.ErrorResp{Error: err.Error()}
			// FIXME internal error
			saveErrCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, http.StatusBadRequest, resp)
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		saveOKCache(cacheConf.IsEnabled, cacheConf, cacheProvider, apiLogger, appName, *c.Request, resp)
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
		if !whitelist.ValidateChannelID(item.Snippet.ChannelId) {
			err = fmt.Errorf("channelId(%s) is invalid", item.Snippet.ChannelId)
			return err
		}
	}
	return nil
}
