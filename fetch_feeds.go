package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/LucasCoppola/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Post []Post `xml:"item"`
}

type Post struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
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
	fmt.Println("Start feed scrapping...")
	rss, err := fetchFromFeed(feed.Url)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, post := range rss.Channel.Post {
		var publishedAt time.Time

		if post.PubDate != "" {
			parsedTime, err := parsePublishedAt(post.PubDate)
			if err != nil {
				log.Printf("Error parsing published date for post %s: %v", post.Title, err)
				// set published date to now
				publishedAt = time.Now().UTC()
			} else {
				publishedAt = parsedTime
			}
		}

		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       post.Title,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Url:         post.Link,
			Description: post.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	feed.LastFetchedAt.Time = time.Now().UTC()
	feed.LastFetchedAt.Valid = true
	db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: feed.LastFetchedAt,
		UpdatedAt:     time.Now().UTC(),
		ID:            feed.ID,
	})

	fmt.Println("Finish feed scrapping...")
}

func fetchFromFeed(url string) (RSS, error) {
	fmt.Println("Fetching from feed...")
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

func parsePublishedAt(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	formats := []string{
		time.RFC1123Z,               // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,                // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,                // "02 Jan 06 15:04 -0700"
		time.RFC822,                 // "02 Jan 06 15:04 MST"
		"2006-01-02T15:04:05Z07:00", // ISO 8601 / RFC 3339
		"2006-01-02T15:04:05Z",      // ISO 8601 / RFC 3339 without timezone
		"2006-01-02 15:04:05",       // MySQL datetime format
		"2006-01-02",                // Just date
		"02 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 -0700",
		"02.01.2006 15:04:05",
		"Monday, January 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
