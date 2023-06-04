package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/takokun778/client-side-caching-in-redis/internal/redis"
)

func main() {
	rds := redis.New(os.Getenv("REDIS_URL"))

	defer rds.Close()

	hdl := &Handler{
		rds: rds,
	}

	http.HandleFunc("/set", hdl.Set)

	http.HandleFunc("/get", hdl.Get)

	http.HandleFunc("/del", hdl.Del)

	http.HandleFunc("/get/cache", hdl.GetCache)

	http.ListenAndServe(":8080", nil)
}

type Handler struct {
	rds *redis.Client
}

func (hdl *Handler) Set(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := r.URL.Query().Get("key")

	val := r.URL.Query().Get("val")

	log.Printf("key: %s, val: %s", key, val)

	cmd := hdl.rds.B().Set().Key(key).Value(val).Build()

	if err := hdl.rds.Do(ctx, cmd).Error(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (hdl *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := r.URL.Query().Get("key")

	cmd := hdl.rds.B().Get().Key(key).Build()

	val, err := hdl.rds.Do(ctx, cmd).ToString()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write([]byte(val))
}

func (hdl *Handler) GetCache(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := r.URL.Query().Get("key")

	cmd := hdl.rds.B().Get().Key(key).Cache()

	val, err := hdl.rds.DoCache(ctx, cmd, time.Hour).ToString()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write([]byte(val))
}

func (hdl *Handler) Del(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := r.URL.Query().Get("key")

	cmd := hdl.rds.B().Del().Key(key).Build()

	if err := hdl.rds.Do(ctx, cmd).Error(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
