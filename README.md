# Bionic
Bionic is a tool to convert data exports from web apps to a single SQLite database. Bionic currently supports data exports from Google, Apple Health, Spotify, Telegram, RescueTime, Instagram, Twitter and Netflix.

[Schema and demo data](https://bionic-db.vercel.app/db_public).

![Example of bionic usage](https://user-images.githubusercontent.com/6896447/108840862-f9932f80-75e7-11eb-9014-70afc55ff302.png)

**Fun**: you can use Bionic to explore your own data and discover insights about yourself. Join tables between different sources to create reports like "Songs I listen in different locations" or dive deep into a single source to create "How amount of Telegram messages per week with different people changed over time" report.

**Research**: if you research human behaviour, subjective metrics could heavily impact your findings. If your respondents are able to run Bionic and send you aggregated results from their data, you can collect new objectives datasets describing important parts of life: transportation, social media, knowledge work and others.

**Development**: you can use Bionic as a Go package to implement personal data import in your apps.

**Education**: you can include Bionic exercises in your articles, courses or books. Learning to process data on personal records is much more exciting than processing artificial datasets.

## Install

### Homebrew

```bash
brew install bionic-dev/tap/bionic
```

Update:

```bash
brew upgrade bionic-dev/tap/bionic
```

### cURL

```bash
curl -L https://raw.githubusercontent.com/bionic-dev/bionic/main/install.sh | bash -s -- -b /usr/local/bin
```

## Usage

### Import data

Use the following syntax to convert downloaded data to a SQLite database:
```bash
bionic import [provider] [path to downloaded directory or an archive] --db [path to sqlite db]
```

If the database doesn't exist, Bionic will create a new one. If it already exists, Bionic will create tables if needed and append new rows.

Examples:
```bash
bionic import google /Users/seva/gdpr_exports/Takeout/ --db db.sqlite
bionic import health /Users/seva/gdpr_exports/apple-health.zip --db db.sqlite
bionic import spotify /Users/seva/gdpr_exports/MyData/ --db db.sqlite
```

### Generate views

Bionic provides helper tables ("views") to make processing data easier. 

For example, `google_searches` is a view based on original `google_activity` table, 
but filtered only to include search queries and altered to have the search query as a column.

To generate or update views run:
```bash
bionic generate-views --db db.sqlite
```

### Query

You can query the database with  ```sqlite3``` client:
```bash
$ sqlite3 db.sqlite                                               
SQLite version 3.28.0 2019-04-15 14:49:49
Enter ".help" for usage hints.

sqlite> select * from netflix_playback_related_events limit 1;
1|2021-01-22 20:46:21.696934+03:00|2021-01-22 20:46:21.696934+03:00||Seva|How I Met Your Mother: Season 1: "Come On"|Apple iPhone XR|RU|2020-12-30 20:14:21+00:00
```

Alternatively, you can use [datasette](https://github.com/simonw/datasette) to build a web UI to view and query data:

```bash
$ datasette serve db.sqlite
INFO:     Started server process [23975]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://127.0.0.1:8001 (Press CTRL+C to quit)
```

![datasette screenshot](https://user-images.githubusercontent.com/6896447/108776053-eb56fc00-7572-11eb-9081-1732cdc4a3bd.png)

You can also use Python and [pandas](https://pandas.pydata.org/) to process data:
```python
import pandas as pd
import sqlite3

DATABASE_PATH = '/Users/seva/db.sqlite'
db_connection = sqlite3.connect(DATABASE_PATH)

messages_df = pd.read_sql('select * from telegram_messages;', con=db_connection)
```

## Providers


| Name | Export link | Created tables | Notes 
|------|-------------|----------------|------
| Google: `google` | https://takeout.google.com/ |`google_activity`, `google_activity_details`, `google_activity_location_infos`, `google_activity_path_points`, `google_activity_products`, `google_activity_products_assoc`, `google_activity_segments`, `google_activity_subtitles`, `google_activity_type_candidates`, `google_candidate_locations`, `google_location_activity`, `google_location_activity_type_candidates`, `google_location_history`, `google_place_path_points`, `google_place_visits`, `google_transit_stops`, `google_waypoints` | Only Activity and Location data is processed. You should specify the JSON format.
| Apple Health: `health` | Apple Health iOS app settings | `health_activity_summaries`, `health_beats_per_minutes`, `health_data_exports`, `health_devices`, `health_entries`, `health_entry_metadata`, `health_me_records`, `health_metadata_entries`, `health_workout_events`, `health_workout_metadata`, `health_workout_route_metadata`, `health_workout_route_track_points`, `health_workout_routes`, `health_workouts`, `health_activity_summaries`, `health_beats_per_minutes`, `health_data_exports`, `health_devices`, `health_entries`, `health_entry_metadata`, `health_me_records`, `health_metadata_entries`, `health_workout_events`, `health_workout_metadata`, `health_workout_route_metadata`, `health_workout_route_track_points`, `health_workout_routes`, `health_workouts`
| Instagram: `instagram` | https://www.instagram.com/download/request/ | - `instagram_account_history`, `instagram_comment_hashtag_mentions`, `instagram_comment_user_mentions`, `instagram_comments`, `instagram_hashtags`, `instagram_likes`, `instagram_media`, `instagram_media_hashtag_mentions`, `instagram_media_user_mentions`, `instagram_profile_photos`, `instagram_registration_info`, `instagram_stories_activities`, `instagram_users`
| Netflix: `netflix` | https://www.netflix.com/account/getmyinfo | `netflix_billing_history`, `netflix_clickstream`, `netflix_devices`, `netflix_indicated_preferences`, `netflix_interactive_titles`, `netflix_ip_addresses`, `netflix_my_list`, `netflix_playback_related_events`, `netflix_playtraces`, `netflix_ratings`, `netflix_search_history`, `netflix_subscription_history`, `netflix_viewing_activity`
| Spotify: `spotify` | https://www.spotify.com/us/account/privacy/ | `spotify_streaming_history`
| Telegram: `telegram` | Telegram Desktop app => Settings | `telegram_chats`, `telegram_members`, `telegram_messages`, `telegram_poll_answers`, `telegram_text_attachments`
| Twitter: `twitter` |https://twitter.com/settings/download_your_data | `twitter_ad_impressions`, `twitter_ad_impressions_matched_targeting_criteria`, `twitter_advertisers`, `twitter_age_info_records`, `twitter_audience_and_advertiser_records`, `twitter_audience_and_advertisers`, `twitter_audience_and_lookalike_advertisers`, `twitter_conversations`, `twitter_device_infos`, `twitter_direct_message_reactions`, `twitter_direct_message_urls`, `twitter_direct_messages`, `twitter_email_address_changes`, `twitter_gender_info`, `twitter_hashtags`, `twitter_inferred_age_info_records`, `twitter_interest_records`, `twitter_language_records`, `twitter_likes`, `twitter_locations`, `twitter_login_ips`, `twitter_personalization_locations`, `twitter_personalization_records`, `twitter_personalization_shows`, `twitter_screen_name_changes`, `twitter_shows`, `twitter_targeting_criteria`, `twitter_tweet_entities`, `twitter_tweet_hashtags`, `twitter_tweet_media`, `twitter_tweet_urls`, `twitter_tweet_user_mentions`, `twitter_tweets`, `twitter_urls`, `twitter_users`
| Chrome: `chrome` | OS X: ~/Library/Application Support/Google/Chrome/Default/History<br />Windows: C:\\Users\%USERNAME%\AppData\Local\Google\Chrome\User Data\Default\History<br />Linux: ~/.config/google-chrome/Default/databases | `chrome_segments`, `chrome_urls`, `chrome_visits`

## As a package

## Contributing

We appreciate contributions a lot! Here are some ways you can contribute:

* **Providers**. You can create new sources of data. Check out [#new-provider issues](https://github.com/bionic-dev/bionic/issues?q=is%3Aissue+is%3Aopen+label%3Anew-provider) and [an example PR with a new provider](https://github.com/bionic-dev/bionic/pull/41). Many existing providers lack some of the data: for example, the Google provider only proccesses a small subset of the Google export. Feel free to change it! We also target to test all providers and adding tests (especially, with unusual corner cases you found in your data) could be a very helpful contribution.
* **Views**. Views are additional SQL tables based on data from providers. Check out [an example PR with new views](https://github.com/bionic-dev/bionic/pull/29/files).
* **Docs**. 
* **Ecosystem**. Create and release your own tools based on Bionic databases. Think a web UI to visualize life or a custom Spotify Year In Review report generator.

When contributing, feel free to create issues and discussions with any questions. We promise to be helpful and kind!
