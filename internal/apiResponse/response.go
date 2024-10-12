package apiResponse

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
)

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	TimeStamp int64  `json:"timeStamp"`
}

func Fail(message string) *Response {
	return &Response{Code: 1, Message: message, Data: "", TimeStamp: time.Now().Unix()}
}

func Success(data any, message string) *Response {
	return &Response{Code: 0, Message: message, Data: data, TimeStamp: time.Now().Unix()}
}

func (response *Response) SetCode(code int) *Response {
	response.Code = code
	return response
}

func (response *Response) Json(c *gin.Context) {
	//locale := lang.GetLang(c)
	//if response.Message != "" {
	//	response.Message = lang.T(locale, response.Message)
	//}
	c.Set("message", response.Message)
	c.Set("code", response.Code)
	c.AsciiJSON(200, response)
}

func (response *Response) JsonF(c *gin.Context, v ...interface{}) {
	//locale := lang.GetLang(c)
	//if response.Message != "" {
	//	response.Message = fmt.Sprintf(lang.T(locale, response.Message), v)
	//}

	c.Set("message", response.Message)
	c.Set("code", response.Code)
	c.AsciiJSON(200, response)
}

func QueryConditionWhere(c *gin.Context, query *gorm.DB, queryParams map[string]string) *gorm.DB {
	for param, condition := range queryParams {
		if value := c.Query(param); value != "" {
			if strings.Contains(condition, "OR") {
				q := strings.Split(condition, "OR")
				placeholderCount := strings.Count(condition, "?")
				values := make([]interface{}, placeholderCount)
				for i := 0; i < placeholderCount; i++ {
					if strings.Contains(q[i], "LIKE") {
						values[i] = "%" + value + "%"
					} else {
						values[i] = value
					}
				}
				query = query.Where(gorm.Expr(condition, values...))
			} else {
				if strings.Contains(condition, "LIKE") {
					query = query.Where(condition, "%"+value+"%")
				} else {
					query = query.Where(condition, value)
				}
			}
		}

		if p, ok := c.Get(param); ok {
			if value, valid := p.([]string); valid {
				finalParam := make([]string, 0, len(value))
				for _, v := range value {
					if v != "" {
						finalParam = append(finalParam, v)
					}
				}
				if len(finalParam) > 0 {
					query = query.Where(condition, finalParam)
				}
			}
		}
	}

	return query
}

func QueryDateCustomParam(c *gin.Context, start_query string, end_query string, isFilter bool) (carbon.Carbon, carbon.Carbon, bool) {
	startTime := carbon.Parse(c.DefaultQuery(start_query, ""))
	endTime := carbon.Parse(c.DefaultQuery(end_query, ""))
	exist := true
	if startTime.Timestamp() == 0 {
		startTime = carbon.Now()
		exist = false
	}
	if endTime.Timestamp() == 0 {
		endTime = carbon.Now()
		exist = false
	}
	start := startTime
	end := endTime
	if isFilter {
		start = startTime.StartOfDay()
		end = endTime.EndOfDay()
	}
	return start, end, exist
}

func QueryDate(c *gin.Context) (carbon.Carbon, carbon.Carbon, bool) {
	startTime := carbon.Parse(c.DefaultQuery("start_time", ""))
	endTime := carbon.Parse(c.DefaultQuery("end_time", ""))
	exist := true
	if startTime.Timestamp() == 0 {
		startTime = carbon.Now()
		exist = false
	}
	if endTime.Timestamp() == 0 {
		endTime = carbon.Now()
		exist = false
	}
	start := startTime.StartOfDay()
	end := endTime.EndOfDay()
	return start, end, exist
}

func QueryDateTime(c *gin.Context) (string, string, bool) {
	startTime := carbon.Parse(c.DefaultQuery("start_time", ""))
	endTime := carbon.Parse(c.DefaultQuery("end_time", ""))
	exist := true
	if startTime.Timestamp() == 0 {
		startTime = carbon.Now()
		exist = false
	}
	if endTime.Timestamp() == 0 {
		endTime = carbon.Now()
		exist = false
	}

	return startTime.ToDateTimeString(), endTime.ToDateTimeString(), exist
}

func PostDateTime(startTime, endTime string) (carbon.Carbon, carbon.Carbon, bool) {
	start := carbon.Parse(startTime)
	end := carbon.Parse(endTime)
	exist := true

	if start.IsZero() {
		start = carbon.Now()
		exist = false
	}

	if end.IsZero() {
		end = carbon.Now()
		exist = false
	}

	start = start.StartOfDay()
	end = end.EndOfDay()

	return start, end, exist
}
