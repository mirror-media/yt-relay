# yt-relay
   * For YouTube relay

## Configuration

The server can be configured either by a YAML file or by environment variables.

```sh
# with a config file (see configs/config.example.yml)
./yt-relay serve -config ./configs/config.yml -address 0.0.0.0 -port 8080

# with environment variables (no -config flag)
./yt-relay serve -address 0.0.0.0 -port 8080
```

### Environment variables

| Variable | Required | Description |
| --- | --- | --- |
| `APP_NAME` | yes | App name, used as the cache namespace. Alphanumeric, dot, and dash only |
| `API_KEY` | yes | YouTube Data API key |
| `WHITELIST_CHANNEL_IDS` | yes | Comma-separated channel IDs allowed in search/videos APIs |
| `WHITELIST_PLAYLIST_IDS` | yes | Comma-separated playlist IDs allowed in playlistItems API |
| `CACHE_ENABLED` | no | `true` to enable the redis cache (default `false`) |
| `CACHE_TTL` | if cache enabled | Default cache TTL in seconds |
| `CACHE_ERROR_TTL` | if cache enabled | TTL in seconds for cached error responses |
| `CACHE_DISABLED_APIS` | no | Comma-separated API paths excluded from cache, e.g. `/youtube/v3/search` |
| `CACHE_OVERWRITE_TTL` | no | Per-API TTL overrides, e.g. `/youtube/v3/playlistItems:300` |
| `REDIS_TYPE` | if cache enabled | `single`, `cluster`, `sentinel`, or `replica` |
| `REDIS_ADDRESSES` | if redis set | Comma-separated `host:port` list (writers for `replica` type) |
| `REDIS_READER_ADDRESSES` | only for `replica` | Comma-separated `host:port` list of readers |
| `REDIS_PASSWORD` | no | Redis password |

## Deployment (Cloud Build + Cloud Run)

`cloudbuild.yaml` builds the image, pushes it to `gcr.io/$PROJECT_ID/yt-relay`,
and deploys it to the Cloud Run service matching the branch being built.

Create a single Cloud Build trigger watching `^(dev|stag|master)$`;
the deploy step picks the target service from `$BRANCH_NAME`:

| Branch | Cloud Run service |
| --- | --- |
| `dev` | `yt-relay-dev` |
| `stag` | `yt-relay-staging` |
| `master` | `yt-relay-prod` |

Runtime settings (environment variables above, VPC connector for redis, etc.)
are managed directly on the Cloud Run services and are preserved across deploys.
