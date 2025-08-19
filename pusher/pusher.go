package pusher

import (
	"github.com/pusher/pusher-http-go/v5"
)

type PusherService struct {
	Client *pusher.Client
}

func NewPusherService(appId string, key string, secret string, cluster string, secure bool) *PusherService {
	client := &pusher.Client{
		AppID:   appId,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  secure,
	}

	return &PusherService{Client: client}
}
