package models

import (
	"context"
	"ginchat/utils"
	"time"
)

func setUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	utils.Red.Set(ctx, key, val, timeTTL)
}
