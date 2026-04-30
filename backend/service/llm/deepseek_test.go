package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newServer(handler http.HandlerFunc) (*Client, *httptest.Server) {
	s := httptest.NewServer(handler)
	c := &Client{
		BaseURL: s.URL,
		APIKey:  "sk-test",
		Model:   "deepseek-v4-flash",
		Timeout: 5 * time.Second,
		HTTP:    s.Client(),
	}
	return c, s
}

func TestGenerate_Success(t *testing.T) {
	c, s := newServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sk-test" {
			t.Errorf("missing or wrong Authorization header: %q", r.Header.Get("Authorization"))
		}
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req["model"] != "deepseek-v4-flash" {
			t.Errorf("model = %v, want deepseek-v4-flash", req["model"])
		}
		rf, _ := req["response_format"].(map[string]any)
		if rf == nil || rf["type"] != "json_object" {
			t.Errorf("response_format = %v, want {type: json_object}", req["response_format"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{
					"content": `{"summary":"All good.","note":"Continue plan."}`,
				}},
			},
		})
	})
	defer s.Close()

	out, err := c.Generate(context.Background(), "sys", "user")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if out.Summary != "All good." {
		t.Errorf("summary = %q", out.Summary)
	}
	if out.Note != "Continue plan." {
		t.Errorf("note = %q", out.Note)
	}
}

func TestGenerate_Truncates(t *testing.T) {
	long := strings.Repeat("x", 250)   // > 200 char summary
	longer := strings.Repeat("y", 600) // > 500 char note
	c, s := newServer(func(w http.ResponseWriter, r *http.Request) {
		body := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{
					"content": `{"summary":"` + long + `","note":"` + longer + `"}`,
				}},
			},
		}
		json.NewEncoder(w).Encode(body)
	})
	defer s.Close()

	out, err := c.Generate(context.Background(), "s", "u")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(out.Summary) != 200 {
		t.Errorf("summary len = %d, want 200", len(out.Summary))
	}
	if len(out.Note) != 500 {
		t.Errorf("note len = %d, want 500", len(out.Note))
	}
}

func TestGenerate_Failures(t *testing.T) {
	cases := []struct {
		name    string
		handler http.HandlerFunc
		substr  string
	}{
		{"500", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "boom", 500) }, "status 500"},
		{"429", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "rate", 429) }, "status 429"},
		{"empty content", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]string{"content": ""}}},
			})
		}, "empty"},
		{"bad json", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]string{"content": "not json"}}},
			})
		}, "parse"},
		{"missing summary", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]string{"content": `{"note":"x"}`}}},
			})
		}, "summary"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, s := newServer(tc.handler)
			defer s.Close()
			_, err := c.Generate(context.Background(), "s", "u")
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(strings.ToLower(err.Error()), tc.substr) {
				t.Errorf("err = %v, want contains %q", err, tc.substr)
			}
		})
	}
}
