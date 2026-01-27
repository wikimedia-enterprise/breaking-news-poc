package models

import (
	"breaking-news-poc/utils"
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/wikimedia-enterprise/wmf"
	eventstream "github.com/wikimedia-enterprise/wmf-event-stream-sdk-go"
)

var brkNewsTmps = []string{
	"Template:Cite news",
	"Template:In the news/footer",
	"Template:Wikinews",
	"Template:Editnotice",
	"Template:ITN candidate",
	"Template:In the news",
	"Template:In use",
	"Template:Recent death",
	"Template:Recent death presumed",
}

var brkNewsStartsWithTmps = []string{
	"Template:Current",
}

var brkNewsCats = []string{"Category:Current events",
	"Category:News",
	"Category:Current events portal",
}

type Article struct {
	ID                              int            `gorm:"primaryKey" json:"id"`
	Name                            string         `json:"name"`
	Editors                         pq.StringArray `json:"editors" gorm:"type:varchar(255)[]"`
	EditorsWithinHour               int            `json:"editors_within_hour" gorm:"not null"`
	AnonymousEditorsWithinHour      int            `json:"anonymous_editors_within_hour" gorm:"not null"`
	AnonymousEditorsRatioWithinHour float64        `json:"anonymous_editors_ratio_within_hour" gorm:"not null"`
	QID                             string         `json:"qid,omitempty" gorm:"column:qid;index:idx_qid"`
	NumberOfEdits                   int            `gorm:"index" json:"number_of_edits"`
	Project                         string         `gorm:"primaryKey" json:"project"`
	ProjectURL                      string         `json:"project_url"`
	URL                             string         `json:"url"`
	Indications                     pq.StringArray `json:"indications" gorm:"type:varchar(255)[]"`
	DateCreated                     *time.Time     `gorm:"index" json:"date_created"`
	DateModified                    *time.Time     `json:"date_modified"`
	DateNamespaceMoved              *time.Time     `gorm:"index" json:"date_namespace_moved"`
}

func (a *Article) GetDataFromAPI(ctx context.Context, clt wmf.PageGetter) error {
	pge, err := clt.GetPage(ctx, a.Project, a.Name, func(v *url.Values) {
		if a.DateCreated == nil {
			v.Set("prop", "templates|categories|revisions|pageprops")
			v.Set("ppprop", "wikibase_item")
			v.Set("rvprop", "ids|timestamp")
			v.Set("rvdir", "newer")
			v.Set("rvlimit", "1")
		} else {
			v.Set("prop", "templates|categories|pageprops")
			v.Set("ppprop", "wikibase_item")
		}

		v.Set("cllimit", "500")
		v.Set("tllimit", "500")
	})

	if err != nil {
		return err
	}

	if pge.PageProps != nil && len(pge.PageProps.WikiBaseItem) > 0 {
		a.QID = pge.PageProps.WikiBaseItem
	}

	if len(pge.Revisions) > 0 && a.DateCreated == nil {
		a.DateCreated = pge.Revisions[0].Timestamp
	}

	for _, tmp := range pge.Templates {
		for _, bTmp := range brkNewsTmps {
			if tmp.Title == bTmp {
				a.Indications = utils.AppendToSet(a.Indications, tmp.Title)
			}
		}

		for _, bTmp := range brkNewsStartsWithTmps {
			if strings.HasPrefix(tmp.Title, bTmp) {
				a.Indications = utils.AppendToSet(a.Indications, tmp.Title)
			}
		}
	}

	for _, cat := range pge.Categories {
		for _, bCat := range brkNewsCats {
			if cat.Title == bCat {
				a.Indications = utils.AppendToSet(a.Indications, cat.Title)
			}
		}
	}

	return nil
}

func (a *Article) CalculateEditingRatio(evt *eventstream.RevisionCreate) {
	a.EditorsWithinHour += 1

	if evt.Data.Performer.UserID == 0 {
		a.AnonymousEditorsWithinHour += 1
	}

	a.AnonymousEditorsRatioWithinHour = float64(a.AnonymousEditorsWithinHour) / float64(a.EditorsWithinHour)
}
