# はじめに

以前「KVSとしてのRedisに入門しgo-redisとrueidisから触ってみる」という記事を投稿したところ

https://zenn.dev/takokun/articles/a3bdeee4f570f9

https://github.com/redis/rueidis

rueidisの作者さまからコメントをいただきました！

> Hi takokun, thank you for trying rueidis and sharing a detailed comparison with go-redis with us!
> 
> I am also encouraging you to try the client-side caching as well which is a killer feature since Redis 6. Hoping that you can get huge benefit from it.

前回の記事ももともとrueidisの`supports client side caching.`という一文がきっかけで記事を書いてみましたが今回の記事でいよいよクライアントサイドキャッシュについて触れてみたいと思います。

# Client-side caching in Redis

そもそもRedisのクライアントサイドキャッシュとはなんでしょうか？

Redisの公式ドキュメントに説明が記載されていました。

https://redis.io/docs/manual/client-side-caching/

[Bard](https://bard.google.com/)の力を借りつつ説明します。

クライアントサイドキャッシュはハイパフォーマンスなサービスを実現するために使用する技術です。
Redisのキャッシュを利用するのではなく、クライアントサイド（= アプリケーション）のキャッシュを利用します。

クライアントサイドキャッシュを利用することによるメリットは以下の2点があります。

- 非常に小さなレイテンシーでデータの利用
- Redisへの負荷を軽減

しかしクライアントサイドキャッシュを利用することにはキャッシュをどのように無効化するかという問題があります。（更新された場合に最新の値を取得したい、削除された場合に値が存在しないように振るまいたい）
アプリケーションによってはクライアントサイドにてTTLを設定することで問題を解決することが可能です。

※ TTL ... Time To Live の略で、有効期限を意味します。

しかし、TTLを設定することで問題が解決するとは限りません。
キャッシュの有効期限が切れる前に更新された場合に最新の値を取得できません。
そこでRedisではPub/Subを利用してクライアントに無効化メッセージを送信することができます。クライアントは無効化メッセージを受信することでキャッシュを無効化し不整合な値の取得を防ぐことができます。

# rueidisを使った動作確認

## クライアントサイドキャッシュを使うための DoCache() メソッド

```go
// DoCache is similar to Do, but it uses opt-in client side caching and requires a client side TTL.
// The explicit client side TTL specifies the maximum TTL on the client side.
// If the key's TTL on the server is smaller than the client side TTL, the client side TTL will be capped.
//  client.Do(ctx, client.B().Get().Key("k").Cache(), time.Minute).ToString()
// The above example will send the following command to redis if cache miss:
//  CLIENT CACHING YES
//  PTTL k
//  GET k
// The in-memory cache size is configured by ClientOption.CacheSizeEachConn.
// The cmd parameter is recycled after passing into DoCache() and should not be reused.
DoCache(ctx context.Context, cmd Cacheable, ttl time.Duration) (resp RedisResult)
```

DoCacheメソッドは、Doメソッドに似ていますが、クライアント側のキャッシュを使用します。サーバー上のキーの TTL がクライアント側の TTL より小さい場合、クライアント側の TTL には上限が設定されます。

DoCacheメソッドは、Redisコマンドをキャッシュします。
キャッシュが存在する場合、キャッシュから結果を返します。
キャッシュが存在しない場合、Redisコマンドを実行し、結果をキャッシュします。
DoCacheメソッドは、Redisコマンドの結果をキャッシュすることで、Redisへのリクエストを減らすことができます。これにより、アプリケーションのパフォーマンスが向上します。

## サンプルアプリの構築

以下のエンドポイントを持つサーバーを構築します。

- `GET /get?key=xxx` ... キーに紐づく値を取得します。
- `GET /get/cache?key=xxx` ... キーに紐づく値を取得します。クライアントサイドキャッシュを有効にします。
- `GET /set?key=xxx&val=yyy` ... キーに値を設定します。
- `GET /del?key=xxx` ... キーと値を削除します。

```go
type Client struct {
	rueidis.Client
}

func New(url string) *Client {
	ctx := context.Background()

	option := rueidis.ClientOption{
		InitAddress: []string{url},
	}

	cli, err := rueidis.NewClient(option)
	if err != nil {
		panic(err)
	}

	if err := cli.Do(ctx, cli.B().Ping().Build()).Error(); err != nil {
		panic(err)
	}

	return &Client{
		Client: cli,
	}
}
```

```go
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
```

`compose.yaml`を以下のように設定しサンプルアプリ×2 + Redisのコンテナを起動します。

AppA ... `localhost:8081`
AppB ... `localhost:8082`

```yaml
services:
  redis:
    container_name: redis
    image: redis:7.0.11-alpine
    ports:
      - 6379:6379
    restart: always
  app-a:
    container_name: app-a
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 8081:8080
    restart: always
    environment:
      REDIS_URL: redis:6379
    volumes:
      - ../:/app
  app-b:
    container_name: app-b
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 8082:8080
    restart: always
    environment:
      REDIS_URL: redis:6379
    volumes:
      - ../:/app
```

## クライアントサイドキャッシュの動作確認

## 値更新の動作確認

## 値削除の動作確認

# おわりに

Redisのクライアントサイドキャッシュの仕様を確認し、rueidisを使って実際に動作確認を行いました。
今回は説明をしませんでしたがクライアントサイドキャッシュを有効にした場合はRedisとの接続モードには2種類あります。
