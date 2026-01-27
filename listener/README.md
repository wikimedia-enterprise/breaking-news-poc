#### For connecting to the db container and running queries:

1. Connect to the running db container

    ```bash
    docker exec -it <container> bash
    ```

1. Connect to db

    ```bash
    psql -h localhost breaking_news
    ```

1. Some example queries:

* *Articles with indications:*

```bash
SELECT "name", "project", "editors", "edits", "date_created", "date_modified", "indications" FROM articles WHERE ARRAY_LENGTH("indications",1) > 0  LIMIT 2;
            name            | project |     editors      | edits |      date_created      |     date_modified      |      indications       
----------------------------+---------+------------------+-------+------------------------+------------------------+------------------------
 New_Zealand_swan           | enwiki  | {AmberDragon030} |     1 | 2004-02-11 13:27:53+00 | 2022-11-11 02:40:35+00 | {"Template:Cite news"}
 Mujhe_Khuda_Pay_Yaqeen_Hai | enwiki  | {Dl2000}         |     1 | 2021-04-02 14:19:19+00 | 2022-11-11 02:40:36+00 | {"Template:Cite news"}
(2 rows)
```

* *Articles with multiple editors:*

```bash
SELECT "name", "project", "editors", "edits", "date_created", "date_modified", "indications" FROM articles WHERE ARRAY_LENGTH("editors",1) > 1  LIMIT 2;

          name           | project |           editors           | edits | date_created |     date_modified      | indications 
-------------------------+---------+-----------------------------+-------+--------------+------------------------+-------------
 User_talk:Im_bored6823  | enwiki  | {"Im bored6823",Acroterion} |     4 |              | 2022-11-11 02:40:40+00 | {}
 Morte_de_Lázaro_Barbosa | ptwiki  | {Editoramado005,Horcoff}    |     3 |              | 2022-11-11 02:40:38+00 | {}
(2 rows)
```

* *Articles with creation date after Sep, 2022*


```bash
 SELECT "name", "project", "editors", "edits", "date_created", "date_modified", "indications" FROM articles WHERE   EXTRACT(MONTH FROM date_created) > 9 AND EXTRACT(YEAR FROM date_created) > 2021  LIMIT 2;
      name      |   project    |   editors    | edits |      date_created      |     date_modified      | indications 
----------------+--------------+--------------+-------+------------------------+------------------------+-------------
 Anny_Ogrezeanu | dewiki       | {RAL1028}    |     1 | 2022-11-05 14:06:57+00 | 2022-11-11 02:40:44+00 | {}
 goruk          | mswiktionary | {Lilindas00} |     2 | 2022-11-11 02:41:07+00 | 2022-11-11 02:41:07+00 | {}
(2 rows)
```