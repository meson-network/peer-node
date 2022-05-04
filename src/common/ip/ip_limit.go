package ip

// import (
// 	"context"
// 	"time"

// 	"github.com/meson-network/peer-node/plugin/redis_plugin"
// )

// const redis_ip_limit_prefix = "ip_limit"

// func LimitAction(ip string, action string, duration time.Duration) {
// 	key := redis_plugin.GetInstance().GenKey(redis_ip_limit_prefix, ip, action)
// 	redis_plugin.GetInstance().Set(context.Background(), key, 1, duration)
// }

// func HasLimitedAction(ip string, action string) bool {
// 	key := redis_plugin.GetInstance().GenKey(redis_ip_limit_prefix, ip, action)
// 	_, err := redis_plugin.GetInstance().Get(context.Background(), key).Result()
// 	if err == redis.Nil {
// 		return false
// 	} else {
// 		return true
// 	}
// }
