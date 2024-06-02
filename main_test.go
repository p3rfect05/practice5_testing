package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

// TestAddTwoNumbers выполняет 1 задание
func TestAddTwoNumbers(t *testing.T) {
	float_tests := []struct {
		x        float64
		y        float64
		expected float64
	}{
		{x: 1.2, y: 5.7, expected: 6.9},
		{x: -1.2, y: 5.7, expected: 4.5},
		{x: 0.2, y: 0.1, expected: 0.3},
		{x: -0.3, y: 1.2, expected: 0.9},
		{x: 11.2, y: 52323.7, expected: 52334.9},
	}
	for _, tc := range float_tests {
		if math.Abs(AddTwoNumbers(tc.x, tc.y)-tc.expected) > math.Pow10(-6) {
			t.Fatalf("Sum of %f + %f is not equal to %f but should be", tc.x, tc.y, tc.expected)
		}
	}
	int_tests := []struct {
		x        int
		y        int
		expected int
	}{
		{x: 3, y: 6, expected: 9},
		{x: -2, y: 5, expected: 3},
		{x: 14, y: 19999, expected: 20013},
		{x: 51, y: 49, expected: 100},
		{x: 11, y: 52323, expected: 52334},
	}
	for _, tc := range int_tests {
		if tc.x+tc.y != tc.expected {
			t.Fatalf("Sum of %d + %d is not equal to %d but should be", tc.x, tc.y, tc.expected)
		}
	}
}

// TestThrowMistakeIfInputIsOdd выполняет 2 задание
func TestThrowMistakeIfInputIsOdd(t *testing.T) {
	error_test := 1
	success_test := 2
	if err := ThrowMistakeIfInputIsOdd(error_test); err == nil {
		t.Fatalf("Input %d should have thrown an error", error_test)
	}
	if err := ThrowMistakeIfInputIsOdd(success_test); err != nil {
		t.Fatalf("Input %d should not have thrown an error", success_test)
	}
}

// TestGetUserInfo выполняет 3 задание
func TestCountInts(t *testing.T) {
	tests := []struct {
		input    []int
		expected map[int]int
	}{
		{input: []int{1, 5, 7, 3, 1, 1, 2}, expected: map[int]int{1: 3, 5: 1, 3: 1, 7: 1, 2: 1}},
		{input: []int{-1, 3, 14, 5, 5, 5, 5}, expected: map[int]int{-1: 1, 3: 1, 14: 1, 5: 4}},
		{input: []int{}, expected: map[int]int{}},
		{input: []int{0}, expected: map[int]int{0: 1}},
	}

	for _, tc := range tests {
		res := CountInts(tc.input)
		if !reflect.DeepEqual(res, tc.expected) {
			t.Errorf("input: %#v, got %#v, expected %#v", tc.input, res, tc.expected)
		}
	}
}

// TestGetUserInfo выполняет 5 задание
func TestGetUserInfo(t *testing.T) {
	test_urls := []struct {
		input_url              string
		should_contain_message string
		expected_status_code   int
	}{
		{input_url: "/?user_id=2", should_contain_message: "2", expected_status_code: http.StatusOK},
		{input_url: "/?user_id=-2", should_contain_message: "non-existent", expected_status_code: http.StatusNotFound},
		{input_url: "/?u_id=2", should_contain_message: "no user_id", expected_status_code: http.StatusBadRequest},
		{input_url: "/?user_id=ab", should_contain_message: "invalid", expected_status_code: http.StatusBadRequest},
		{input_url: "/?user_id=", should_contain_message: "invalid", expected_status_code: http.StatusBadRequest},
	}
	for _, tc := range test_urls {
		req := httptest.NewRequest("GET", tc.input_url, nil)
		w := httptest.NewRecorder()
		GetUserInfo(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("could not read body when making request to: %s", tc.input_url)
		}
		if res.StatusCode != tc.expected_status_code {
			t.Fatalf("input_url=%s, got status %d, expected %d", tc.input_url, res.StatusCode, tc.expected_status_code)
		} else if !strings.Contains(string(data), tc.should_contain_message) {
			t.Fatalf("got response: %s, expected: %s", string(data), tc.should_contain_message)
		}
	}

}

type MockClock struct {
	CurrentTime time.Time
}

func (m *MockClock) Now() time.Time {
	return m.CurrentTime
}

// TestAcceptTicket выполняет 7 задание
func TestAcceptTicket(t *testing.T) {
	acceptedTime, expected_accepted := time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC), true
	refusedTime1, expected_refused1 := time.Date(2024, 1, 3, 6, 0, 0, 0, time.UTC), false
	refusedTime2, expected_refused2 := time.Date(2024, 1, 3, 19, 0, 0, 0, time.UTC), false

	res1 := !AcceptTicket(&MockClock{CurrentTime: acceptedTime})
	if res1 {
		t.Fatalf("time: %#v, got %v, expected %v", acceptedTime, res1, expected_accepted)
	}

	res2 := AcceptTicket(&MockClock{CurrentTime: refusedTime1})
	if res2 {
		t.Fatalf("time: %#v, got %v, expected %v", refusedTime1, res2, expected_refused1)
	}

	res3 := AcceptTicket(&MockClock{CurrentTime: refusedTime2})
	if res3 {
		t.Fatalf("time: %#v, got %v, expected %v", refusedTime2, res3, expected_refused2)
	}
}

// TestSpecificErrors выполняет 8 задание
func TestMergeChannels(t *testing.T) {
	sl1, sl2, sl3, sl4, sl5 := []int{1, 3, 9, 10}, []int{-1, 5, 6, 7, 8}, []int{}, []int{0}, []int{1}
	expected := slices.Concat(sl1, sl2, sl3, sl4, sl5)
	ch1, ch2, ch3, ch4, ch5 := CreateAndCloseChannel(sl1...), CreateAndCloseChannel(sl2...),
		CreateAndCloseChannel(sl3...), CreateAndCloseChannel(sl4...), CreateAndCloseChannel(sl5...)

	res := make([]int, 0, 11)
	for val := range MergeChannels(ch1, ch2, ch3, ch4, ch5) {
		res = append(res, val)
	}

	slices.SortStableFunc(res, func(a int, b int) int {
		return a - b
	})

	slices.SortStableFunc(expected, func(a int, b int) int {
		return a - b
	})
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("got %#v, expected %#v", res, expected)
	}

}

// TestSpecificErrors выполняет 9 задание
func TestSpecificErrors(t *testing.T) {
	err := GetFileOpenError()
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("got \"%v\" error, expected \"%v\"", err, fs.ErrNotExist)
	}
}

// TestPostLoginHandler выполняет 6 задание
func TestNewsStatistics(t *testing.T) {
	m := &MockNewsService{}

	m.On("GetTotalNewsByHours", 24).Return(10)
	res := NewsStatistics(24, m)
	if res != 10 {
		t.Errorf("got %d, expected %d", res, 10)
	}

	m.AssertCalled(t, "GetTotalNewsByHours", 24)
}

// TestPostLoginHandler выполняет 10 задание
func TestPostLoginHandler(t *testing.T) {
	type postData struct {
		key   string
		value string
	}
	tests := []struct {
		params                 []postData
		expected_status_code   int
		should_contain_message string
	}{
		{params: []postData{
			{key: "login", value: "alex"},
			{key: "password", value: "1234"},
		}, expected_status_code: http.StatusForbidden, should_contain_message: "invalid credentials"},
		{params: []postData{
			{key: "logi", value: "alex"}, // non-existent field
			{key: "password", value: "1234"},
		}, expected_status_code: http.StatusBadRequest, should_contain_message: "invalid body"},
		{params: []postData{
			{key: "login", value: "admin"}, // no password field
		}, expected_status_code: http.StatusBadRequest, should_contain_message: "invalid body"},
		{params: []postData{}, expected_status_code: http.StatusBadRequest, should_contain_message: "invalid body"}, // empty body
		{params: []postData{
			{key: "login", value: "admin"},
			{key: "password", value: "1234"},
		}, expected_status_code: http.StatusOK, should_contain_message: "successful"},
	}
	route := chi.NewRouter()
	route.Post("/", PostLoginHandler)
	srv := httptest.NewServer(route)
	defer srv.Close()
	for _, tc := range tests {
		values := make(map[string]string)
		for _, x := range tc.params {
			values[x.key] = x.value
		}
		byte_json, err := json.Marshal(values)
		if err != nil {
			t.Errorf("error while marshalling %v", values)
		}

		resp, err := srv.Client().Post(srv.URL+"/", "application/json", bytes.NewBuffer(byte_json))
		if err != nil {
			t.Fatalf("got error %v, when making request with %#v", err, tc.params)
		}
		var jsonResp jsonResponse
		err = json.NewDecoder(resp.Body).Decode(&jsonResp)
		defer resp.Body.Close()

		if err != nil {
			t.Fatalf("got error %v, when decoding body with %#v", err, tc.params)
		}
		if resp.StatusCode != tc.expected_status_code || !strings.Contains(jsonResp.Message, tc.should_contain_message) {
			t.Fatalf("code: got %d, expected %d\nmessage: got %s, expected %s when params %#v",
				resp.StatusCode, tc.expected_status_code, jsonResp.Message, tc.should_contain_message, values)
		}

	}

}
