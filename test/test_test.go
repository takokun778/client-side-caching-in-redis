package test_test

import (
	"io"
	"net/http"
	"testing"
)

func TestTest(t *testing.T) {
	t.Parallel()

	t.Run("クライアントサイドキャッシュが有効であることを確認する", func(t *testing.T) {
		t.Skip()

		client := http.DefaultClient

		{
			url := "http://localhost:8081/set?key=key&val=foo"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.Do(req); err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してキー(key)と値(foo)を設定")
		}

		{
			url := "http://localhost:8081/get?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8081/get?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8081/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8081/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8081/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		defer func() {
			url := "http://localhost:8081/del?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.Do(req); err != nil {
				t.Fatal(err)
			}
		}()
	})

	t.Run("クライアントキャッシュを有効にしてからAppAで値を削除してAppBで値を取得する", func(t *testing.T) {
		t.Skip()

		client := http.DefaultClient

		{
			url := "http://localhost:8081/set?key=key&val=foo"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.Do(req); err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してキー(key)と値(foo)を設定")
		}

		{
			url := "http://localhost:8081/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8082/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppBに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8081/del?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.Do(req); err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してキー(key)を指定して値を削除")
		}

		{
			url := "http://localhost:8082/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppBに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}
	})

	t.Run("クライアントサイドキャッシュを有効にしてからAppAで値を更新してAppBで値を取得する", func(t *testing.T) {
		t.Skip()

		client := http.DefaultClient

		{
			url := "http://localhost:8081/set?key=key&val=foo"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.Do(req); err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してキー(key)と値(foo)を設定")
		}

		{
			url := "http://localhost:8081/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8082/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppBに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8082/set?key=key&val=bar"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.Do(req); err != nil {
				t.Fatal(err)
			}

			t.Log("AppBに対してキー(key)と値(bar)を設定")
		}

		{
			url := "http://localhost:8081/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppAに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		{
			url := "http://localhost:8082/get/cache?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("AppBに対してクライアントサイドキャッシュ利用のキー(key)を指定して値を取得")
			t.Logf("val: %s", string(body))
		}

		defer func() {
			url := "http://localhost:8081/del?key=key"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.Do(req); err != nil {
				t.Fatal(err)
			}
		}()
	})
}
