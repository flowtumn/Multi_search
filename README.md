## Multi_search

���̌����ŕ����̃T�C�g���������A���ʂ�iFrame�ɂĉ�ʂɕ\�����܂��B
�����悪iFrame��e���Ă��Ă��A�{�v���O�������o�R����ƕ\�����邱�Ƃ��\�ł��B

* * *

## �g����

html��̂œ��삷��̂ŁAmulti_search��template�̒��ɂɂ���

index.html��search.html��main.go�Ɠ����ꏊ�ɒu���ĉ������B

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

## ����}

![preview](https://github.com/flowtumn/multi_search/blob/doc/image.gif "preview")