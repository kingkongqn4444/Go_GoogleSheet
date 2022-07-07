package utilities

import (
	google "awesomeProject/src/googleSheets"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strconv"
	"time"
)

type location struct {
	Name string `json:"name"`
	Link string `json:link`
}

type Company struct {
	Name    string `json:"name"`
	Mst     string `json:"mst"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Link    string `json:link`
}

type Companys struct {
	TotalCompany int       `json:"total_company"`
	List         []Company `json:companys`
}

func NewCompany() *Companys {
	return &Companys{}
}

type Person struct {
	Link string    `json:"link"`
	List []Company `json:companys`
}

func (companys *Companys) getCompanyUrl(url string,page int) error {
	//eg := errgroup.Group{}
	t := strconv.Itoa(page)
	doc, err := goquery.NewDocument(url + t)
	if err != nil {
		return err
	}

	doc.Find(".company-name").Each(func(i int, s *goquery.Selection) { // Lặp dữ liệu qua DOM: .col-left ._2pin
		linkDoc, exists := s.Find("a").Attr("href") // Lấy link của ebook
		if !exists {
			linkDoc = "#"
		}
		companys.getCompanyDetail(linkDoc) // Thu thập thông tin company qua url của page
	})
	return nil
}

func (companys *Companys) getCompanyDetail(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	checkPhone := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	doc.Find(".company-info").Each(func(i int, s *goquery.Selection) {
		nameCompany := s.Find(".col-xs-12 .col-md-9").Slice(0, 1).Text()
		mst := s.Find(".col-xs-12 .col-md-9").Slice(2, 3).Text()
		address := s.Find(".col-xs-12 .col-md-9").Slice(6, 7).Text()
		phone := s.Find(".col-xs-12 .col-md-9").Slice(7, 8).Text()
		if checkPhone.MatchString(phone) == true {
			// Gán dữ liệu thu thập được vào struct Ebook
			company := Company{
				Name:    nameCompany,
				Mst:     mst,
				Address: address,
				Link:    url,
				Phone:   phone,
			}
			google.SaveToGoogle(company.Name,company.Mst,company.Address,company.Link,company.Phone)
			companys.TotalCompany += 1
			companys.List = append(companys.List, company)
		}
	})
	return nil
}

type locations struct {
	totalLocation int        `json:"totalLocation"`
	List          []location `json:location`
}

func (locations *locations) getLocationUrl(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	doc.Find(".list-link li").Each(func(i int, s *goquery.Selection) {
		linkLocation, exist := s.Find("a").Attr("href")
		nameLocation := s.Find("a").Text()
		if !exist {
			linkLocation = "#w"
		}
		location := location{
			Name: nameLocation,
			Link: linkLocation,
		}
		locations.totalLocation += 1
		locations.List = append(locations.List, location)
	})
	return nil
}

func checkError(err error) {
	if err != nil {
		log.Print(err)
	}
}


// Lấy danh sach thành phố cần crawle
func ListLocation(c *gin.Context) {
	var locations locations
	err2 := locations.getLocationUrl("https://doanhnghiepmoi.vn")
	c.Header("Access-Control-Allow-Origin","text/html; charset=utf-8")
	c.Header("Content-Type","*")
	checkError(err2)
	json.NewEncoder(c.Writer).Encode(locations)
}

//func ListCompany(c *gin.Context) {
//	decoder := json.NewDecoder(c.Request.Body)
//	var person Person
//	err := decoder.Decode(&person)
//		if err != nil {
//			fmt.Print(err)
//			panic(err)
//		}
//	company := NewCompany()
//	err2 := company.getCompanyUrl(person.Link)
//	checkError(err2)
//	json.NewEncoder(c.Writer).Encode(company)
//}

// lấy danh sách công ty
func ListCompany(c *gin.Context) {
	GetAllPage(787);
	json.NewEncoder(c.Writer).Encode("Done")
}

func GetAllPage(currentPage int){
	if currentPage == 0 {
		return
	} else {
		company := NewCompany()
		err2 := company.getCompanyUrl("https://doanhnghiepmoi.vn/Dac-Lac/trang-",currentPage)
		checkError(err2)
		time.Sleep(3600)
		GetAllPage(currentPage - 1)
	}
}