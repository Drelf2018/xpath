# xpath

使用 tag 获取对应 xpath 元素。

### 使用

```go
package xpath_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/xpath"
)

type Attachment struct {
	CollegeName string
	Name        string `xpath:"./a//text()"`
	Url         string `xpath:"./a/@href"`
}

type College struct {
	Name        string       `xpath:"./td[1]//span/text()" gorm:"primaryKey"`
	Url         string       `xpath:"./td[2]//a/@href"`
	Attachments []Attachment `xpath:"//td//span[a] | //form//li[a] | //ul[@style='list-style-type:none;']//li[a] | //ul[@class='attach']//li[a]"`
}

func (c College) String() string {
	return fmt.Sprintf("%s(%s, %d attachments)", c.Name, c.Url, len(c.Attachments))
}

type Page struct {
	Colleges []College `xpath:"//*[@id='vsb_content']/div/table/tbody/tr[position() > 1]"`
}

func TestMap(t *testing.T) {
	p := xpath.MustLoadURL[Page]("https://gs.dlut.edu.cn/info/1173/14056.htm")
	for _, c := range p.Colleges {
		if c.Url == "" {
			continue
		}
		err := xpath.LoadURL(c.Url, &c)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(c)
	}
}
```

#### 控制台

```
（主校区）数学科学学院(https://math.dlut.edu.cn/info/1083/17868.htm, 6 attachments)
（主校区）物理学院(https://physics.dlut.edu.cn/info/1055/9954.htm, 6 attachments)
（主校区）机械工程学院(http://me.dlut.edu.cn/info/1064/7736.htm, 7 attachments)
（主校区）材料科学与工程学院(https://mse.dlut.edu.cn/info/1037/4581.htm, 6 attachments)
（主校区）建设工程学院(https://sche.dlut.edu.cn/info/1225/26840.htm, 7 attachments)
（主校区）能源与动力学院(http://power.dlut.edu.cn/info/4659/22678.htm, 8 attachments)
（主校区）经济管理学院(https://sem.dlut.edu.cn/info/1006/10570.htm, 5 attachments)
（主校区）MBA/EMBA教育中心(https://mba.dlut.edu.cn/info/1021/11212.htm, 6 attachments)
（主校区）外国语学院(https://fld.dlut.edu.cn/info/1094/9128.htm, 6 attachments)
体育与健康学院(https://kh.dlut.edu.cn/info/1124/2928.htm, 7 attachments)
（主校区）建筑与艺术学院(https://aaschool.dlut.edu.cn/info/1029/4918.htm, 5 attachments)
（主校区）马克思主义学院(http://marx.dlut.edu.cn/info/1019/20871.htm, 5 attachments)
（主校区）光电工程与仪器科学学院(https://oeis.dlut.edu.cn/info/1145/5752.htm, 2 attachments)
（主校区）国际教育学院(http://sie.dlut.edu.cn/info/1020/22071.htm, 6 attachments)
（主校区）化工学院(http://chemeng.dlut.edu.cn/info/1118/15683.htm, 6 attachments)
（主校区）环境学院(http://est.dlut.edu.cn/info/1087/5115.htm, 6 attachments)
（主校区）生物工程学院(https://biotech.dlut.edu.cn/info/1343/6362.htm, 5 attachments)
（主校区）大连理工大学白俄罗斯国立大学联合学院(https://dbji.dlut.edu.cn/info/1005/5899.htm, 5 attachments)
（主校区）医学部(http://med.dlut.edu.cn/info/1003/3817.htm, 7 attachments)
（主校区）化学学院(https://zdysc.dlut.edu.cn/info/1010/4649.htm, 7 attachments)
人文学院(https://fhss.dlut.edu.cn/info/1411/32741.htm, 5 attachments)
（主校区）公共管理学院(https://spap.dlut.edu.cn/info/1124/4755.htm, 4 attachments)
（主校区）MPA教育中心(http://mpa.dlut.edu.cn/info/1024/3891.htm, 4 attachments)
高等教育研究院(https://gdjyyjzx.dlut.edu.cn/info/1053/2798.htm, 1 attachments)
（主校区）卓越工程师学院(https://gs.dlut.edu.cn/info/1173/14076.htm, 4 attachments)
（主校区）力学与航空航天学院(https://lihang.dlut.edu.cn/info/1050/16291.htm, 5 attachments)
（主校区）船舶工程学院(http://naoe.dlut.edu.cn/info/1068/3751.htm, 4 attachments)
电气工程学院(https://eee.dlut.edu.cn/info/1153/6973.htm, 5 attachments)
控制科学与工程(http://scse.dlut.edu.cn/info/1181/3831.htm, 4 attachments)
信息与通信工程学院(https://ice.dlut.edu.cn/info/1208/3186.htm, 7 attachments)
计算机科学与技术学院(https://cs.dlut.edu.cn/info/1432/3238.htm, 4 attachments)
（开发区校区）软件学院(https://ss.dlut.edu.cn/info/1122/26632.htm, 5 attachments)
集成电路学院(http://ic.dlut.edu.cn/info/1034/6957.htm, 5 attachments)
（盘锦校区）化工海洋与生命学院(https://hyxy.dlut.edu.cn/info/1184/7737.htm, 7 attachments)
（盘锦校区）商学院(http://business.dlut.edu.cn/info/1023/4497.htm, 3 attachments)
```
