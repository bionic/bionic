# Bionic
Bionic is a command line tool to load your personal data exports from different services to a single SQLite database. Bionic currently supports data exports from Google, Apple Health, Spotify, Telegram and Netflix.

According to GDPR, every service has to provide an ability for users to export their data in a machine-readable format. However, the data format varies between different services providers, making it hard for advanced users to play around with the data. Bionic solves this problem by unifiying different GDPR exports to a single SQL schema.

![Example of bionic usage](https://user-images.githubusercontent.com/6896447/108840862-f9932f80-75e7-11eb-9014-70afc55ff302.png)

**Fun**: you can use Bionic to explore your own data and discover insights about yourself. Join tables between different sources to create reports like "Songs I listen in different locations" or dive deep into a single source to create "How amount of Telegram messages per week with different people changed over time" report.

**Research**: if you research human behaviour, subjective metrics could heavily impact your findings. If your respondents are able to run Bionic and send you aggregated results from their data, you can collect new objectives datasets describing important parts of life: transportation, social media, knowledge work and others.

**Development**: you can use Bionic as a Go package to implement personal data import in your apps.

**Education**: you can include Bionic exercises in your articles, courses or books. Learning to process data on personal records is much more exciting than processing artificial datasets. 

## Install


## Usage

### Import data

### Generate views

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

## Supported exports


## As a package

## Contributing

We appreciate contributions a lot! Here are some of the ways you can contribute:

* **Providers**. You can create new sources of data. Check out [#new-provider issues](https://github.com/bionic-dev/bionic/issues?q=is%3Aissue+is%3Aopen+label%3Anew-provider) and [an example PR with a new provider](https://github.com/bionic-dev/bionic/pull/41). Many existing providers lack some of the data: for example, the Google provider only proccesses a small subset of the Google export. Feel free to change it! We also target to test all providers and adding tests (especially, with unusual corner cases you found in your data) could be a very helpful contribution.
* **Views**. Views are additional SQL tables based on data from providers. Check out [an example PR with new views](https://github.com/bionic-dev/bionic/pull/29/files).
* **Docs**. 
* **Recipes**.
* **Ecosystem**. Create and release your own tools based on Bionic databases. Think a web UI to visualize life or a custom Spotify Year In Review report generator.

When contributing, feel free to create issues and discussions with any questions. We promise to be helpful and kind!
