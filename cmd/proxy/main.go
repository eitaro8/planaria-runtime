package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 全 HTTP ヘッダーを出力して Envoy が付与しているトレースヘッダーを特定
	for k, v := range r.Header {
		log.Printf("[Header Log] %s: %s", k, v)
	}

	// Envoy が自動付与したトレース関連ヘッダーをログ出力
	log.Printf("[Header Check] traceparent: %s", r.Header.Get("traceparent"))
	log.Printf("[Header Check] x-request-id: %s", r.Header.Get("x-request-id"))

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for i := 1; i <= 10; i++ {
		select {
		case <-r.Context().Done():
			log.Println("[LLM] クライアント切断を検知し処理を停止しました")
			return
		case <-time.After(time.Second):
			if _, err := fmt.Fprintf(w, "data: {\"token\": \"chunk-%d\"}\n\n", i); err != nil {
				log.Printf("[LLM] failed to write event: %v", err)
				return
			}
			flusher.Flush()
		}
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", eventsHandler)
	mux.HandleFunc("/healthz", healthHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server starting on :8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
