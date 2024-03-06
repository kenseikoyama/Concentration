package server_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"kadai/server"
)


func Test_Handler(t *testing.T) {
	url := "http://localhost:8080"
	data := []byte(`{"username":"test","pass":"test"}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request in test server. req: %+v", r)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello in test server."))
	}))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ts.Config.Handler.ServeHTTP(w, r)
	ew, err := server.New("localhost:8080")
	if err != nil {
		t.Fatal("予期せぬエラー:", err)
	}

	if err = ew.Start(); err != nil {
		t.Fatal("予期せぬエラー:", err)
	}

	ew.HandleIndex(w, r)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})
	if res.StatusCode != http.StatusOK {
		t.Error("期待しないステータスコード:", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("予期せぬエラー:", err)
	}

	// 改行区切り
	expected := "Hello, World!\n"
	if string(body) != expected {
		t.Errorf("期待する内容: %s, 実際の内容: %s", expected, string(body))
	}

}
