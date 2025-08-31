package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aaronkarr/aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBeweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	fmt.Printf("Collecting feeds every %s...\n", timeBeweenReqs)

	ticker := time.NewTicker(timeBeweenReqs)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("couldn't get next feed to fetch: %w", err)
		return
	}
	fmt.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("couldn't mark feed %s as fetched: %w", feed.Name, err)
	}

	feedItem, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("error fetching feed %s: %w", feed.Name, err)
	}

	for _, item := range feedItem.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	fmt.Printf("Feed %s collected, %v posts found", feed.Name, len(feedItem.Channel.Item))
}
