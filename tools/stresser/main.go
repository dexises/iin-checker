package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Person struct {
	IIN   string `json:"iin"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	workers := getEnvInt("WORKERS", 10)
	total := getEnvInt("TOTAL_REQUESTS", 1000)
	serviceURL := getEnv("SERVICE_URL", "http://app:8080")

	jobs := make(chan int)
	var wg sync.WaitGroup

	// стартуем воркеры
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := &http.Client{Timeout: 5 * time.Second}
			for range jobs {
				// иногда генерим неправильный IIN
				var iin string
				if rand.Intn(5) == 0 {
					iin = generateInvalidIIN()
				} else {
					iin = generateValidIIN()
				}

				person := Person{
					IIN:   iin,
					Name:  fmt.Sprintf("User-%s", iin),
					Phone: fmt.Sprintf("+7%010d", rand.Int63n(1e10)),
				}
				body, _ := json.Marshal(person)
				resp, err := client.Post(serviceURL+"/people/info", "application/json", bytes.NewReader(body))
				log.Println(person.IIN + resp.Status)
				if err != nil {
					fmt.Printf("worker %d: error %v\n", id, err)
					continue
				}
				resp.Body.Close()
			}
		}(i)
	}

	// шлём задания
	for j := 0; j < total; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	fmt.Println("Stress test completed")
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func generateValidIIN() string {
	// генерим 11 цифр и считаем контрольную
	digits := make([]int, 11)
	for i := range digits {
		digits[i] = rand.Intn(10)
	}
	c := calcControlDigit(digits)
	return fmt.Sprintf("%s%d", intsToString(digits), c)
}

func generateInvalidIIN() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = byte('0' + rand.Intn(10))
	}
	return string(b)
}

func calcControlDigit(digits []int) int {
	w1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	w2 := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}
	sum := 0
	for i := 0; i < 11; i++ {
		sum += digits[i] * w1[i]
	}
	mod := sum % 11
	if mod == 10 {
		sum = 0
		for i := 0; i < 11; i++ {
			sum += digits[i] * w2[i]
		}
		mod = sum % 11
		if mod == 10 {
			return rand.Intn(10)
		}
	}
	return mod
}

func intsToString(d []int) string {
	s := ""
	for _, n := range d {
		s += strconv.Itoa(n)
	}
	return s
}
