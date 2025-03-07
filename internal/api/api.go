package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/noueii/rss-gator/internal/app"
	"github.com/noueii/rss-gator/internal/db"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func ScrapeFeeds(a *app.App) error {
	fmt.Printf("Fetching feed...\n")
	feed, err := a.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	if err := a.DB.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return err
	}

	fetchedFeed, err := FetchFeed(context.Background(), feed.Url)

	fmt.Printf("Fetched feed %s\n", fetchedFeed.Channel.Title)
	if err := saveFeed(a, *fetchedFeed, feed); err != nil {
		return err
	}
	return nil

}

func saveFeed(a *app.App, feed RSSFeed, dbFeed db.Feed) error {

	for _, item := range feed.Channel.Items {
		// fmt.Printf("\t- Saving post %s\n", item.Title)
		if !isValidPost(item) {
			fmt.Println("Post not valid")
			continue
		}
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err

		}

		if _, err := a.DB.CreatePost(context.Background(), db.CreatePostParams{
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: parsedTime,
			FeedID:      dbFeed.ID,
		}); err != nil {
			if strings.Contains(err.Error(), `duplicate key value violates unique constraint "posts_url_key"`) {
				continue
			}
			return err
		}
	}

	fmt.Println("All posts saved")

	return nil
}

func isValidPost(post RSSItem) bool {
	if len(post.Title) == 0 || len(post.Description) == 0 || len(post.Link) == 0 {
		return false
	}
	return true
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	body, err := fetchXMLWithRequest(ctx, "GET", feedURL)
	if err != nil {
		return nil, err
	}

	rss := RSSFeed{}
	if err = xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	for i := range rss.Channel.Items {
		item := &rss.Channel.Items[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		item.PubDate = html.UnescapeString(item.PubDate)
		item.Link = html.EscapeString(item.Link)
	}

	return &rss, nil
}

func fetchXMLWithRequest(ctx context.Context, method string, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("user-agent", "rss-gator")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f *RSSFeed) Print() {
	fmt.Println(f)
}
