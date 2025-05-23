package config

import (
	"breyting/blog-aggregator/internal/database"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func HandlerFetchFeed(s *State, cmd Command) error {
	res, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	resFeed := RSSFeed{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &resFeed, err
	}
	req.Header.Set("user-agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &resFeed, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &resFeed, err
	}

	err = xml.Unmarshal(data, &resFeed)
	if err != nil {
		return &resFeed, err
	}

	resFeed.Channel.Title = html.UnescapeString(resFeed.Channel.Title)
	resFeed.Channel.Description = html.UnescapeString(resFeed.Channel.Description)

	for _, val := range resFeed.Channel.Item {
		val.Title = html.UnescapeString(val.Title)
		val.Description = html.UnescapeString(val.Description)
	}

	return &resFeed, nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("add feed needs a name and a url")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.Db.CreateFeeds(context.Background(), database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	follow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return nil
	}

	fmt.Printf("added and starting the feed %s\n", follow.FeedName)
	return nil
}

func HandlerGetFeeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("feeds takes no arguments")
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, val := range feeds {
		fmt.Printf("name : %s, url : %s, user : %s\n", val.FeedName, val.FeedUrl, val.UserName.String)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("follow needs the url of the feed")
	}

	url := cmd.Args[0]

	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("no feed with this url")
	}

	follow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("name current user : %s feed : %s\n", follow.UserName, follow.FeedName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("following takes no arguments")
	}

	follow, err := s.Db.GetFeedFollowsByUsers(context.Background(), user.Name)
	if err != nil {
		return err
	}

	for _, val := range follow {
		fmt.Printf("- %s\n", val.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("unfollow needs an url")
	}
	url := cmd.Args[0]

	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	err = s.Db.DeleteFollowByUserAndUrl(context.Background(), database.DeleteFollowByUserAndUrlParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("succesfully unfollow the feed \"%s\"", feed.Name)
	return nil
}
