package cache_test

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	cache "microup.ru/vcache/src"
)

func TestCache_Add(t *testing.T) {
	t.Parallel()

	testCache := cache.New(time.Second, time.Minute)
	keys := []string{}

	for index := 0; index < 1000; index++ {
		key := generateRandomKey(32)
		value, err := generateRandomValue()

		if err != nil {
			t.Fatalf("failed test, get err: %v", err)
		}

		keys = append(keys, key)
		testCache.Add(key, value)

		_, foundKey := testCache.Get(key)
		if !foundKey {
			t.Fatalf("Key %s not found in cache", key)
		}

		// Ensure that each key added to the cache is unique
		for _, existingKey := range keys[:index] {
			if existingKey == key {
				t.Fatalf("Key %s is not unique", key)
			}
		}
	}
}

func generateRandomKey(length int) string {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	key := make([]byte, length)
	_, err := rand.Read(key)

	if err != nil {
		return ""
	}

	for i, b := range key {
		key[i] = chars[b%byte(len(chars))]
	}

	return base64.RawURLEncoding.EncodeToString(key)
}

func generateRandomValue() (interface{}, error) {
	randomInt := make([]byte, 4)
	_, err := rand.Read(randomInt)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return int(randomInt[0])*256*256*256 + int(randomInt[1])*256*256 + int(randomInt[2])*256 + int(randomInt[3]), nil
}
