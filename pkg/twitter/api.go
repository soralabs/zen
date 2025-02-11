package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

// GetTweetDetails gets the details of a tweet
func (c *Client) GetTweetDetails(tweetID string) (*TweetDetailsResponse, error) {

	variables := map[string]interface{}{
		"focalTweetId":                           tweetID,
		"referrer":                               "home",
		"with_rux_injections":                    false,
		"rankingMode":                            "Relevance",
		"includePromotedContent":                 true,
		"withCommunity":                          true,
		"withQuickPromoteEligibilityTweetFields": true,
		"withBirdwatchNotes":                     true,
		"withVoice":                              true,
	}

	features := map[string]interface{}{
		"rweb_tipjar_consumption_enabled":                                         true,
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"communities_web_enable_tweet_community_results_fetch":                    true,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"articles_preview_enabled":                                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"creator_subscriptions_quote_tweet_preview_enabled":                       false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"rweb_video_timestamps_enabled":                                           true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	fieldToggles := map[string]interface{}{
		"withArticleRichContentState": true,
		"withArticlePlainText":        false,
		"withGrokAnalyze":             false,
		"withDisallowedReplyControls": false,
	}

	data, _ := json.Marshal(variables)
	featuresData, _ := json.Marshal(features)
	fieldTogglesData, _ := json.Marshal(fieldToggles)

	res, err := resty.New().R().
		SetHeaders(map[string]string{
			"sec-ch-ua-platform":        "\"macOS\"",
			"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
			"x-csrf-token":              c.twitterCredential.CT0,
			"x-client-uuid":             "1fc2b334-3a88-4e1d-bdc7-f0390bd14dba",
			"sec-ch-ua":                 "\"Google Chrome\";v=\"131",
			"x-twitter-client-language": "en",
			"sec-ch-ua-mobile":          "?0",
			"x-twitter-active-user":     "yes",
			"x-client-transaction-id":   "iEx0NdN4jP+RSDiS8mIl16L7LrUW/92VH91Vqw2Z+UYUvs0YYQkuYzDD29/trpmGs2o5f4ooVtt9HvuqWjW8k7YpRlDDiw",
			"x-twitter-auth-type":       "OAuth2Session",
			"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"content-type":              "application/json",
			"accept":                    "*/*",
			"sec-fetch-site":            "same-origin",
			"sec-fetch-mode":            "cors",
			"sec-fetch-dest":            "empty",
			"referer":                   fmt.Sprintf("https://x.com/~/status/%s", tweetID),
			"accept-language":           "en-US,en;q=0.9",
			"priority":                  "u=1, i",
		}).
		SetContext(c.ctx).
		SetQueryParams(map[string]string{
			"variables":    string(data),
			"features":     string(featuresData),
			"fieldToggles": string(fieldTogglesData),
		}).
		SetCookies([]*http.Cookie{
			{Name: "ct0", Value: c.twitterCredential.CT0},
			{Name: "auth_token", Value: c.twitterCredential.AuthToken},
		}).
		Get("https://x.com/i/api/graphql/nBS-WpgA6ZG0CyNHD517JQ/TweetDetail")
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get tweet details: %s", res.String())
	}

	var response TweetDetailsResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SearchTimeline searches the timeline for tweets from the given accounts
func (c *Client) SearchTimeline(accounts []string) (*SearchTimelineResponse, error) {

	variables := map[string]interface{}{
		"rawQuery":    fmt.Sprintf("(from:%s) -filter:replies", strings.Join(accounts, " OR from:")),
		"count":       20,
		"querySource": "typed_query",
		"product":     "Latest",
	}

	features := map[string]interface{}{
		"rweb_tipjar_consumption_enabled":                                         true,
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"communities_web_enable_tweet_community_results_fetch":                    true,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"articles_preview_enabled":                                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"creator_subscriptions_quote_tweet_preview_enabled":                       false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"rweb_video_timestamps_enabled":                                           true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	data, _ := json.Marshal(variables)
	featuresData, _ := json.Marshal(features)

	res, err := resty.New().R().
		SetHeaders(map[string]string{
			"Host":                      "x.com",
			"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
			"x-csrf-token":              c.twitterCredential.CT0,
			"x-client-uuid":             "1fc2b334-3a88-4e1d-bdc7-f0390bd14dba",
			"sec-ch-ua":                 "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
			"x-twitter-client-language": "en",
			"sec-ch-ua-mobile":          "?0",
			"x-twitter-active-user":     "yes",
			"x-twitter-auth-type":       "OAuth2Session",
			"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"content-type":              "application/json",
			"accept":                    "*/*",
			"sec-fetch-site":            "same-origin",
			"sec-fetch-mode":            "cors",
			"sec-fetch-dest":            "empty",
			"referer":                   fmt.Sprintf("https://x.com/search?f=live&q=%s&src=typed_query", strings.Join(accounts, " OR ")),
			"accept-language":           "en-US,en;q=0.9",
			"priority":                  "u=1, i",
		}).
		SetContext(c.ctx).
		SetQueryParams(map[string]string{
			"variables": string(data),
			"features":  string(featuresData),
		}).
		SetCookies([]*http.Cookie{
			{Name: "ct0", Value: c.twitterCredential.CT0},
			{Name: "auth_token", Value: c.twitterCredential.AuthToken},
		}).
		Get("https://x.com/i/api/graphql/MJpyQGqgklrVl_0X9gNy3A/SearchTimeline")
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to search timeline: %s", res.String())
	}

	var response SearchTimelineResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateTweet creates a tweet
func (c *Client) CreateTweet(tweetContent string, opts *TweetOptions) (*CreateTweetResponse, error) {

	req := CreateTweetRequest{
		QueryID: "sQhC96QH-hw_8TjKE3kW3g",
		Variables: CreateTweetVariables{
			TweetText:   tweetContent,
			DarkRequest: false,
			Media: struct {
				MediaEntities     []any `json:"media_entities,omitempty"`
				PossiblySensitive bool  `json:"possibly_sensitive,omitempty"`
			}{
				MediaEntities:     []any{},
				PossiblySensitive: false,
			},
			SemanticAnnotationIds: []any{},
		},
		Features: CreateTweetFeatures{
			CommunitiesWebEnableTweetCommunityResultsFetch:                 true,
			C9STweetAnatomyModeratorBadgeEnabled:                           true,
			ResponsiveWebEditTweetAPIEnabled:                               true,
			GraphqlIsTranslatableRwebTweetIsTranslatableEnabled:            true,
			ViewCountsEverywhereAPIEnabled:                                 true,
			LongformNotetweetsConsumptionEnabled:                           true,
			ResponsiveWebTwitterArticleTweetConsumptionEnabled:             true,
			TweetAwardsWebTippingEnabled:                                   false,
			CreatorSubscriptionsQuoteTweetPreviewEnabled:                   false,
			LongformNotetweetsRichTextReadEnabled:                          true,
			LongformNotetweetsInlineMediaEnabled:                           true,
			ArticlesPreviewEnabled:                                         true,
			RwebVideoTimestampsEnabled:                                     true,
			ProfileLabelImprovementsPcfLabelInPostEnabled:                  false,
			RwebTipjarConsumptionEnabled:                                   true,
			ResponsiveWebGraphqlExcludeDirectiveEnabled:                    true,
			VerifiedPhoneLabelEnabled:                                      false,
			FreedomOfSpeechNotReachFetchEnabled:                            true,
			StandardizedNudgesMisinfo:                                      true,
			TweetWithVisibilityResultsPreferGqlLimitedActionsPolicyEnabled: true,
			ResponsiveWebGraphqlSkipUserProfileImageExtensionsEnabled:      false,
			ResponsiveWebGraphqlTimelineNavigationEnabled:                  true,
			ResponsiveWebEnhanceCardsEnabled:                               false,
		},
	}

	if opts != nil && opts.ReplyToTweetID != "" {
		req.Variables.Reply = &struct {
			InReplyToTweetID    string `json:"in_reply_to_tweet_id,omitempty"`
			ExcludeReplyUserIds []any  `json:"exclude_reply_user_ids,omitempty"`
		}{
			InReplyToTweetID:    opts.ReplyToTweetID,
			ExcludeReplyUserIds: []any{},
		}
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	res, err := resty.New().R().
		SetHeaders(map[string]string{
			"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
			"x-csrf-token":              c.twitterCredential.CT0,
			"content-type":              "application/json",
			"x-twitter-auth-type":       "OAuth2Session",
			"x-twitter-client-language": "en",
			"x-twitter-active-user":     "yes",
			"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"referer":                   "https://x.com/compose/post",
		}).
		SetContext(c.ctx).
		SetBody(reqBody).
		SetCookies([]*http.Cookie{
			{Name: "ct0", Value: c.twitterCredential.CT0},
			{Name: "auth_token", Value: c.twitterCredential.AuthToken},
		}).
		Post("https://x.com/i/api/graphql/hOcJ79jSQLm1A06VxjNtZQ/CreateTweet")
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("invalid status code creating tweet: %s", res.String())
	}

	var response CreateTweetResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Data.CreateTweet.TweetResults.Result.Legacy.IDStr == "" {
		return nil, fmt.Errorf("failed to create tweet: %s", res.String())
	}

	// log.Println(res.String())

	return &response, nil
}

// FavoriteTweet favorites a tweet
func (c *Client) FavoriteTweet(tweetID string) error {
	req := FavoriteTweetRequest{
		Variables: struct {
			TweetID string `json:"tweet_id,omitempty"`
		}{
			TweetID: tweetID,
		},
		QueryID: "lI07N6Otwv1PhnEgXILM7A",
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	res, err := resty.New().R().
		SetHeaders(map[string]string{
			"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
			"x-csrf-token":              c.twitterCredential.CT0,
			"content-type":              "application/json",
			"x-twitter-auth-type":       "OAuth2Session",
			"x-twitter-client-language": "en",
			"x-twitter-active-user":     "yes",
			"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"referer":                   "https://x.com/compose/post",
		}).
		SetContext(c.ctx).
		SetBody(reqBody).
		SetCookies([]*http.Cookie{
			{Name: "ct0", Value: c.twitterCredential.CT0},
			{Name: "auth_token", Value: c.twitterCredential.AuthToken},
		}).
		Post("https://x.com/i/api/graphql/lI07N6Otwv1PhnEgXILM7A/FavoriteTweet")
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("invalid status code creating tweet: %s", res.String())
	}

	var response FavoriteTweetResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Data.FavoriteTweet != "Done" {
		return fmt.Errorf("failed to favorite tweet: %s", res.String())
	}

	return nil
}

func (c *Client) GetUserDetails(username string) (*GetUserDetailsResponse, error) {
	variables := map[string]interface{}{
		"screen_name": username,
	}

	features := map[string]interface{}{
		"hidden_profile_subscriptions_enabled":                              true,
		"profile_label_improvements_pcf_label_in_post_enabled":              true,
		"rweb_tipjar_consumption_enabled":                                   true,
		"responsive_web_graphql_exclude_directive_enabled":                  true,
		"verified_phone_label_enabled":                                      false,
		"subscriptions_verification_info_is_identity_verified_enabled":      true,
		"subscriptions_verification_info_verified_since_enabled":            true,
		"highlights_tweets_tab_ui_enabled":                                  true,
		"responsive_web_twitter_article_notes_tab_enabled":                  true,
		"subscriptions_feature_can_gift_premium":                            true,
		"creator_subscriptions_tweet_preview_api_enabled":                   true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled": false,
		"responsive_web_graphql_timeline_navigation_enabled":                true,
	}

	fieldToggles := map[string]interface{}{
		"withAuxiliaryUserLabels": false,
	}

	data, _ := json.Marshal(variables)
	featuresData, _ := json.Marshal(features)
	fieldTogglesData, _ := json.Marshal(fieldToggles)

	res, err := resty.New().R().
		SetHeaders(map[string]string{
			"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
			"x-csrf-token":              c.twitterCredential.CT0,
			"x-twitter-client-language": "en",
			"x-twitter-active-user":     "yes",
			"x-twitter-auth-type":       "OAuth2Session",
			"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"content-type":              "application/json",
			"accept":                    "*/*",
			"referer":                   fmt.Sprintf("https://x.com/%s", username),
		}).
		SetContext(c.ctx).
		SetQueryParams(map[string]string{
			"variables":    string(data),
			"features":     string(featuresData),
			"fieldToggles": string(fieldTogglesData),
		}).
		SetCookies([]*http.Cookie{
			{Name: "ct0", Value: c.twitterCredential.CT0},
			{Name: "auth_token", Value: c.twitterCredential.AuthToken},
		}).
		Get("https://x.com/i/api/graphql/32pL5BWe9WKeSK1MoPvFQQ/UserByScreenName")
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get user details: %s", res.String())
	}

	var response GetUserDetailsResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SearchReplies searches for replies to the specified user
func (c *Client) SearchReplies(username string, limit int) (*SearchTimelineResponse, error) {
	variables := map[string]interface{}{
		"rawQuery":    fmt.Sprintf("to:%s", username),
		"count":       limit,
		"querySource": "typed_query",
		"product":     "Latest",
	}

	features := map[string]interface{}{
		"rweb_tipjar_consumption_enabled":                                         true,
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"communities_web_enable_tweet_community_results_fetch":                    true,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"articles_preview_enabled":                                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"creator_subscriptions_quote_tweet_preview_enabled":                       false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"rweb_video_timestamps_enabled":                                           true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	data, _ := json.Marshal(variables)
	featuresData, _ := json.Marshal(features)

	res, err := resty.New().R().
		SetHeaders(map[string]string{
			"Host":                      "x.com",
			"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
			"x-csrf-token":              c.twitterCredential.CT0,
			"x-client-uuid":             "1fc2b334-3a88-4e1d-bdc7-f0390bd14dba",
			"sec-ch-ua":                 "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
			"x-twitter-client-language": "en",
			"sec-ch-ua-mobile":          "?0",
			"x-twitter-active-user":     "yes",
			"x-twitter-auth-type":       "OAuth2Session",
			"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"content-type":              "application/json",
			"accept":                    "*/*",
			"sec-fetch-site":            "same-origin",
			"sec-fetch-mode":            "cors",
			"sec-fetch-dest":            "empty",
			"referer":                   fmt.Sprintf("https://x.com/search?q=to:%s&src=typed_query&f=live", username),
			"accept-language":           "en-US,en;q=0.9",
			"priority":                  "u=1, i",
		}).
		SetContext(c.ctx).
		SetQueryParams(map[string]string{
			"variables": string(data),
			"features":  string(featuresData),
		}).
		SetCookies([]*http.Cookie{
			{Name: "ct0", Value: c.twitterCredential.CT0},
			{Name: "auth_token", Value: c.twitterCredential.AuthToken},
		}).
		Get("https://x.com/i/api/graphql/MJpyQGqgklrVl_0X9gNy3A/SearchTimeline")
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to search replies: %s", res.String())
	}

	var response SearchTimelineResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, err
	}

	return &response, nil
}
