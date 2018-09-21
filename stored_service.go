package magicland

import (
	"fmt"

	"github.com/go-redis/redis"
)

const (
	localRedisURL        = "localhost:6379"
	defaultRedisPassword = ""
	defaultRedisDB       = 0
)

var (
	serviceFields = []string{"user", "branchname", "repositoryurl", "servicename"}

	redisClient *redis.Client
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     envOr("REDIS_URL", localRedisURL),
		Password: envOr("REDIS_PASSWORD", defaultRedisPassword),
		DB:       envOrI("REDIS_DB", defaultRedisDB),
	})
}

func saveService(gitConfig GitConfiguration) error {
	key := "svc:" + gitConfig.ServiceName
	err := redisClient.HMSet(key, gitConfig.AsRedisHashMap()).Err()
	if err != nil {
		return err
	}
	return nil
}

func selectServiceByName(serviceName string) (GitConfiguration, error) {
	gitConfiguration := GitConfiguration{}
	key := "svc:" + serviceName
	svcDef, err := redisClient.HMGet(key, serviceFields...).Result()
	if err != nil {
		fmt.Printf("Error getting %s: %v\n", key, err)
		return gitConfiguration, err
	}

	return gitConfiguration.FromRedisHashMap(svcDef), nil
}
