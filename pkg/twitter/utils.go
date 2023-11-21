package twitter

import (
	"fmt"
	"time"
)

type ParsedTweet struct {
	// User Info
	UserID      string `json:"user_id,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`

	// Tweet Info
	TweetCreatedAt      int64         `json:"tweet_created_at,omitempty"`
	TweetID             string        `json:"tweet_id,omitempty"`
	TweetConversationID string        `json:"tweet_conversation_id,omitempty"`
	TweetText           string        `json:"tweet_text,omitempty"`
	TweetImages         []string      `json:"tweet_images,omitempty"`
	TweetLinks          []string      `json:"tweet_links,omitempty"`
	TweetReplies        []ParsedTweet `json:"reply_tweets,omitempty"`

	InReplyToTweetID    string `json:"in_reply_to_tweet_id,omitempty"`
	InReplyToScreenName string `json:"in_reply_to_screen_name,omitempty"`
}

func (c *Client) ParseSearchTimelineResponse(res *SearchTimelineResponse) ([]*ParsedTweet, error) {

	var parsedTweets []*ParsedTweet
	for _, instruction := range res.Data.SearchByRawQuery.SearchTimeline.Timeline.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				if entry.Content.EntryType == "TimelineTimelineItem" || entry.Content.EntryType == "TimelineTimelineModule" {
					timestamp, err := parseTwitterTimestamp(entry.Content.ItemContent.TweetResults.Result.Legacy.CreatedAt)
					if err != nil {
						return nil, fmt.Errorf("failed to parse tweet timestamp: %w", err)
					}

					var urls []string
					for _, url := range entry.Content.ItemContent.TweetResults.Result.Legacy.Entities.Urls {
						urlInterface, ok := url.(map[string]interface{})
						if !ok {
							continue
						}
						urls = append(urls, urlInterface["expanded_url"].(string))
					}

					var images []string
					for _, media := range entry.Content.ItemContent.TweetResults.Result.Legacy.Entities.Media {
						if media.Type == "photo" {
							images = append(images, media.MediaURLHTTPS)
						}
					}

					parsedTweetInfo := &ParsedTweet{
						UserID:              entry.Content.ItemContent.TweetResults.Result.Core.UserResults.Result.RestID,
						UserName:            entry.Content.ItemContent.TweetResults.Result.Core.UserResults.Result.Legacy.ScreenName,
						DisplayName:         entry.Content.ItemContent.TweetResults.Result.Core.UserResults.Result.Legacy.Name,
						TweetCreatedAt:      timestamp,
						TweetID:             entry.Content.ItemContent.TweetResults.Result.Legacy.IDStr,
						TweetConversationID: entry.Content.ItemContent.TweetResults.Result.Legacy.ConversationIDStr,
						TweetText:           entry.Content.ItemContent.TweetResults.Result.Legacy.FullText,
						TweetImages:         images,
						TweetLinks:          urls,
						InReplyToTweetID:    entry.Content.ItemContent.TweetResults.Result.Legacy.InReplyToStatusIDStr,
						InReplyToScreenName: entry.Content.ItemContent.TweetResults.Result.Legacy.InReplyToScreenName,
					}
					parsedTweets = append(parsedTweets, parsedTweetInfo)
				}
			}
		}
	}

	return parsedTweets, nil
}

// ParseTweetReplies parses the replies of a tweet. userIDFilter is the user ID of the user whose replies we want to parse.
// If userIDFilter is empty, all replies are parsed.
func (c *Client) ParseTweetReplies(res *TweetDetailsResponse, userIDFilter string) ([]ParsedTweet, error) {
	var parsedReplies []ParsedTweet
	for _, reply := range res.Data.ThreadedConversationWithInjectionsV2.Instructions {
		if reply.Type == "TimelineAddEntries" {
			for _, entry := range reply.Entries {
				if entry.Content.EntryType == "TimelineTimelineModule" {
					for _, item := range entry.Content.Items {
						if userIDFilter == "" || item.Item.ItemContent.TweetResults.Result.Legacy.UserIDStr == userIDFilter {
							if item.Item.ItemContent.TweetResults.Result.Legacy.FullText == "" {
								continue
							}

							timestamp, err := parseTwitterTimestamp(item.Item.ItemContent.TweetResults.Result.Legacy.CreatedAt)
							if err != nil {
								return nil, fmt.Errorf("failed to parse tweet timestamp: %w", err)
							}

							var links []string
							for _, url := range item.Item.ItemContent.TweetResults.Result.Legacy.Entities.Urls {
								links = append(links, url.ExpandedURL)
							}

							var images []string
							for _, media := range item.Item.ItemContent.TweetResults.Result.Legacy.Entities.Media {
								if media.Type == "photo" {
									images = append(images, media.MediaURLHTTPS)
								}
							}

							parsedReplies = append(parsedReplies, ParsedTweet{
								UserID:              item.Item.ItemContent.TweetResults.Result.Core.UserResults.Result.RestID,
								UserName:            item.Item.ItemContent.TweetResults.Result.Core.UserResults.Result.Legacy.ScreenName,
								DisplayName:         item.Item.ItemContent.TweetResults.Result.Core.UserResults.Result.Legacy.Name,
								TweetCreatedAt:      timestamp,
								TweetID:             item.Item.ItemContent.TweetResults.Result.Legacy.IDStr,
								TweetConversationID: item.Item.ItemContent.TweetResults.Result.Legacy.ConversationIDStr,
								TweetText:           item.Item.ItemContent.TweetResults.Result.Legacy.FullText,
								TweetImages:         images,
								TweetLinks:          links,
								InReplyToTweetID:    item.Item.ItemContent.TweetResults.Result.Legacy.InReplyToStatusIDStr,
								InReplyToScreenName: item.Item.ItemContent.TweetResults.Result.Legacy.InReplyToScreenName,
							})
						}
					}
				}
			}
		}
	}
	return parsedReplies, nil
}

// IsReply checks if a tweet is a reply to another tweet
func IsReply(tweetID string, tweetDetailsRes *TweetDetailsResponse) bool {
	for _, instruction := range tweetDetailsRes.Data.ThreadedConversationWithInjectionsV2.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				if entry.Content.Typename == "TimelineTimelineItem" && entry.Content.ItemContent.TweetResults.Result.RestID == tweetID {
					entities := entry.Content.ItemContent.TweetResults.Result.Legacy.Entities
					if len(entities.UserMentions) > 0 {
						return true
					}
				}
			}
		}
	}
	return false
}

// GetRootTweetID gets the tweet ID of the root tweet of a conversation
func GetRootTweetID(tweetID string, tweetDetailsRes *TweetDetailsResponse) string {
	for _, instruction := range tweetDetailsRes.Data.ThreadedConversationWithInjectionsV2.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for i := range instruction.Entries {
				if instruction.Entries[i].Content.Typename == "TimelineTimelineItem" {
					return instruction.Entries[i].Content.ItemContent.TweetResults.Result.RestID
				}
			}
		}
	}
	return ""
}

func parseTwitterTimestamp(timestamp string) (int64, error) {
	layout := "Mon Jan 2 15:04:05 -0700 2006"
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return 0, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	return t.Unix(), nil
}

// ParseTweet parses a single tweet from a TweetDetailsResponse
func (c *Client) ParseTweet(res *TweetDetailsResponse) (*ParsedTweet, error) {
	for _, instruction := range res.Data.ThreadedConversationWithInjectionsV2.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				if entry.Content.Typename == "TimelineTimelineItem" {
					timestamp, err := parseTwitterTimestamp(entry.Content.ItemContent.TweetResults.Result.Legacy.CreatedAt)
					if err != nil {
						return nil, fmt.Errorf("failed to parse tweet timestamp: %w", err)
					}

					var links []string
					for _, url := range entry.Content.ItemContent.TweetResults.Result.Legacy.Entities.Urls {
						links = append(links, url.ExpandedURL)
					}

					var images []string
					for _, media := range entry.Content.ItemContent.TweetResults.Result.Legacy.Entities.Media {
						if media.Type == "photo" {
							images = append(images, media.MediaURLHTTPS)
						}
					}

					return &ParsedTweet{
						UserID:              entry.Content.ItemContent.TweetResults.Result.Core.UserResults.Result.RestID,
						UserName:            entry.Content.ItemContent.TweetResults.Result.Core.UserResults.Result.Legacy.ScreenName,
						DisplayName:         entry.Content.ItemContent.TweetResults.Result.Core.UserResults.Result.Legacy.Name,
						TweetCreatedAt:      timestamp,
						TweetID:             entry.Content.ItemContent.TweetResults.Result.Legacy.IDStr,
						TweetConversationID: entry.Content.ItemContent.TweetResults.Result.Legacy.ConversationIDStr,
						TweetText:           entry.Content.ItemContent.TweetResults.Result.Legacy.FullText,
						TweetImages:         images,
						TweetLinks:          links,
						InReplyToTweetID:    entry.Content.ItemContent.TweetResults.Result.Legacy.InReplyToStatusIDStr,
						InReplyToScreenName: entry.Content.ItemContent.TweetResults.Result.Legacy.InReplyToScreenName,
					}, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("tweet not found in response")
}
