package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
)

func main() {

}

type Number interface {
	float32 | float64 | int
}

func AddTwoNumbers[T Number](x, y T) T {
	return x + y
}

func ThrowMistakeIfInputIsOdd(x int) error {
	if x%2 == 1 {
		return errors.New("input is odd")
	}
	return nil
}

// CountInts считает количество переменных в мапе
func CountInts(input []int) map[int]int {
	res := make(map[int]int)
	for _, val := range input {
		res[val]++
	}
	return res
}

type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// GetRecordById мокает обращение к БД
func GetRecordById(id int) error {
	if id <= 0 || id > 5 {
		return fmt.Errorf("user with id %d does not exist", id)
	}
	return nil
}

// GetUserInfo возвращает json с id юзера либо с ошибкой
func GetUserInfo(w http.ResponseWriter, r *http.Request) {

	urlValues := r.URL.Query()
	if string_user_id, ok := urlValues["user_id"]; ok {
		user_id, err := strconv.Atoi(string_user_id[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid user_id"))
			return
		}
		if err := GetRecordById(user_id); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("non-existent user_id"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string_user_id[0]))

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no user_id"))
	}

}

type Clock interface {
	Now() time.Time
}

// AcceptTicket принимает тикет, если время находится в нужном промежутке
func AcceptTicket(clock Clock) bool {
	currentTime := clock.Now()
	if currentTime.Hour() >= 9 && currentTime.Hour() <= 17 {
		// submit ticket somewhere
		return true
	}
	return false
}

// CreateAndCloseChannel создаёт канал, заполняет его значением и закрывает от записи
func CreateAndCloseChannel[T any](values ...T) chan T {
	ch := make(chan T, len(values))
	for _, val := range values {
		ch <- val
	}
	close(ch)
	return ch
}

func MergeChannels[T any](chans ...chan T) chan T {
	mergedChan := make(chan T)
	wg := sync.WaitGroup{}
	wg.Add(len(chans))
	go func() {

		for _, ch := range chans {
			go func(ch chan T) {
				for val := range ch {
					mergedChan <- val
				}
				wg.Done()
			}(ch)
		}
		wg.Wait()
		close(mergedChan)
	}()

	return mergedChan
}

// GetFileOpenError открывает несуществующий файл, возвращая ошибку
func GetFileOpenError() error {
	if _, err := os.Open("lol.kek"); err != nil {
		return err
	}
	return nil
}

// NewsService общий интерфейс для работы с новостями, также с помощью него можно мокать внешние зависимости
type NewsService interface {
	GetTotalNewsByHours(hours int) int
}

type NewYorkTimesService struct {
}

// GetTotalNewsByHours должно обращаться к внешнему серверу и возвращать какие-либо статистические данные
func (news *NewYorkTimesService) GetTotalNewsByHours(hours int) int {
	// getting actual news by hours
	return 100
}

func NewsStatistics(hours int, news NewsService) int {
	return news.GetTotalNewsByHours(hours)

}

type MockNewsService struct {
	mock.Mock
}

// GetTotalNewsByHours мокает функцию NewYorkTimesService.GetTotalNewsByHours
func (m *MockNewsService) GetTotalNewsByHours(hours int) int {
	// getting actual news by hours
	args := m.Called(hours)
	return args.Int(0)
}

type jsonResponse struct {
	Message string `json:"msg"`
}

type LoginCredentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

// PostLoginHandler проверяет наличие логина и пароля, проверят их корректность
func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	defer r.Body.Close()
	if err != nil || creds.Password == "" || creds.Username == "" {

		resp, _ := json.Marshal(jsonResponse{Message: "invalid body values"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	if creds.Username == "admin" && creds.Password == "1234" {
		resp, _ := json.Marshal(jsonResponse{Message: "successful login!"})
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else {
		resp, _ := json.Marshal(jsonResponse{Message: "invalid credentials!"})
		w.WriteHeader(http.StatusForbidden)
		w.Write(resp)
	}

}
