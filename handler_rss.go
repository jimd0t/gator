package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jimd0t/gator/internal/database"
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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s name url", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feed, err := s.queries.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Printf("Feed entry added correctly! - %v | %s | %s\n", feed.ID, feed.Name, feed.Url)

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
	}

	_, err = s.queries.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.queries.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("- %s (%s) | %s\n", feed.Name, feed.Url, feed.Username)
	}
	return nil
}

func handlerFeedFollows(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s url", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.queries.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	feedFollow, err := s.queries.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		return err
	}
	fmt.Println(feedFollow)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.queries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("List of feeds for %s\n", s.config.CurrentUserName)
	for _, f := range feeds {
		fmt.Printf(" - %s\n", f.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.queries.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	unfollowArgs := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.queries.UnfollowFeed(context.Background(), unfollowArgs)
	if err != nil {
		return err
	}
	return nil
}
