package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"kadai/server"
)

func Test_Handler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ew, err := server.New(":8080")
	if err != nil {
		t.Fatal("予期せぬエラー:", err)
	}
	ew.HandleIndex(w, r)
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
	t.Log(string(body))
}
