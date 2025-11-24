package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	rssFeed := RSSFeed{}
	if err != nil {
		return &rssFeed, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req.Header.Set("User-Anget", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &rssFeed, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &rssFeed, err
	}
	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return &rssFeed, err
	}
	// fmt.Println(rssFeed)
	rssFeed.parseFeed()
	return &rssFeed, err
}

func (rssFeed *RSSFeed) parseFeed() {
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}
}

func handlerAgg(s *state, cmd command) error {
	// if len(cmd.Args) < 2 {
	// 	return fmt.Errorf("usage: %s <rss-url>", cmd.Name)
	// }
	// url := cmd.Args[0]
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
