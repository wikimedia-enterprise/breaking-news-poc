package main

import (
	"breaking-news-poc/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/wikimedia-enterprise/wmf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// Article cache expiration
	exp = 30 * time.Minute

	// Row limit
	lmt = 100
)

func main() {
	time.Sleep(time.Second * 5)

	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)), &gorm.Config{})

	if err != nil {
		log.Panic(err)
	}

	str := persistence.NewInMemoryStore(time.Hour)
	rtr := gin.Default()
	rtr.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowWildcard:   true,
		AllowOrigins:    []string{"https://*.wikipediaenterprise.org", "https://*.enterprise.wikimedia.com"},
		AllowMethods:    []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders:    []string{"*"},
	}))

	rtr.GET("/articles", func(gcx *gin.Context) {
		pms := gcx.Request.URL.Query()
		// Construct a custom cache key based on the URL path and query parameters.
		key := gcx.Request.URL.Path + "?" + pms.Encode()

		ars := []*models.Article{}
		// Try to fetch from the cache first.
		err := str.Get(key, &ars)

		if err == nil {
			gcx.JSON(http.StatusOK, ars)
			return
		}

		if err == persistence.ErrNotStored {
			log.Printf("key %s was found in cache, but unable to set to the passed interface with error '%s'\n", key, err)
			gcx.JSON(http.StatusInternalServerError, err)
			return
		}

		ors := map[string]string{
			"name":                                "name",
			"editors":                             "(CASE WHEN ARRAY_LENGTH(editors, 1) IS NULL THEN 0 ELSE ARRAY_LENGTH(editors, 1) END)",
			"number_of_edits":                     "number_of_edits",
			"project":                             "project",
			"editors_within_hour":                 "editors_within_hour",
			"anonymous_editors_within_hour":       "anonymous_editors_within_hour",
			"anonymous_editors_ratio_within_hour": "anonymous_editors_ratio_within_hour",
			"url":                                 "url",
			"indications":                         "(CASE WHEN ARRAY_LENGTH(indications, 1) IS NULL THEN 0 ELSE ARRAY_LENGTH(indications, 1) END)",
			"date_created":                        "(CASE WHEN date_created IS NULL THEN TIMESTAMP '2004-10-19 10:23:54' ELSE date_created END)",
			"date_modified":                       "(CASE WHEN date_modified IS NULL THEN TIMESTAMP '2004-10-19 10:23:54' ELSE date_modified END)",
			"date_namespace_moved":                "(CASE WHEN date_namespace_moved IS NULL THEN TIMESTAMP '2004-10-19 10:23:54' ELSE date_namespace_moved END)",
		}

		qlm := gcx.Query("limit")

		if slm, err := strconv.Atoi(qlm); err == nil && len(qlm) > 0 {
			lmt = slm
		}

		qry := db.
			Model(new(models.Article))

		for _, key := range gcx.QueryArray("sort") {
			if pts := strings.Split(key, "|"); len(pts) == 2 {
				if srt, ok := ors[pts[0]]; ok {
					qry = qry.Order(fmt.Sprintf("%s %s", srt, strings.ToUpper(pts[1])))
				}
			}
		}

		if prs := gcx.QueryArray("project"); len(prs) > 0 {
			qry = qry.Where("project IN ?", prs)
		}

		if qds := gcx.QueryArray("qid"); len(qds) > 0 {
			qry = qry.Where("qid IN ?", qds)
		}

		if edt := gcx.Query("editors"); len(edt) > 0 {
			if eds := gcx.Query("editors_sign"); len(eds) > 0 {
				edn, _ := strconv.Atoi(edt)
				qry = qry.Where(fmt.Sprintf("array_length(editors, 1) %s ?", eds), edn)
			}
		}

		if ned := gcx.Query("number_of_edits"); len(ned) > 0 {
			if eds := gcx.Query("edits_sign"); len(eds) > 0 {
				edn, _ := strconv.Atoi(ned)
				qry = qry.Where(fmt.Sprintf("number_of_edits %s ?", eds), edn)
			}
		}
		if idc := gcx.Query("indications_count"); len(idc) > 0 {
			if ids := gcx.Query("indications_sign"); len(ids) > 0 {
				edn, _ := strconv.Atoi(idc)
				// when edn == 0 we don't want to filter by indications count, default=null
				if edn > 0 {
					qry = qry.Where(fmt.Sprintf("array_length(indications, 1) %s ?", ids), edn)
				}
			}
		}

		if ins := gcx.Query("indications"); len(ins) > 0 {
			if ina := strings.Split(ins, ","); len(ina) > 0 {
				for _, ind := range ina {
					qry = qry.Where("array_to_string(indications, ',') like ?", fmt.Sprintf("%%%s%%", ind))
				}
			}
		}

		dct := time.Now().UTC().Add(-24 * time.Hour)
		rws, err := qry.
			Where("(date_namespace_moved IS NOT null AND date_namespace_moved > ?) OR date_created > ?", dct, dct).
			Limit(lmt).
			Rows()

		if err != nil {
			gcx.JSON(http.StatusInternalServerError, err)
			return
		}

		for rws.Next() {
			art := new(models.Article)

			if err := db.ScanRows(rws, art); err != nil {
				gcx.JSON(http.StatusInternalServerError, err)
				return
			}

			ars = append(ars, art)
		}

		// Store the result in the cache with the custom cache key and expiration time.
		if err := str.Set(key, ars, exp); err != nil {
			log.Printf("Failed to set cache for key '%s': %v\n", key, err)
		}

		gcx.JSON(http.StatusOK, ars)
	})

	rtr.GET("/languages", cache.CachePage(str, time.Hour*24, func(gcx *gin.Context) {
		lns, err := wmf.NewClient().GetLanguages(gcx, "enwiki")

		if err != nil {
			gcx.JSON(http.StatusInternalServerError, err)
			return
		}

		gcx.JSON(http.StatusOK, lns)
	}))

	rtr.GET("/status", func(gcx *gin.Context) {
		gcx.Status(http.StatusOK)
	})

	rtr.POST("/rate", func(gcx *gin.Context) {
		rqb := new(models.RateQuery)
		err := gcx.ShouldBindJSON(rqb)

		if err != nil {
			gcx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fdk := new(models.Feedback)
		fdk.ID = rqb.Article.ID
		fdk.Project = rqb.Article.Project

		if db.Find(&fdk).RowsAffected == 0 {
			_ = copier.Copy(&fdk, &rqb.Article)
		}

		der := db.Save(&fdk).Error
		if der != nil {
			gcx.JSON(http.StatusInternalServerError, gin.H{"error": der.Error()})
		}

		db.Create(&models.Reaction{
			Rating:   rqb.Rating,
			Feedback: *fdk,
		})

		gcx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	rtr.POST("/comment", func(gcx *gin.Context) {
		cqy := new(models.CommentQuery)
		err := gcx.ShouldBindJSON(cqy)

		if err != nil {
			gcx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fdk := new(models.Feedback)
		fdk.ID = cqy.ID
		fdk.Project = cqy.Project

		if rad := db.Find(&fdk).RowsAffected > 0; rad {
			r := new(models.Reaction)
			db.Where("feedback_id = ? AND feedback_project = ? AND comment = ?", cqy.ID, cqy.Project, "").First(r)
			r.Comment = cqy.Comment
			err = db.Save(&r).Error

			if err != nil {
				gcx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			gcx.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			gcx.JSON(http.StatusBadRequest, gin.H{"error": "feedback does not exist"})
		}
	})

	if err := rtr.Run(":4042"); err != nil {
		log.Panic(err)
	}
}
