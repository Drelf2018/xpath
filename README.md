# xpath

使用 tag 获取对应 xpath 元素

### 使用

```go
package xpath_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/xpath"
)

type College struct {
	Name string `xpath:"./td[1]//span/text()"`
	URL  string `xpath:"./td[2]//a/@href"`
}

func (c College) String() string {
	return fmt.Sprintf("%s(%s)", c.Name, c.URL)
}

type Colleges []College

func (Colleges) XPath() string {
	return "//*[@id='vsb_content']/div/table/tbody/tr[position() > 1]"
}

func TestColleges(t *testing.T) {
	colleges, err := xpath.LoadURLWith[Colleges]("https://gs.dlut.edu.cn/info/1173/14056.htm")
	if err != nil {
		t.Fatal(err)
	}
	for _, college := range colleges {
		t.Log(college)
	}
}
```