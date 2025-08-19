package pusher

import (
	"github.com/pusher/pusher-http-go/v5"
)

var (
	Client *pusher.Client
)

func InitPusher(appId string, key string, secret string, cluster string, secure bool) {
	Client = &pusher.Client{
		AppID:   appId,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  secure,
	}
}
