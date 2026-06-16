package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, kv map[string]string) {
	t.Helper()
	for k, v := range kv {
		old, had := os.LookupEnv(k)
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("failed to set env %s: %v", k, err)
		}
		k, old, had := k, old, had
		t.Cleanup(func() {
			if had {
				os.Setenv(k, old)
			} else {
				os.Unsetenv(k)
			}
		})
	}
}

func TestLoadFromEnv(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_NAME":               "mm-yt-relay.test",
		"API_KEY":                "test-api-key",
		"WHITELIST_CHANNEL_IDS":  "channelA,channelB",
		"WHITELIST_PLAYLIST_IDS": "playlistA,playlistB",
		"CACHE_ENABLED":          "true",
		"CACHE_TTL":              "1800",
		"CACHE_ERROR_TTL":        "60",
		"CACHE_OVERWRITE_TTL":    "/youtube/v3/playlistItems:300",
		"REDIS_TYPE":             "single",
		"REDIS_ADDRESSES":        "10.0.0.1:6379",
	})

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load from env failed: %v", err)
	}

	if cfg.ApiKey != "test-api-key" {
		t.Errorf("ApiKey = %q, want %q", cfg.ApiKey, "test-api-key")
	}
	if cfg.AppName != "mm-yt-relay.test" {
		t.Errorf("AppName = %q, want %q", cfg.AppName, "mm-yt-relay.test")
	}
	if !cfg.Whitelists.ChannelIDs["channelA"] || !cfg.Whitelists.ChannelIDs["channelB"] {
		t.Errorf("ChannelIDs = %v, want channelA and channelB enabled", cfg.Whitelists.ChannelIDs)
	}
	if !cfg.Whitelists.PlaylistIDs["playlistA"] || !cfg.Whitelists.PlaylistIDs["playlistB"] {
		t.Errorf("PlaylistIDs = %v, want playlistA and playlistB enabled", cfg.Whitelists.PlaylistIDs)
	}
	if !cfg.Cache.IsEnabled || cfg.Cache.TTL != 1800 || cfg.Cache.ErrorTTL != 60 {
		t.Errorf("Cache = %+v, want enabled with ttl 1800 and errorTtl 60", cfg.Cache)
	}
	if got := cfg.Cache.OverwriteTTL["/youtube/v3/playlistItems"]; got != 300 {
		t.Errorf("OverwriteTTL = %d, want 300", got)
	}
	if cfg.Redis == nil || cfg.Redis.Type != Single {
		t.Fatalf("Redis = %+v, want single instance", cfg.Redis)
	}
	if addr := cfg.Redis.SingleInstance.Instance; addr.Addr != "10.0.0.1" || addr.Port != 6379 {
		t.Errorf("Redis address = %+v, want 10.0.0.1:6379", addr)
	}
}

func TestLoadFromEnvMissingPlaylistWhitelist(t *testing.T) {
	setEnv(t, map[string]string{
		"APP_NAME":               "mm-yt-relay.test",
		"API_KEY":                "test-api-key",
		"WHITELIST_CHANNEL_IDS":  "channelA",
		"WHITELIST_PLAYLIST_IDS": "",
	})

	if _, err := Load(""); err == nil {
		t.Fatal("Load should fail when playlist whitelist is empty")
	}
}

func TestLoadFromFile(t *testing.T) {
	yml := `
appName: mm-yt-relay.file
apiKey: file-api-key
whitelists:
  channelIDs:
    channelA: true
  playlistIDs:
    playlistA: true
`
	f, err := os.CreateTemp(t.TempDir(), "config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(yml); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load from file failed: %v", err)
	}
	if cfg.ApiKey != "file-api-key" {
		t.Errorf("ApiKey = %q, want %q", cfg.ApiKey, "file-api-key")
	}
	if !cfg.Whitelists.PlaylistIDs["playlistA"] {
		t.Errorf("PlaylistIDs = %v, want playlistA enabled", cfg.Whitelists.PlaylistIDs)
	}
}
