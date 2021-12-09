# yt-relay

## Responsibility

`yt-relay` is an overlay above the YouTube APIs. It implemens cache, whitelist, scope the parameters of APIs to satisfy bussiness needs coincisely.

## API

`yt-relay` provides three APIs coresponding to the YouTube APIs:

- /youtube/v3/search
- /youtube/v3/videos
- /youtube/v3/playlistItems

## Whitelist

To prevent malicious usage, `yt-relay` embraces the whitelist mechanism.

Whitlists are defines in [Config](https://github.com/mirror-media/yt-relay/blob/f66310cdae732dcb2cc47a419755a9314fff9cf4/config/config.go#L25). They key could be channel id or playlist id. If the value is `true`, then it mean such element in the whitelist is activated.

## Detailed Document

Please visit [paper](https://paper.dropbox.com/doc/YouTube-Relay-API--BXs8b5XF~jHmLhKBbB7a8y_LAg-kTE4r0dLXThlRF2hVN5Cb)
