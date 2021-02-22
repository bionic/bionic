# Bionic
Bionic is a command line tool to load your personal data exports from different services to a single SQLite database.

![Example of bionic usage](https://user-images.githubusercontent.com/6896447/108770008-dc6c4b80-756a-11eb-97ae-b6b84b21831f.png)

**Fun**: you can use Bionic to explore your own data and discover insights about yourself. Join tables between different sources to create reports like "Songs I listen in different locations" or dive deep into a single source to create "How amount of Telegram messages per week with different people changed over time" report.

**Research**: if you research human behaviour, subjective metrics could heavily impact your findings. If your respondents are able to run Bionic and send you aggregated results from their data, you can collect new objectives datasets describing important parts of life: transportation, social media, knowledge work and others.

**Education**: you can include Bionic exercises in your articles, courses or books. Learning to process data on personal records is much more exciting than processing artificial datasets. 

Bionic currently supports data exports from Google, Apple Health, Spotify, Telegram and Netflix.

## Install


## Usage

### Import data

### Generate views

### Query

You can query data with  ```sqlite3``` client:
```bash
$ sqlite3 db.sqlite                                               
SQLite version 3.28.0 2019-04-15 14:49:49
Enter ".help" for usage hints.

sqlite> select * from netflix_playback_related_events limit 1;
1|2021-01-22 20:46:21.696934+03:00|2021-01-22 20:46:21.696934+03:00||Seva|How I Met Your Mother: Season 1: "Come On"|Apple iPhone XR|RU|2020-12-30 20:14:21+00:00
```

Alternatively, you can use [datasette](https://github.com/simonw/datasette) to build a web ui to view and query data:

```bash
$ datasette serve db.sqlite
INFO:     Started server process [23975]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://127.0.0.1:8001 (Press CTRL+C to quit)
```

![datasette screenshot](https://user-images.githubusercontent.com/6896447/108776053-eb56fc00-7572-11eb-9081-1732cdc4a3bd.png)

You can also use Python and pandas to process data:
```python
import pandas as pd
import sqlite3

DATABASE_PATH = '/Users/seva/db.sqlite'
db_connection = sqlite3.connect(DATABASE_PATH)

messages_df = pd.read_sql('select * from telegram_messages;', con=db_connection)
```

## Supported exports


## Package

## Contributing

We need help!

providers => views scheme
