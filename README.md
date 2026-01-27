# WME Breaking News PoC

Create `.env` file in the project root with following content:

    ```bash
    VUE_APP_API_URL=http://localhost:4042
    MSGS_PER_SECOND=2
    ```

If running locally, update the CORS config in `api/main.go` to allow `localhost`:
  ```bash
    AllowOrigins:    []string{"https://*.wikipediaenterprise.org", "https://*.enterprise.wikimedia.com", "http://localhost/*"},
  ```

After that you can run:

```bash
  docker compose up -d --build
```

Once it's running, access in `http://localhost:8082`.

To stop the execution:

```bash
  docker compose stop
```

To get the latest updates simply:

```bash
  git pull origin main
```



## Debugging
To inspect the DB tables, connect to the container using `psql` cli and run sql queries as follows:

```bash
  # Connect as user root
  docker exec -it breaking-news-poc_db_1 psql -U root

  # Once inside the container to get all the tables, use \dt
  breaking_news=# \dt
         List of relations
  Schema |   Name    | Type  | Owner 
  --------+-----------+-------+-------
  public | articles  | table | root
  public | feedbacks | table | root
  public | reactions | table | root

  # Connect to the database
  breaking_news=# \c breaking_news
  You are now connected to database "breaking_news" as user "root".

  # Run queries
  breaking_news=# SELECT * FROM reactions;
  id | rating | comment | feedback_id | feedback_project 
  ----+--------+---------+-------------+------------------
    3 | like   |         |    74643188 | enwiki
    1 | like   | tst1    |     8486756 | zhwiki
    2 | like   | test1   |    74648015 | enwiki
```