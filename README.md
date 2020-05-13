Companion code for the article [This article has x reactions and y comments](https://dev.to/napicella/this-article-has-x-reactions-and-y-comments-1jlf)

The program uses the [dev.to](https://docs.dev.to/api/) API to keep the title of the [article](https://dev.to/napicella/this-article-has-x-reactions-and-y-comments-1jlf) up to date with the number of reactions and comments it received.

## Build
```bash
make build
```

## Run
```bash
ARTICLE_ID=123456 DEV_TO_API_KEY=ABCDEFGHILMN123456789012 ./bin/update-article
```

## Cron job 
```bash
SHELL=/bin/bash
* * * * * /home/user/wrapper.sh
```
