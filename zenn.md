目次
1. はじめに
1. そもそもRedisのクライアントサイドキャッシュとは??
1. rueidisを使った動作確認
1. おわりに

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

クライアントサイドキャッシュはハイパフォーマンスなサービスを実現するために使用する技術です。image.png
Redisのキャッシュを利用するのではなく、クライアント（=アプリケーション）のキャッシュを利用します。

クライアントサイドキャッシュを利用することによるメリットは以下の2点です。

- 非常に小さなレイテンシーでデータの利用
- Redisへの負荷を軽減

しかしクライアントサイドキャッシュを利用することには問題があります。
それはクライアントサイドキャッシュをどのように無効化するかということです。
アプリケーションによってはTTLを設定することで問題を解決できますが、

※ TTL...

# rueidisを使った動作確認

## DoCache() メソッド

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

# 

## ほかのサーバーのローカルキャッシュって消せるの??

## 値を変更した場合ってほかのサーバーに伝えられるの??

- ローカルキャッシュは万歳
- キャッシュを消したいときだってある
- TTLで制御するのはよろしくなさそう
- タイミングによってはへんな挙動になる

- aで保存
- bで取得
- aで1秒後に削除
- bで1秒後にキャッシュから取得

# おわりに
