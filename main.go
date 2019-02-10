package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/olivere/elastic"
)

// StatData response data for stat api
type StatData struct {
	MID       int `json:"mid"`
	Following int `json:"following"`
	Follower  int `json:"follower"`
}

// Stat response for stat api
type Stat struct {
	Code int      `json:"code"`
	Data StatData `json:"data"`
}

// Record struct
type Record struct {
	UserID    int       `json:"user_id"`
	Following int       `json:"following"`
	Follower  int       `json:"follower"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	// UserID user id
	UserID = os.Getenv("USER_ID")

	// ESHost elasticsearch host
	ESHost = os.Getenv("ES_HOST")
)

func main() {
	var err error
	var client *elastic.Client

	if client, err = elastic.NewClient(
		elastic.SetURL(ESHost),
		elastic.SetSniff(false),
	); err != nil {
		panic(err)
	}

	var resp *http.Response
	if resp, err = http.Get(fmt.Sprintf("https://api.bilibili.com/x/relation/stat?vmid=%s", UserID)); err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var buf []byte
	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	}

	var stat Stat
	if err = json.Unmarshal(buf, &stat); err != nil {
		panic(err)
	}

	if stat.Code != 0 {
		panic(errors.New("code != 0"))
	}

	record := &Record{
		UserID:    stat.Data.MID,
		Follower:  stat.Data.Follower,
		Following: stat.Data.Following,
		Timestamp: time.Now(),
	}

	if _, err = client.Index().Index(indexName(record.Timestamp)).Type("_doc").BodyJson(record).Do(context.Background()); err != nil {
		panic(err)
	}
}

func indexName(t time.Time) string {
	return fmt.Sprintf("fanstrack-%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
