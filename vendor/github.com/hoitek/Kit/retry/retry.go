package retry

import (
	"log"
	"time"
)

func Get[T interface{}](fn func() (T, error), backOffInterval, limit int) (T, error) {
	result, err := fn()
	count := 0
	for {
		if count >= limit {
			if err != nil {
				return T(result), err
			}
			break
		}
		if err == nil {
			break
		}
		log.Println("Retry to connect...")
		time.Sleep(time.Duration(backOffInterval) * time.Second)
		result, err = fn()
		count++
		continue
	}
	return result, nil
}
