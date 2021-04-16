package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/mirror-media/yt-relay/api"
	"github.com/mirror-media/yt-relay/cache"
	"github.com/mirror-media/yt-relay/config"
	log "github.com/sirupsen/logrus"
)

func Cache(namespace string, cacheConf config.Cache, cacheProvider cache.Rediser) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL

		// check blacklist
		if _, ok := cacheConf.DisabledAPIs[url.Path]; ok {
			log.Infof("cache is disabled for %s", url.Path)
			c.Next()
			return
		}
		// read cache
		uri := c.Request.RequestURI
		key, err := cache.GetCacheKey(namespace, uri)
		if err != nil {
			err = errors.Wrap(err, "Fail to get cache key in cache middleware")
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResp{Error: err.Error()})
			return
		}
		result, err := cacheProvider.Get(c.Request.Context(), key).Result()
		if err != nil {
			c.Next()
			return
		}

		log.Infof("respond with cache for %s", uri)
		c.AbortWithStatusJSON(http.StatusOK, json.RawMessage(result))
	}
}