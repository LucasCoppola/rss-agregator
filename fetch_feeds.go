package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
}

func fetchFeedWorker(db *database.Queries, collectionConcurrency int, collectionInterval time.Duration) {
	ticker := time.NewTicker(collectionInterval)

	for range ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(collectionConcurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}

		log.Printf("Found %v feeds to fetch!", len(feeds))

		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, &wg, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	rss, err := fetchFromFeed(feed.Url)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Feed Title: %s\n", rss.Channel.Title)

	feed.LastFetchedAt.Time = time.Now().UTC()
	feed.LastFetchedAt.Valid = true
	db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: feed.LastFetchedAt,
		UpdatedAt:     time.Now().UTC(),
		ID:            feed.ID,
	})
}

func fetchFromFeed(url string) (RSS, error) {
	res, err := http.Get(url)

	if err != nil {
		return RSS{}, fmt.Errorf("GET error: %v", err)
	}

	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)
	var rss RSS
	err = decoder.Decode(&rss)

	if err != nil {
		return RSS{}, fmt.Errorf("Decode data error: %v", err)
	}

	return rss, nil
}
