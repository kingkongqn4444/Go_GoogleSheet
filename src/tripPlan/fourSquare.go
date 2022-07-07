package tripplan

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AutoCompletePlace(c *gin.Context) {
	name := c.Param("keyword")
	url := "https://api.foursquare.com/v3/autocomplete?query="
	finalUrl := url + name
	req, _ := http.NewRequest("GET", finalUrl, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "fsq3WV/Wi9Ar7WUZIvLMy1CTIZdeNYA37Sg4+xlMo/PmWw4=")

	res, _ := http.DefaultClient.Do(req)
	contentLength := res.ContentLength
	contentType := res.Header.Get("Content-Type")

	reader := res.Body
	defer reader.Close()

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="gopher.png"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
