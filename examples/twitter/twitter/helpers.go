package twitter

import (
	"strings"
	"time"
	"zen/pkg/twitter"

	"golang.org/x/exp/rand"
)

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

func (k *Twitter) sleepWithInterrupt(duration time.Duration) error {
	k.logger.Infof("Waiting %v until next processing", duration)

	select {
	case <-time.After(duration):
		return nil
	case <-k.ctx.Done():
		return k.ctx.Err()
	}
}

func (k *Twitter) isTweetTooOld(tweet *twitter.ParsedTweet) bool {
	return time.Since(time.Unix(tweet.TweetCreatedAt, 0)) > 300*time.Minute
}

func (k *Twitter) isOwnTweet(username string) bool {
	return strings.ToLower(username) == strings.ToLower(k.twitterConfig.Credentials.User)
}
