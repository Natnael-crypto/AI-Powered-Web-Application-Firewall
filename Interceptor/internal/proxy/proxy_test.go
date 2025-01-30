package proxy_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"interceptor/internal/waf"
)

type stubWAF struct {

}

func (s *stubWAF) EvaluateRules(r *http.Request)(bool, int, string, string, int){
	return false, 0, "", "", 0
}

func TestProxy(t *testing.T) {
	const (
		helloWorldMsg = "Hello, World!"
	)

	if _, err := waf.InitializeRuleEngine(""); err != nil {
		t.Fatalf("Failed to initialize WAF: %v", err)
	}

	unsafeTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Transport: unsafeTransport,
	}

	proxyServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, helloWorldMsg)
	}))

	resp, err := client.Get(proxyServer.URL)
	if err != nil {
		t.Fatalf("Error while sending request to the proxy server, %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error while reading response body, %v", err)
	}
	got := string(body)

	if got != helloWorldMsg {
		t.Errorf("got %q, want %q", got, helloWorldMsg)
	}
}
