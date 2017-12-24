## Multi_search

一回の検索で複数のサイトを検索し、結果をiFrameにて画面に表示します。
検索先がiFrameを弾いていても、本プログラムを経由すると表示することが可能です。

* * *

## 使い方

html主体で動作するので、multi_searchのtemplateの中ににある

index.htmlとsearch.htmlをmain.goと同じ場所に置いて下さい。

```
go get github.com/flowtumn/multi_search
```

```main.go
package main

import (
	"os"

	"github.com/flowtumn/multi_search"
)

func main() {
	pwd, _ := os.Getwd()
	ser, _ := muls.CreateSearchProxyServer(pwd, 20202)
	ser.Listen("localhost", 20202)
}
```

```
go run main.go
```

## 動作図

![preview](https://github.com/flowtumn/multi_search/blob/doc/image.gif "preview")