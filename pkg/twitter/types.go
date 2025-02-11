package twitter

type TweetDetailsResponse struct {
	Data struct {
		ThreadedConversationWithInjectionsV2 struct {
			Instructions []struct {
				Type    string `json:"type,omitempty"`
				Entries []struct {
					EntryID   string `json:"entryId,omitempty"`
					SortIndex string `json:"sortIndex,omitempty"`
					Content   struct {
						EntryType   string `json:"entryType,omitempty"`
						Typename    string `json:"__typename,omitempty"`
						ItemContent struct {
							ItemType     string `json:"itemType,omitempty"`
							Typename     string `json:"__typename,omitempty"`
							TweetResults struct {
								Result struct {
									Typename          string `json:"__typename,omitempty"`
									RestID            string `json:"rest_id,omitempty"`
									HasBirdwatchNotes bool   `json:"has_birdwatch_notes,omitempty"`
									Core              struct {
										UserResults struct {
											Result struct {
												Typename                   string `json:"__typename,omitempty"`
												ID                         string `json:"id,omitempty"`
												RestID                     string `json:"rest_id,omitempty"`
												AffiliatesHighlightedLabel struct {
												} `json:"affiliates_highlighted_label,omitempty"`
												HasGraduatedAccess bool   `json:"has_graduated_access,omitempty"`
												IsBlueVerified     bool   `json:"is_blue_verified,omitempty"`
												ProfileImageShape  string `json:"profile_image_shape,omitempty"`
												Legacy             struct {
													FollowedBy          bool   `json:"followed_by,omitempty"`
													Following           bool   `json:"following,omitempty"`
													CanDm               bool   `json:"can_dm,omitempty"`
													CanMediaTag         bool   `json:"can_media_tag,omitempty"`
													CreatedAt           string `json:"created_at,omitempty"`
													DefaultProfile      bool   `json:"default_profile,omitempty"`
													DefaultProfileImage bool   `json:"default_profile_image,omitempty"`
													Description         string `json:"description,omitempty"`
													Entities            struct {
														Description struct {
															Urls []struct {
																DisplayURL  string `json:"display_url,omitempty"`
																ExpandedURL string `json:"expanded_url,omitempty"`
																URL         string `json:"url,omitempty"`
																Indices     []int  `json:"indices,omitempty"`
															} `json:"urls,omitempty"`
														} `json:"description,omitempty"`
														URL struct {
															Urls []struct {
																DisplayURL  string `json:"display_url,omitempty"`
																ExpandedURL string `json:"expanded_url,omitempty"`
																URL         string `json:"url,omitempty"`
																Indices     []int  `json:"indices,omitempty"`
															} `json:"urls,omitempty"`
														} `json:"url,omitempty"`
													} `json:"entities,omitempty"`
													FastFollowersCount      int      `json:"fast_followers_count,omitempty"`
													FavouritesCount         int      `json:"favourites_count,omitempty"`
													FollowersCount          int      `json:"followers_count,omitempty"`
													FriendsCount            int      `json:"friends_count,omitempty"`
													HasCustomTimelines      bool     `json:"has_custom_timelines,omitempty"`
													IsTranslator            bool     `json:"is_translator,omitempty"`
													ListedCount             int      `json:"listed_count,omitempty"`
													Location                string   `json:"location,omitempty"`
													MediaCount              int      `json:"media_count,omitempty"`
													Name                    string   `json:"name,omitempty"`
													NormalFollowersCount    int      `json:"normal_followers_count,omitempty"`
													PinnedTweetIdsStr       []string `json:"pinned_tweet_ids_str,omitempty"`
													PossiblySensitive       bool     `json:"possibly_sensitive,omitempty"`
													ProfileBannerURL        string   `json:"profile_banner_url,omitempty"`
													ProfileImageURLHTTPS    string   `json:"profile_image_url_https,omitempty"`
													ProfileInterstitialType string   `json:"profile_interstitial_type,omitempty"`
													ScreenName              string   `json:"screen_name,omitempty"`
													StatusesCount           int      `json:"statuses_count,omitempty"`
													TranslatorType          string   `json:"translator_type,omitempty"`
													URL                     string   `json:"url,omitempty"`
													Verified                bool     `json:"verified,omitempty"`
													WantRetweets            bool     `json:"want_retweets,omitempty"`
													WithheldInCountries     []any    `json:"withheld_in_countries,omitempty"`
												} `json:"legacy,omitempty"`
												Professional struct {
													RestID           string `json:"rest_id,omitempty"`
													ProfessionalType string `json:"professional_type,omitempty"`
													Category         []struct {
														ID       int    `json:"id,omitempty"`
														Name     string `json:"name,omitempty"`
														IconName string `json:"icon_name,omitempty"`
													} `json:"category,omitempty"`
												} `json:"professional,omitempty"`
												TipjarSettings struct {
													IsEnabled      bool   `json:"is_enabled,omitempty"`
													EthereumHandle string `json:"ethereum_handle,omitempty"`
												} `json:"tipjar_settings,omitempty"`
												SuperFollowEligible bool `json:"super_follow_eligible,omitempty"`
											} `json:"result,omitempty"`
										} `json:"user_results,omitempty"`
									} `json:"core,omitempty"`
									UnmentionData struct {
									} `json:"unmention_data,omitempty"`
									EditControl struct {
										EditTweetIds       []string `json:"edit_tweet_ids,omitempty"`
										EditableUntilMsecs string   `json:"editable_until_msecs,omitempty"`
										IsEditEligible     bool     `json:"is_edit_eligible,omitempty"`
										EditsRemaining     string   `json:"edits_remaining,omitempty"`
									} `json:"edit_control,omitempty"`
									IsTranslatable bool `json:"is_translatable,omitempty"`
									Views          struct {
										Count string `json:"count,omitempty"`
										State string `json:"state,omitempty"`
									} `json:"views,omitempty"`
									Source string `json:"source,omitempty"`
									Legacy struct {
										BookmarkCount     int    `json:"bookmark_count,omitempty"`
										Bookmarked        bool   `json:"bookmarked,omitempty"`
										CreatedAt         string `json:"created_at,omitempty"`
										ConversationIDStr string `json:"conversation_id_str,omitempty"`
										DisplayTextRange  []int  `json:"display_text_range,omitempty"`
										Entities          struct {
											Hashtags []any `json:"hashtags,omitempty"`
											Media    []struct {
												DisplayURL           string `json:"display_url,omitempty"`
												ExpandedURL          string `json:"expanded_url,omitempty"`
												IDStr                string `json:"id_str,omitempty"`
												Indices              []int  `json:"indices,omitempty"`
												MediaKey             string `json:"media_key,omitempty"`
												MediaURLHTTPS        string `json:"media_url_https,omitempty"`
												Type                 string `json:"type,omitempty"`
												URL                  string `json:"url,omitempty"`
												ExtMediaAvailability struct {
													Status string `json:"status,omitempty"`
												} `json:"ext_media_availability,omitempty"`
												Features struct {
													Large struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"large,omitempty"`
													Medium struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"medium,omitempty"`
													Small struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"small,omitempty"`
													Orig struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"orig,omitempty"`
												} `json:"features,omitempty"`
												Sizes struct {
													Large struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"large,omitempty"`
													Medium struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"medium,omitempty"`
													Small struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"small,omitempty"`
													Thumb struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"thumb,omitempty"`
												} `json:"sizes,omitempty"`
												OriginalInfo struct {
													Height     int `json:"height,omitempty"`
													Width      int `json:"width,omitempty"`
													FocusRects []struct {
														X int `json:"x,omitempty"`
														Y int `json:"y,omitempty"`
														W int `json:"w,omitempty"`
														H int `json:"h,omitempty"`
													} `json:"focus_rects,omitempty"`
												} `json:"original_info,omitempty"`
												AllowDownloadStatus struct {
													AllowDownload bool `json:"allow_download,omitempty"`
												} `json:"allow_download_status,omitempty"`
												MediaResults struct {
													Result struct {
														MediaKey string `json:"media_key,omitempty"`
													} `json:"result,omitempty"`
												} `json:"media_results,omitempty"`
											} `json:"media,omitempty"`
											Symbols []struct {
												Indices []int  `json:"indices,omitempty"`
												Text    string `json:"text,omitempty"`
											} `json:"symbols,omitempty"`
											Timestamps []any `json:"timestamps,omitempty"`
											Urls       []struct {
												DisplayURL  string `json:"display_url,omitempty"`
												ExpandedURL string `json:"expanded_url,omitempty"`
												URL         string `json:"url,omitempty"`
												Indices     []int  `json:"indices,omitempty"`
											} `json:"urls,omitempty"`
											UserMentions []struct {
												IDStr      string `json:"id_str,omitempty"`
												Name       string `json:"name,omitempty"`
												ScreenName string `json:"screen_name,omitempty"`
												Indices    []int  `json:"indices,omitempty"`
											} `json:"user_mentions,omitempty"`
										} `json:"entities,omitempty"`
										ExtendedEntities struct {
											Media []struct {
												DisplayURL           string `json:"display_url,omitempty"`
												ExpandedURL          string `json:"expanded_url,omitempty"`
												IDStr                string `json:"id_str,omitempty"`
												Indices              []int  `json:"indices,omitempty"`
												MediaKey             string `json:"media_key,omitempty"`
												MediaURLHTTPS        string `json:"media_url_https,omitempty"`
												Type                 string `json:"type,omitempty"`
												URL                  string `json:"url,omitempty"`
												ExtMediaAvailability struct {
													Status string `json:"status,omitempty"`
												} `json:"ext_media_availability,omitempty"`
												Features struct {
													Large struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"large,omitempty"`
													Medium struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"medium,omitempty"`
													Small struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"small,omitempty"`
													Orig struct {
														Faces []any `json:"faces,omitempty"`
													} `json:"orig,omitempty"`
												} `json:"features,omitempty"`
												Sizes struct {
													Large struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"large,omitempty"`
													Medium struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"medium,omitempty"`
													Small struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"small,omitempty"`
													Thumb struct {
														H      int    `json:"h,omitempty"`
														W      int    `json:"w,omitempty"`
														Resize string `json:"resize,omitempty"`
													} `json:"thumb,omitempty"`
												} `json:"sizes,omitempty"`
												OriginalInfo struct {
													Height     int `json:"height,omitempty"`
													Width      int `json:"width,omitempty"`
													FocusRects []struct {
														X int `json:"x,omitempty"`
														Y int `json:"y,omitempty"`
														W int `json:"w,omitempty"`
														H int `json:"h,omitempty"`
													} `json:"focus_rects,omitempty"`
												} `json:"original_info,omitempty"`
												AllowDownloadStatus struct {
													AllowDownload bool `json:"allow_download,omitempty"`
												} `json:"allow_download_status,omitempty"`
												MediaResults struct {
													Result struct {
														MediaKey string `json:"media_key,omitempty"`
													} `json:"result,omitempty"`
												} `json:"media_results,omitempty"`
											} `json:"media,omitempty"`
										} `json:"extended_entities,omitempty"`
										FavoriteCount             int    `json:"favorite_count,omitempty"`
										Favorited                 bool   `json:"favorited,omitempty"`
										FullText                  string `json:"full_text,omitempty"`
										IsQuoteStatus             bool   `json:"is_quote_status,omitempty"`
										Lang                      string `json:"lang,omitempty"`
										PossiblySensitive         bool   `json:"possibly_sensitive,omitempty"`
										PossiblySensitiveEditable bool   `json:"possibly_sensitive_editable,omitempty"`
										QuoteCount                int    `json:"quote_count,omitempty"`
										ReplyCount                int    `json:"reply_count,omitempty"`
										RetweetCount              int    `json:"retweet_count,omitempty"`
										Retweeted                 bool   `json:"retweeted,omitempty"`
										UserIDStr                 string `json:"user_id_str,omitempty"`
										IDStr                     string `json:"id_str,omitempty"`
										InReplyToScreenName       string `json:"in_reply_to_screen_name,omitempty"`
										InReplyToStatusIDStr      string `json:"in_reply_to_status_id_str,omitempty"`
									} `json:"legacy,omitempty"`
									QuickPromoteEligibility struct {
										Eligibility string `json:"eligibility,omitempty"`
									} `json:"quick_promote_eligibility,omitempty"`
								} `json:"result,omitempty"`
							} `json:"tweet_results,omitempty"`
							TweetDisplayType    string `json:"tweetDisplayType,omitempty"`
							HasModeratedReplies bool   `json:"hasModeratedReplies,omitempty"`
						} `json:"itemContent,omitempty"`
						Items []struct {
							EntryID string `json:"entryId,omitempty"`
							Item    struct {
								ItemContent struct {
									ItemType     string `json:"itemType,omitempty"`
									Typename     string `json:"__typename,omitempty"`
									TweetResults struct {
										Result struct {
											Typename          string `json:"__typename,omitempty"`
											RestID            string `json:"rest_id,omitempty"`
											HasBirdwatchNotes bool   `json:"has_birdwatch_notes,omitempty"`
											Core              struct {
												UserResults struct {
													Result struct {
														Typename                   string `json:"__typename,omitempty"`
														ID                         string `json:"id,omitempty"`
														RestID                     string `json:"rest_id,omitempty"`
														AffiliatesHighlightedLabel struct {
														} `json:"affiliates_highlighted_label,omitempty"`
														HasGraduatedAccess bool   `json:"has_graduated_access,omitempty"`
														IsBlueVerified     bool   `json:"is_blue_verified,omitempty"`
														ProfileImageShape  string `json:"profile_image_shape,omitempty"`
														Legacy             struct {
															FollowedBy          bool   `json:"followed_by,omitempty"`
															Following           bool   `json:"following,omitempty"`
															CanDm               bool   `json:"can_dm,omitempty"`
															CanMediaTag         bool   `json:"can_media_tag,omitempty"`
															CreatedAt           string `json:"created_at,omitempty"`
															DefaultProfile      bool   `json:"default_profile,omitempty"`
															DefaultProfileImage bool   `json:"default_profile_image,omitempty"`
															Description         string `json:"description,omitempty"`
															Entities            struct {
																Description struct {
																	Urls []struct {
																		DisplayURL  string `json:"display_url,omitempty"`
																		ExpandedURL string `json:"expanded_url,omitempty"`
																		URL         string `json:"url,omitempty"`
																		Indices     []int  `json:"indices,omitempty"`
																	} `json:"urls,omitempty"`
																} `json:"description,omitempty"`
																URL struct {
																	Urls []struct {
																		DisplayURL  string `json:"display_url,omitempty"`
																		ExpandedURL string `json:"expanded_url,omitempty"`
																		URL         string `json:"url,omitempty"`
																		Indices     []int  `json:"indices,omitempty"`
																	} `json:"urls,omitempty"`
																} `json:"url,omitempty"`
															} `json:"entities,omitempty"`
															FastFollowersCount      int      `json:"fast_followers_count,omitempty"`
															FavouritesCount         int      `json:"favourites_count,omitempty"`
															FollowersCount          int      `json:"followers_count,omitempty"`
															FriendsCount            int      `json:"friends_count,omitempty"`
															HasCustomTimelines      bool     `json:"has_custom_timelines,omitempty"`
															IsTranslator            bool     `json:"is_translator,omitempty"`
															ListedCount             int      `json:"listed_count,omitempty"`
															Location                string   `json:"location,omitempty"`
															MediaCount              int      `json:"media_count,omitempty"`
															Name                    string   `json:"name,omitempty"`
															NormalFollowersCount    int      `json:"normal_followers_count,omitempty"`
															PinnedTweetIdsStr       []string `json:"pinned_tweet_ids_str,omitempty"`
															PossiblySensitive       bool     `json:"possibly_sensitive,omitempty"`
															ProfileBannerURL        string   `json:"profile_banner_url,omitempty"`
															ProfileImageURLHTTPS    string   `json:"profile_image_url_https,omitempty"`
															ProfileInterstitialType string   `json:"profile_interstitial_type,omitempty"`
															ScreenName              string   `json:"screen_name,omitempty"`
															StatusesCount           int      `json:"statuses_count,omitempty"`
															TranslatorType          string   `json:"translator_type,omitempty"`
															URL                     string   `json:"url,omitempty"`
															Verified                bool     `json:"verified,omitempty"`
															WantRetweets            bool     `json:"want_retweets,omitempty"`
															WithheldInCountries     []any    `json:"withheld_in_countries,omitempty"`
														} `json:"legacy,omitempty"`
														Professional struct {
															RestID           string `json:"rest_id,omitempty"`
															ProfessionalType string `json:"professional_type,omitempty"`
															Category         []struct {
																ID       int    `json:"id,omitempty"`
																Name     string `json:"name,omitempty"`
																IconName string `json:"icon_name,omitempty"`
															} `json:"category,omitempty"`
														} `json:"professional,omitempty"`
														TipjarSettings struct {
															IsEnabled      bool   `json:"is_enabled,omitempty"`
															EthereumHandle string `json:"ethereum_handle,omitempty"`
														} `json:"tipjar_settings,omitempty"`
														SuperFollowEligible bool `json:"super_follow_eligible,omitempty"`
													} `json:"result,omitempty"`
												} `json:"user_results,omitempty"`
											} `json:"core,omitempty"`
											UnmentionData struct {
											} `json:"unmention_data,omitempty"`
											EditControl struct {
												EditTweetIds       []string `json:"edit_tweet_ids,omitempty"`
												EditableUntilMsecs string   `json:"editable_until_msecs,omitempty"`
												IsEditEligible     bool     `json:"is_edit_eligible,omitempty"`
												EditsRemaining     string   `json:"edits_remaining,omitempty"`
											} `json:"edit_control,omitempty"`
											IsTranslatable bool `json:"is_translatable,omitempty"`
											Views          struct {
												Count string `json:"count,omitempty"`
												State string `json:"state,omitempty"`
											} `json:"views,omitempty"`
											Source string `json:"source,omitempty"`
											Legacy struct {
												BookmarkCount     int    `json:"bookmark_count,omitempty"`
												Bookmarked        bool   `json:"bookmarked,omitempty"`
												CreatedAt         string `json:"created_at,omitempty"`
												ConversationIDStr string `json:"conversation_id_str,omitempty"`
												DisplayTextRange  []int  `json:"display_text_range,omitempty"`
												Entities          struct {
													Media []struct {
														DisplayURL           string `json:"display_url"`
														ExpandedURL          string `json:"expanded_url"`
														IDStr                string `json:"id_str"`
														Indices              []int  `json:"indices"`
														MediaKey             string `json:"media_key"`
														MediaURLHTTPS        string `json:"media_url_https"`
														Type                 string `json:"type"`
														URL                  string `json:"url"`
														ExtMediaAvailability struct {
															Status string `json:"status"`
														} `json:"ext_media_availability"`
														Features struct {
															Large struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"large"`
															Medium struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"medium"`
															Small struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"small"`
															Orig struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"orig"`
														} `json:"features"`
														Sizes struct {
															Large struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"large"`
															Medium struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"medium"`
															Small struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"small"`
															Thumb struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"thumb"`
														} `json:"sizes"`
														OriginalInfo struct {
															Height     int `json:"height"`
															Width      int `json:"width"`
															FocusRects []struct {
																X int `json:"x"`
																Y int `json:"y"`
																W int `json:"w"`
																H int `json:"h"`
															} `json:"focus_rects"`
														} `json:"original_info"`
														AllowDownloadStatus struct {
															AllowDownload bool `json:"allow_download"`
														} `json:"allow_download_status"`
														MediaResults struct {
															Result struct {
																MediaKey string `json:"media_key"`
															} `json:"result"`
														} `json:"media_results"`
													} `json:"media"`
													Hashtags   []any `json:"hashtags,omitempty"`
													Symbols    []any `json:"symbols,omitempty"`
													Timestamps []any `json:"timestamps,omitempty"`
													Urls       []struct {
														DisplayURL  string `json:"display_url,omitempty"`
														ExpandedURL string `json:"expanded_url,omitempty"`
														URL         string `json:"url,omitempty"`
														Indices     []int  `json:"indices,omitempty"`
													} `json:"urls,omitempty"`
													UserMentions []struct {
														IDStr      string `json:"id_str,omitempty"`
														Name       string `json:"name,omitempty"`
														ScreenName string `json:"screen_name,omitempty"`
														Indices    []int  `json:"indices,omitempty"`
													} `json:"user_mentions,omitempty"`
												} `json:"entities,omitempty"`
												FavoriteCount             int    `json:"favorite_count,omitempty"`
												Favorited                 bool   `json:"favorited,omitempty"`
												FullText                  string `json:"full_text,omitempty"`
												InReplyToScreenName       string `json:"in_reply_to_screen_name,omitempty"`
												InReplyToStatusIDStr      string `json:"in_reply_to_status_id_str,omitempty"`
												InReplyToUserIDStr        string `json:"in_reply_to_user_id_str,omitempty"`
												IsQuoteStatus             bool   `json:"is_quote_status,omitempty"`
												Lang                      string `json:"lang,omitempty"`
												PossiblySensitive         bool   `json:"possibly_sensitive,omitempty"`
												PossiblySensitiveEditable bool   `json:"possibly_sensitive_editable,omitempty"`
												QuoteCount                int    `json:"quote_count,omitempty"`
												ReplyCount                int    `json:"reply_count,omitempty"`
												RetweetCount              int    `json:"retweet_count,omitempty"`
												Retweeted                 bool   `json:"retweeted,omitempty"`
												UserIDStr                 string `json:"user_id_str,omitempty"`
												IDStr                     string `json:"id_str,omitempty"`
											} `json:"legacy,omitempty"`
											QuickPromoteEligibility struct {
												Eligibility string `json:"eligibility,omitempty"`
											} `json:"quick_promote_eligibility,omitempty"`
										} `json:"result,omitempty"`
									} `json:"tweet_results,omitempty"`
									TweetDisplayType string `json:"tweetDisplayType,omitempty"`
								} `json:"itemContent,omitempty"`
								ClientEventInfo struct {
									Details struct {
										ConversationDetails struct {
											ConversationSection string `json:"conversationSection,omitempty"`
										} `json:"conversationDetails,omitempty"`
										TimelinesDetails struct {
											ControllerData string `json:"controllerData,omitempty"`
										} `json:"timelinesDetails,omitempty"`
									} `json:"details,omitempty"`
								} `json:"clientEventInfo,omitempty"`
							} `json:"item,omitempty"`
						}
					} `json:"content,omitempty"`
				} `json:"entries,omitempty"`
				Direction string `json:"direction,omitempty"`
			} `json:"instructions,omitempty"`
			Metadata struct {
				ReaderModeConfig struct {
					IsReaderModeAvailable bool `json:"is_reader_mode_available,omitempty"`
				} `json:"reader_mode_config,omitempty"`
			} `json:"metadata,omitempty"`
		} `json:"threaded_conversation_with_injections_v2,omitempty"`
	} `json:"data,omitempty"`
}

type SearchTimelineResponse struct {
	Data struct {
		SearchByRawQuery struct {
			SearchTimeline struct {
				Timeline struct {
					Instructions []struct {
						Type    string `json:"type"`
						Entries []struct {
							EntryID   string `json:"entryId"`
							SortIndex string `json:"sortIndex"`
							Content   struct {
								EntryType   string `json:"entryType"`
								Typename    string `json:"__typename"`
								ItemContent struct {
									ItemType     string `json:"itemType"`
									Typename     string `json:"__typename"`
									TweetResults struct {
										Result struct {
											Typename string `json:"__typename"`
											RestID   string `json:"rest_id"`
											Core     struct {
												UserResults struct {
													Result struct {
														Typename                   string `json:"__typename"`
														ID                         string `json:"id"`
														RestID                     string `json:"rest_id"`
														AffiliatesHighlightedLabel struct {
														} `json:"affiliates_highlighted_label"`
														HasGraduatedAccess bool   `json:"has_graduated_access"`
														IsBlueVerified     bool   `json:"is_blue_verified"`
														ProfileImageShape  string `json:"profile_image_shape"`
														Legacy             struct {
															FollowedBy          bool   `json:"followed_by"`
															Following           bool   `json:"following"`
															CanDm               bool   `json:"can_dm"`
															CanMediaTag         bool   `json:"can_media_tag"`
															CreatedAt           string `json:"created_at"`
															DefaultProfile      bool   `json:"default_profile"`
															DefaultProfileImage bool   `json:"default_profile_image"`
															Description         string `json:"description"`
															Entities            struct {
																Description struct {
																	Urls []struct {
																		DisplayURL  string `json:"display_url"`
																		ExpandedURL string `json:"expanded_url"`
																		URL         string `json:"url"`
																		Indices     []int  `json:"indices"`
																	} `json:"urls"`
																} `json:"description"`
																URL struct {
																	Urls []struct {
																		DisplayURL  string `json:"display_url"`
																		ExpandedURL string `json:"expanded_url"`
																		URL         string `json:"url"`
																		Indices     []int  `json:"indices"`
																	} `json:"urls"`
																} `json:"url"`
															} `json:"entities"`
															FastFollowersCount      int           `json:"fast_followers_count"`
															FavouritesCount         int           `json:"favourites_count"`
															FollowersCount          int           `json:"followers_count"`
															FriendsCount            int           `json:"friends_count"`
															HasCustomTimelines      bool          `json:"has_custom_timelines"`
															IsTranslator            bool          `json:"is_translator"`
															ListedCount             int           `json:"listed_count"`
															Location                string        `json:"location"`
															MediaCount              int           `json:"media_count"`
															Name                    string        `json:"name"`
															NormalFollowersCount    int           `json:"normal_followers_count"`
															PinnedTweetIdsStr       []string      `json:"pinned_tweet_ids_str"`
															PossiblySensitive       bool          `json:"possibly_sensitive"`
															ProfileBannerURL        string        `json:"profile_banner_url"`
															ProfileImageURLHTTPS    string        `json:"profile_image_url_https"`
															ProfileInterstitialType string        `json:"profile_interstitial_type"`
															ScreenName              string        `json:"screen_name"`
															StatusesCount           int           `json:"statuses_count"`
															TranslatorType          string        `json:"translator_type"`
															URL                     string        `json:"url"`
															Verified                bool          `json:"verified"`
															WantRetweets            bool          `json:"want_retweets"`
															WithheldInCountries     []interface{} `json:"withheld_in_countries"`
														} `json:"legacy"`
														Professional struct {
															RestID           string `json:"rest_id"`
															ProfessionalType string `json:"professional_type"`
															Category         []struct {
																ID       int    `json:"id"`
																Name     string `json:"name"`
																IconName string `json:"icon_name"`
															} `json:"category"`
														} `json:"professional"`
														TipjarSettings struct {
															IsEnabled      bool   `json:"is_enabled"`
															EthereumHandle string `json:"ethereum_handle"`
														} `json:"tipjar_settings"`
														SuperFollowEligible bool `json:"super_follow_eligible"`
													} `json:"result"`
												} `json:"user_results"`
											} `json:"core"`
											UnmentionData struct {
											} `json:"unmention_data"`
											EditControl struct {
												EditTweetIds       []string `json:"edit_tweet_ids"`
												EditableUntilMsecs string   `json:"editable_until_msecs"`
												IsEditEligible     bool     `json:"is_edit_eligible"`
												EditsRemaining     string   `json:"edits_remaining"`
											} `json:"edit_control"`
											IsTranslatable bool `json:"is_translatable"`
											Views          struct {
												Count string `json:"count"`
												State string `json:"state"`
											} `json:"views"`
											Source string `json:"source"`
											Legacy struct {
												BookmarkCount     int    `json:"bookmark_count"`
												Bookmarked        bool   `json:"bookmarked"`
												CreatedAt         string `json:"created_at"`
												ConversationIDStr string `json:"conversation_id_str"`
												DisplayTextRange  []int  `json:"display_text_range"`
												Entities          struct {
													Hashtags []interface{} `json:"hashtags"`
													Media    []struct {
														DisplayURL           string `json:"display_url"`
														ExpandedURL          string `json:"expanded_url"`
														IDStr                string `json:"id_str"`
														Indices              []int  `json:"indices"`
														MediaKey             string `json:"media_key"`
														MediaURLHTTPS        string `json:"media_url_https"`
														Type                 string `json:"type"`
														URL                  string `json:"url"`
														ExtMediaAvailability struct {
															Status string `json:"status"`
														} `json:"ext_media_availability"`
														Features struct {
															Large struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"large"`
															Medium struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"medium"`
															Small struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"small"`
															Orig struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"orig"`
														} `json:"features"`
														Sizes struct {
															Large struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"large"`
															Medium struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"medium"`
															Small struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"small"`
															Thumb struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"thumb"`
														} `json:"sizes"`
														OriginalInfo struct {
															Height     int `json:"height"`
															Width      int `json:"width"`
															FocusRects []struct {
																X int `json:"x"`
																Y int `json:"y"`
																W int `json:"w"`
																H int `json:"h"`
															} `json:"focus_rects"`
														} `json:"original_info"`
														AllowDownloadStatus struct {
															AllowDownload bool `json:"allow_download"`
														} `json:"allow_download_status"`
														MediaResults struct {
															Result struct {
																MediaKey string `json:"media_key"`
															} `json:"result"`
														} `json:"media_results"`
													} `json:"media"`
													Symbols      []interface{} `json:"symbols"`
													Timestamps   []interface{} `json:"timestamps"`
													Urls         []interface{} `json:"urls"`
													UserMentions []struct {
														IDStr      string `json:"id_str"`
														Name       string `json:"name"`
														ScreenName string `json:"screen_name"`
														Indices    []int  `json:"indices"`
													} `json:"user_mentions"`
												} `json:"entities"`
												ExtendedEntities struct {
													Media []struct {
														DisplayURL           string `json:"display_url"`
														ExpandedURL          string `json:"expanded_url"`
														IDStr                string `json:"id_str"`
														Indices              []int  `json:"indices"`
														MediaKey             string `json:"media_key"`
														MediaURLHTTPS        string `json:"media_url_https"`
														Type                 string `json:"type"`
														URL                  string `json:"url"`
														ExtMediaAvailability struct {
															Status string `json:"status"`
														} `json:"ext_media_availability"`
														Features struct {
															Large struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"large"`
															Medium struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"medium"`
															Small struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"small"`
															Orig struct {
																Faces []struct {
																	X int `json:"x"`
																	Y int `json:"y"`
																	H int `json:"h"`
																	W int `json:"w"`
																} `json:"faces"`
															} `json:"orig"`
														} `json:"features"`
														Sizes struct {
															Large struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"large"`
															Medium struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"medium"`
															Small struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"small"`
															Thumb struct {
																H      int    `json:"h"`
																W      int    `json:"w"`
																Resize string `json:"resize"`
															} `json:"thumb"`
														} `json:"sizes"`
														OriginalInfo struct {
															Height     int `json:"height"`
															Width      int `json:"width"`
															FocusRects []struct {
																X int `json:"x"`
																Y int `json:"y"`
																W int `json:"w"`
																H int `json:"h"`
															} `json:"focus_rects"`
														} `json:"original_info"`
														AllowDownloadStatus struct {
															AllowDownload bool `json:"allow_download"`
														} `json:"allow_download_status"`
														MediaResults struct {
															Result struct {
																MediaKey string `json:"media_key"`
															} `json:"result"`
														} `json:"media_results"`
													} `json:"media"`
												} `json:"extended_entities"`
												FavoriteCount             int    `json:"favorite_count"`
												Favorited                 bool   `json:"favorited"`
												FullText                  string `json:"full_text"`
												InReplyToScreenName       string `json:"in_reply_to_screen_name,omitempty"`
												InReplyToStatusIDStr      string `json:"in_reply_to_status_id_str,omitempty"`
												IsQuoteStatus             bool   `json:"is_quote_status"`
												Lang                      string `json:"lang"`
												PossiblySensitive         bool   `json:"possibly_sensitive"`
												PossiblySensitiveEditable bool   `json:"possibly_sensitive_editable"`
												QuoteCount                int    `json:"quote_count"`
												ReplyCount                int    `json:"reply_count"`
												RetweetCount              int    `json:"retweet_count"`
												Retweeted                 bool   `json:"retweeted"`
												UserIDStr                 string `json:"user_id_str"`
												IDStr                     string `json:"id_str"`
											} `json:"legacy"`
										} `json:"result"`
									} `json:"tweet_results"`
									TweetDisplayType string `json:"tweetDisplayType"`
								} `json:"itemContent"`
								ClientEventInfo struct {
									Component string `json:"component"`
									Element   string `json:"element"`
									Details   struct {
										TimelinesDetails struct {
											ControllerData string `json:"controllerData"`
										} `json:"timelinesDetails"`
									} `json:"details"`
								} `json:"clientEventInfo"`
							} `json:"content"`
						} `json:"entries"`
					} `json:"instructions"`
				} `json:"timeline"`
			} `json:"search_timeline"`
		} `json:"search_by_raw_query"`
	} `json:"data"`
}

// TweetOptions contains optional parameters for creating a tweet
type TweetOptions struct {
	ReplyToTweetID string // ID of tweet to reply to, if this is a reply
}

type CreateTweetRequest struct {
	Variables CreateTweetVariables `json:"variables,omitempty"`
	Features  CreateTweetFeatures  `json:"features,omitempty"`
	QueryID   string               `json:"queryId,omitempty"`
}

type CreateTweetVariables struct {
	TweetText string `json:"tweet_text,omitempty"`
	Reply     *struct {
		InReplyToTweetID    string `json:"in_reply_to_tweet_id,omitempty"`
		ExcludeReplyUserIds []any  `json:"exclude_reply_user_ids,omitempty"`
	} `json:"reply,omitempty"`
	DarkRequest bool `json:"dark_request,omitempty"`
	Media       struct {
		MediaEntities     []any `json:"media_entities,omitempty"`
		PossiblySensitive bool  `json:"possibly_sensitive,omitempty"`
	} `json:"media,omitempty"`
	SemanticAnnotationIds  []any `json:"semantic_annotation_ids,omitempty"`
	DisallowedReplyOptions any   `json:"disallowed_reply_options,omitempty"`
}

type CreateTweetFeatures struct {
	CommunitiesWebEnableTweetCommunityResultsFetch                 bool `json:"communities_web_enable_tweet_community_results_fetch"`
	C9STweetAnatomyModeratorBadgeEnabled                           bool `json:"c9s_tweet_anatomy_moderator_badge_enabled"`
	ResponsiveWebEditTweetAPIEnabled                               bool `json:"responsive_web_edit_tweet_api_enabled"`
	GraphqlIsTranslatableRwebTweetIsTranslatableEnabled            bool `json:"graphql_is_translatable_rweb_tweet_is_translatable_enabled"`
	ViewCountsEverywhereAPIEnabled                                 bool `json:"view_counts_everywhere_api_enabled"`
	LongformNotetweetsConsumptionEnabled                           bool `json:"longform_notetweets_consumption_enabled"`
	ResponsiveWebTwitterArticleTweetConsumptionEnabled             bool `json:"responsive_web_twitter_article_tweet_consumption_enabled"`
	TweetAwardsWebTippingEnabled                                   bool `json:"tweet_awards_web_tipping_enabled"`
	CreatorSubscriptionsQuoteTweetPreviewEnabled                   bool `json:"creator_subscriptions_quote_tweet_preview_enabled"`
	LongformNotetweetsRichTextReadEnabled                          bool `json:"longform_notetweets_rich_text_read_enabled"`
	LongformNotetweetsInlineMediaEnabled                           bool `json:"longform_notetweets_inline_media_enabled"`
	ArticlesPreviewEnabled                                         bool `json:"articles_preview_enabled"`
	RwebVideoTimestampsEnabled                                     bool `json:"rweb_video_timestamps_enabled"`
	ProfileLabelImprovementsPcfLabelInPostEnabled                  bool `json:"profile_label_improvements_pcf_label_in_post_enabled"`
	RwebTipjarConsumptionEnabled                                   bool `json:"rweb_tipjar_consumption_enabled"`
	ResponsiveWebGraphqlExcludeDirectiveEnabled                    bool `json:"responsive_web_graphql_exclude_directive_enabled"`
	VerifiedPhoneLabelEnabled                                      bool `json:"verified_phone_label_enabled"`
	FreedomOfSpeechNotReachFetchEnabled                            bool `json:"freedom_of_speech_not_reach_fetch_enabled"`
	StandardizedNudgesMisinfo                                      bool `json:"standardized_nudges_misinfo"`
	TweetWithVisibilityResultsPreferGqlLimitedActionsPolicyEnabled bool `json:"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled"`
	ResponsiveWebGraphqlSkipUserProfileImageExtensionsEnabled      bool `json:"responsive_web_graphql_skip_user_profile_image_extensions_enabled"`
	ResponsiveWebGraphqlTimelineNavigationEnabled                  bool `json:"responsive_web_graphql_timeline_navigation_enabled"`
	ResponsiveWebEnhanceCardsEnabled                               bool `json:"responsive_web_enhance_cards_enabled"`
}

type CreateTweetResponse struct {
	Data struct {
		CreateTweet struct {
			TweetResults struct {
				Result struct {
					RestID string `json:"rest_id"`
					Core   struct {
						UserResults struct {
							Result struct {
								Typename                   string `json:"__typename"`
								ID                         string `json:"id"`
								RestID                     string `json:"rest_id"`
								AffiliatesHighlightedLabel struct {
								} `json:"affiliates_highlighted_label"`
								HasGraduatedAccess bool   `json:"has_graduated_access"`
								IsBlueVerified     bool   `json:"is_blue_verified"`
								ProfileImageShape  string `json:"profile_image_shape"`
								Legacy             struct {
									Following           bool   `json:"following"`
									CanDm               bool   `json:"can_dm"`
									CanMediaTag         bool   `json:"can_media_tag"`
									CreatedAt           string `json:"created_at"`
									DefaultProfile      bool   `json:"default_profile"`
									DefaultProfileImage bool   `json:"default_profile_image"`
									Description         string `json:"description"`
									Entities            struct {
										Description struct {
											Urls []any `json:"urls"`
										} `json:"description"`
									} `json:"entities"`
									FastFollowersCount      int      `json:"fast_followers_count"`
									FavouritesCount         int      `json:"favourites_count"`
									FollowersCount          int      `json:"followers_count"`
									FriendsCount            int      `json:"friends_count"`
									HasCustomTimelines      bool     `json:"has_custom_timelines"`
									IsTranslator            bool     `json:"is_translator"`
									ListedCount             int      `json:"listed_count"`
									Location                string   `json:"location"`
									MediaCount              int      `json:"media_count"`
									Name                    string   `json:"name"`
									NeedsPhoneVerification  bool     `json:"needs_phone_verification"`
									NormalFollowersCount    int      `json:"normal_followers_count"`
									PinnedTweetIdsStr       []string `json:"pinned_tweet_ids_str"`
									PossiblySensitive       bool     `json:"possibly_sensitive"`
									ProfileImageURLHTTPS    string   `json:"profile_image_url_https"`
									ProfileInterstitialType string   `json:"profile_interstitial_type"`
									ScreenName              string   `json:"screen_name"`
									StatusesCount           int      `json:"statuses_count"`
									TranslatorType          string   `json:"translator_type"`
									Verified                bool     `json:"verified"`
									WantRetweets            bool     `json:"want_retweets"`
									WithheldInCountries     []any    `json:"withheld_in_countries"`
								} `json:"legacy"`
								TipjarSettings struct {
									IsEnabled     bool   `json:"is_enabled"`
									BitcoinHandle string `json:"bitcoin_handle"`
								} `json:"tipjar_settings"`
							} `json:"result"`
						} `json:"user_results"`
					} `json:"core"`
					UnmentionData struct {
					} `json:"unmention_data"`
					EditControl struct {
						EditTweetIds       []string `json:"edit_tweet_ids"`
						EditableUntilMsecs string   `json:"editable_until_msecs"`
						IsEditEligible     bool     `json:"is_edit_eligible"`
						EditsRemaining     string   `json:"edits_remaining"`
					} `json:"edit_control"`
					IsTranslatable bool `json:"is_translatable"`
					Views          struct {
						State string `json:"state"`
					} `json:"views"`
					Source string `json:"source"`
					Legacy struct {
						BookmarkCount     int    `json:"bookmark_count"`
						Bookmarked        bool   `json:"bookmarked"`
						CreatedAt         string `json:"created_at"`
						ConversationIDStr string `json:"conversation_id_str"`
						DisplayTextRange  []int  `json:"display_text_range"`
						Entities          struct {
							Hashtags     []any `json:"hashtags"`
							Symbols      []any `json:"symbols"`
							Timestamps   []any `json:"timestamps"`
							Urls         []any `json:"urls"`
							UserMentions []struct {
								IDStr      string `json:"id_str"`
								Name       string `json:"name"`
								ScreenName string `json:"screen_name"`
								Indices    []int  `json:"indices"`
							} `json:"user_mentions"`
						} `json:"entities"`
						FavoriteCount        int    `json:"favorite_count"`
						Favorited            bool   `json:"favorited"`
						FullText             string `json:"full_text"`
						InReplyToScreenName  string `json:"in_reply_to_screen_name"`
						InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
						InReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`
						IsQuoteStatus        bool   `json:"is_quote_status"`
						Lang                 string `json:"lang"`
						QuoteCount           int    `json:"quote_count"`
						ReplyCount           int    `json:"reply_count"`
						RetweetCount         int    `json:"retweet_count"`
						Retweeted            bool   `json:"retweeted"`
						UserIDStr            string `json:"user_id_str"`
						IDStr                string `json:"id_str"`
					} `json:"legacy"`
					UnmentionInfo struct {
					} `json:"unmention_info"`
				} `json:"result"`
			} `json:"tweet_results"`
		} `json:"create_tweet"`
	} `json:"data"`
}

type FavoriteTweetRequest struct {
	Variables struct {
		TweetID string `json:"tweet_id,omitempty"`
	} `json:"variables,omitempty"`
	QueryID string `json:"queryId,omitempty"`
}

type FavoriteTweetResponse struct {
	Data struct {
		FavoriteTweet string `json:"favorite_tweet,omitempty"`
	} `json:"data,omitempty"`
}

type GetUserDetailsResponse struct {
	Data struct {
		User struct {
			Result struct {
				Typename                   string `json:"__typename"`
				ID                         string `json:"id"`
				RestID                     string `json:"rest_id"`
				AffiliatesHighlightedLabel struct {
				} `json:"affiliates_highlighted_label"`
				HasGraduatedAccess       bool   `json:"has_graduated_access"`
				ParodyCommentaryFanLabel string `json:"parody_commentary_fan_label"`
				IsBlueVerified           bool   `json:"is_blue_verified"`
				ProfileImageShape        string `json:"profile_image_shape"`
				Legacy                   struct {
					Following           bool   `json:"following"`
					CanDm               bool   `json:"can_dm"`
					CanMediaTag         bool   `json:"can_media_tag"`
					CreatedAt           string `json:"created_at"`
					DefaultProfile      bool   `json:"default_profile"`
					DefaultProfileImage bool   `json:"default_profile_image"`
					Description         string `json:"description"`
					Entities            struct {
						Description struct {
							Urls []any `json:"urls"`
						} `json:"description"`
					} `json:"entities"`
					FastFollowersCount      int    `json:"fast_followers_count"`
					FavouritesCount         int    `json:"favourites_count"`
					FollowersCount          int    `json:"followers_count"`
					FriendsCount            int    `json:"friends_count"`
					HasCustomTimelines      bool   `json:"has_custom_timelines"`
					IsTranslator            bool   `json:"is_translator"`
					ListedCount             int    `json:"listed_count"`
					Location                string `json:"location"`
					MediaCount              int    `json:"media_count"`
					Name                    string `json:"name"`
					NeedsPhoneVerification  bool   `json:"needs_phone_verification"`
					NormalFollowersCount    int    `json:"normal_followers_count"`
					PinnedTweetIdsStr       []any  `json:"pinned_tweet_ids_str"`
					PossiblySensitive       bool   `json:"possibly_sensitive"`
					ProfileBannerURL        string `json:"profile_banner_url"`
					ProfileImageURLHTTPS    string `json:"profile_image_url_https"`
					ProfileInterstitialType string `json:"profile_interstitial_type"`
					ScreenName              string `json:"screen_name"`
					StatusesCount           int    `json:"statuses_count"`
					TranslatorType          string `json:"translator_type"`
					Verified                bool   `json:"verified"`
					WantRetweets            bool   `json:"want_retweets"`
					WithheldInCountries     []any  `json:"withheld_in_countries"`
				} `json:"legacy"`
				TipjarSettings struct {
				} `json:"tipjar_settings"`
				LegacyExtendedProfile struct {
					Birthdate struct {
						Day            int    `json:"day"`
						Month          int    `json:"month"`
						Year           int    `json:"year"`
						Visibility     string `json:"visibility"`
						YearVisibility string `json:"year_visibility"`
					} `json:"birthdate"`
				} `json:"legacy_extended_profile"`
				IsProfileTranslatable           bool `json:"is_profile_translatable"`
				HasHiddenSubscriptionsOnProfile bool `json:"has_hidden_subscriptions_on_profile"`
				VerificationInfo                struct {
					IsIdentityVerified bool `json:"is_identity_verified"`
					Reason             struct {
						Description struct {
							Text     string `json:"text"`
							Entities []struct {
								FromIndex int `json:"from_index"`
								ToIndex   int `json:"to_index"`
								Ref       struct {
									URL     string `json:"url"`
									URLType string `json:"url_type"`
								} `json:"ref"`
							} `json:"entities"`
						} `json:"description"`
						VerifiedSinceMsec string `json:"verified_since_msec"`
					} `json:"reason"`
				} `json:"verification_info"`
				HighlightsInfo struct {
					CanHighlightTweets bool   `json:"can_highlight_tweets"`
					HighlightedTweets  string `json:"highlighted_tweets"`
				} `json:"highlights_info"`
				UserSeedTweetCount     int  `json:"user_seed_tweet_count"`
				PremiumGiftingEligible bool `json:"premium_gifting_eligible"`
				BusinessAccount        struct {
				} `json:"business_account"`
				CreatorSubscriptionsCount int `json:"creator_subscriptions_count"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}
