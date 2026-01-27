package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"breaking-news-poc/models"
	"breaking-news-poc/utils"

	"github.com/wikimedia-enterprise/wmf"
	eventstream "github.com/wikimedia-enterprise/wmf-event-stream-sdk-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cron "github.com/robfig/cron/v3"
)

type entry struct {
	prj *wmf.Project
	evt *eventstream.RevisionCreate
}

func notInNamespace(nsp int) bool {
	return nsp != 0
}

func printErrors(errs chan error) {
	for err := range errs {
		if err != nil && err != wmf.ErrProjectNotFound {
			log.Println(err)
		}
	}
}

func main() {
	time.Sleep(time.Second * 5)

	ctx := context.Background()
	snc := time.Now()
	sts := eventstream.NewClient()
	sts.SetUserAgent("Breaking News PoC (WME)")
	clt := wmf.NewClient()

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

	if err := db.AutoMigrate(new(models.Article), new(models.Feedback), new(models.Reaction)); err != nil {
		log.Panic(err)
	}

	msgsPerSecond := math.MaxFloat64
	env, set := os.LookupEnv("MSGS_PER_SECOND")
	if set {
		msgsPerSecond, err = strconv.ParseFloat(env, 64)
		if err != nil {
			log.Panic(err)
		}
	}

	t := rate.NewLimiter(rate.Limit(msgsPerSecond), 1)
	que := make(chan entry, 10000)
	hdl := func(prj *wmf.Project, evt *eventstream.RevisionCreate) error {
		art := new(models.Article)
		art.ID = evt.Data.PageID
		art.Project = evt.Data.Database
		art.URL = evt.Data.Meta.URI

		if err := db.Find(art).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		art.Name = evt.Data.PageTitle
		art.ProjectURL = prj.URL

		if err = t.Wait(ctx); err != nil {
			if err != nil {
				log.Printf("failed to throttle: %v\n", err)
			}
		}

		if err := art.GetDataFromAPI(ctx, clt); err != nil {
			return err
		}

		art.DateModified = &evt.Data.RevTimestamp
		art.NumberOfEdits += 1
		editors := utils.AppendToSet(art.Editors, evt.Data.Performer.UserText)

		if len(editors) > len(art.Editors) {
			if art.DateCreated != nil && art.DateNamespaceMoved == nil {
				if art.DateModified.Before(art.DateCreated.Add(time.Minute * 61)) {
					art.CalculateEditingRatio(evt)
				}
			} else if art.DateNamespaceMoved != nil {
				if art.DateModified.Before(art.DateNamespaceMoved.Add(time.Minute * 61)) {
					art.CalculateEditingRatio(evt)
				}
			}
		}

		art.Editors = editors

		return db.Save(art).Error
	}

	go func() {
		rct := sts.RevisionCreate(ctx, snc, func(evt *eventstream.RevisionCreate) error {
			if notInNamespace(evt.Data.PageNamespace) {
				return nil
			}

			if evt.Data.PageIsRedirect {
				return nil
			}

			prj, err := clt.GetProject(ctx, evt.Data.Database)

			if err != nil {
				return err
			}

			que <- entry{
				evt: evt,
				prj: prj,
			}

			return nil
		})

		printErrors(rct.Sub())
	}()

	for i := 0; i < 25; i++ {
		go func() {
			for ent := range que {
				if err := hdl(ent.prj, ent.evt); err != nil {
					log.Println(err)
				}
			}
		}()
	}

	go func() {
		pgm := sts.PageMove(ctx, snc, func(evt *eventstream.PageMove) error {
			if notInNamespace(evt.Data.PageNamespace) {
				return nil
			}

			if _, err := clt.GetProject(ctx, evt.Data.Database); err != nil {
				return err
			}

			if evt.Data.PriorState.PageNamespace == 0 {
				return nil
			}

			time.Sleep(time.Second * 2)
			art := new(models.Article)
			art.ID = evt.Data.PageID
			art.Project = evt.Data.Database
			art.Name = evt.Data.PageTitle

			if err := db.Find(art).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			art.DateNamespaceMoved = &evt.Data.Meta.Dt

			return db.Save(art).Error
		})

		printErrors(pgm.Sub())
	}()

	go func() {
		c := cron.New()
		_, err := c.AddFunc("@every 1h00m00s", func() {
			db.Where("date_namespace_moved < ? OR date_created < ?",
				time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, -1)).Delete(&models.Article{})
		})
		c.Start()

		if err != nil {
			log.Panic(err)
		}
	}()

	swg := new(sync.WaitGroup)
	swg.Add(1)
	swg.Wait()
}
