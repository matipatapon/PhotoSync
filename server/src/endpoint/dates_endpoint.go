package endpoint

import (
	"encoding/json"
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/jwt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Date struct {
	Date      string `json:"date"`
	FileCount string `json:"file_count"`
}

var query_without_filtration string = "SELECT TO_CHAR(creation_date, 'YYYY.MM.DD') AS date, COUNT(*) AS file_count FROM files WHERE user_id = $1 GROUP BY date ORDER BY date DESC"
var query_with_year_filtration string = "SELECT TO_CHAR(creation_date, 'YYYY.MM.DD') AS date, COUNT(*) AS file_count FROM files WHERE user_id = $1 AND DATE_PART('year', creation_date) = $2 GROUP BY date ORDER BY date DESC"
var query_with_year_and_month_filtration string = "SELECT TO_CHAR(creation_date, 'YYYY.MM.DD') AS date, COUNT(*) AS file_count FROM files WHERE user_id = $1 AND DATE_PART('year', creation_date) = $2 AND DATE_PART('month', creation_date) = $3 GROUP BY date ORDER BY date DESC"

type DatesEndpoint struct {
	db     database.IDataBase
	jm     jwt.IJwtManager
	logger *log.Logger
}

func NewDatesEndpoint(db database.IDataBase, jm jwt.IJwtManager) DatesEndpoint {
	return DatesEndpoint{db: db, jm: jm, logger: log.New(os.Stdout, "[DatesEndpoint]: ", log.LstdFlags)}
}

func (fe *DatesEndpoint) Options(c *gin.Context) {
	// c.Header("Access-Control-Allow-Headers", "Authorization")
	// c.Header("Access-Control-Allow-Methods", "GET")
	// c.Status(200)
}

func (de *DatesEndpoint) Get(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	jwt, err := de.jm.Decode(token)
	if err != nil {
		c.Status(403)
		de.logger.Printf("Token is invalid: '%s'", err.Error())
		return
	}

	query := query_without_filtration
	queryParams := []any{jwt.UserId}

	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr != "" {
		year, err := strconv.ParseInt(yearStr, 10, 64)
		if err != nil {
			c.Status(400)
			de.logger.Printf("Invalid year: '%s'", err.Error())
			return
		}
		if monthStr != "" {
			month, err := strconv.ParseInt(monthStr, 10, 64)
			if err != nil {
				c.Status(400)
				de.logger.Printf("Invalid month: '%s'", err.Error())
				return
			}
			query = query_with_year_and_month_filtration
			queryParams = append(queryParams, year, month)
		} else {
			query = query_with_year_filtration
			queryParams = append(queryParams, year)
		}
	} else {
		if monthStr != "" {
			c.Status(400)
			de.logger.Print("Month specified without a year")
			return
		}
	}

	rows, err := de.db.Query(query, queryParams...)
	if err != nil {
		c.Status(500)
		de.logger.Printf("Query failed: '%s'", err.Error())
		return
	}

	body := []Date{}
	for _, row := range rows {
		date := row[0].(string)
		fileCount := strconv.FormatInt(row[1].(int64), 10)
		body = append(body, Date{Date: date, FileCount: fileCount})
	}

	de.logger.Printf("Successfully gathered '%d' dates for '%s'", len(body), jwt.Username)

	bytes, _ := json.Marshal(body)
	c.Writer.Write(bytes)
}
