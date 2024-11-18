package Redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func SetKey(client *redis.Client, key, value string) error {
	ctx := context.Background()
	err := client.Set(ctx, key, value, 30*time.Minute).Err()
	if err != nil {
		log.Printf("Error setting key '%s' in Redis: %v", key, err)
	} else {
		log.Printf("Successfully set key '%s' in Redis with value '%s'", key, value)
	}
	return err
}

func GetKey(client *redis.Client, key string) (string, error) {
	ctx := context.Background()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("Key '%s' not found in Redis", key)
		} else {
			log.Printf("Error getting key '%s' from Redis: %v", key, err)
		}
		return "", err
	}
	log.Printf("Successfully retrieved key '%s' from Redis with value '%s'", key, value)
	return value, nil
}

func SetMultipleKey(client *redis.Client, key string, values interface{}) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("error serializing data to JSON: %w", err)
	}

	if err = client.Set(ctx, key, jsonData, 30*time.Minute).Err(); err != nil {
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
