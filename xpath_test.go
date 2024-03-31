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

// func (c *College) Download(path string) {
// 	u, _ := url.Parse(c.Url)
// 	for _, a := range c.Attachments {
// 		u.Path = a.Url
// 		referer := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
// 		resp := request.Get(u.String(), request.SetReferer(referer))
// 		resp.Write(filepath.Join(path, c.Name, a.Name))
// 	}
// }

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
