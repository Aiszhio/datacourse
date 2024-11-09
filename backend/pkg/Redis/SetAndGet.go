package Redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func SetKey(client *redis.Client, key, value string) error {
	ctx := context.Background()
	return client.Set(ctx, key, value, 0).Err()
}

func GetKey(client *redis.Client, key string) (string, error) {
	ctx := context.Background()
	return client.Get(ctx, key).Result()
}

func SetMultipleKey(client *redis.Client, key string, values interface{}) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("error serializing data to JSON: %w", err)
	}

	if err = client.Set(ctx, key, jsonData, 5*time.Minute).Err(); err != nil {
		fmt.Println("Error saving to Redis:", err)
		return err
	}

	fmt.Println("Data saved to Redis successfully:", string(jsonData))

	savedData, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error retrieving saved data from Redis:", err)
		return err
	}

	fmt.Println("Retrieved data from Redis:", savedData)
	return nil
}

func GetMultipleKey(client *redis.Client, key string) (map[string]interface{}, error) {
	ctx := context.Background()

	jsonData, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error retrieving data from Redis:", err)
	}

	fmt.Println("Retrieved data from Redis:", jsonData)

	var result map[string]interface{}
	if err = json.Unmarshal([]byte(jsonData), &result); err != nil {
		fmt.Println("Error unmarshalling data from Redis:", err)
	}

	return result, nil
}
