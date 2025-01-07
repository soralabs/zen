package twitter

import (
	"strings"
	"time"

	"github.com/soralabs/zen/pkg/twitter"

	"golang.org/x/exp/rand"
)

// getRandomInterval returns a random duration between the configured Min and Max intervals
func (k *Twitter) getRandomInterval() time.Duration {
	min := k.twitterConfig.MonitorInterval.Min
	max := k.twitterConfig.MonitorInterval.Max

	if min == max {
		return min
	}

	// Convert to int64 nanoseconds for random calculation
	minNanos := min.Nanoseconds()
	maxNanos := max.Nanoseconds()

	// Generate random duration between min and max
	randomNanos := minNanos + rand.Int63n(maxNanos-minNanos+1)

	return time.Duration(randomNanos)
}

// sleepWithInterrupt waits for the specified duration unless the context is canceled
func (k *Twitter) sleepWithInterrupt(duration time.Duration) error {
	k.logger.Infof("Waiting %v until next processing", duration)

	select {
	case <-time.After(duration):
		return nil
	case <-k.ctx.Done():
		return k.ctx.Err()
	}
}

// isTweetTooOld checks if the tweet's creation time is older than 300 minutes
func (k *Twitter) isTweetTooOld(tweet *twitter.ParsedTweet) bool {
	return time.Since(time.Unix(tweet.TweetCreatedAt, 0)) > 300*time.Minute
}

// isOwnTweet verifies if the provided username matches the authenticated user's name
func (k *Twitter) isOwnTweet(username string) bool {
	return strings.ToLower(username) == strings.ToLower(k.twitterConfig.Credentials.User)
}
