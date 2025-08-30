package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aaronkarr/aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error retrieving feed: %w", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create follow entry: %w", err)
	}
	println("Following:")
	fmt.Printf(" * Feed: %s\n", follow.FeedName)
	fmt.Printf(" * User: %s\n", follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving feed list: %w", err)
	}

	if len(feeds) > 0 {
		fmt.Printf("%s is following:\n", user.Name)
		for i, feed := range feeds {
			fmt.Printf("Feed %v: %s\n", i+1, feed.FeedName)
		}
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	err := s.db.RemoveFollow(context.Background(), database.RemoveFollowParams{
		Name: user.Name,
		Url:  cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("couldn't unfollow: %w", err)
	}
	fmt.Printf("Unfollowed %s", cmd.Args[0])
	return nil
}
