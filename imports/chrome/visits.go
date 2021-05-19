package chrome

import (
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
)

type Visit struct {
	gorm.Model
	URLID   int
	URL     URL
	Time    types.DateTime
	VisitID int // Parent visit
	Visit   *Visit

	// See https://chromium.googlesource.com/chromium/+/trunk/content/public/common/page_transition_types.h
	TransitionType          TransitionType
	TransitionQualifierType TransitionQualifierType
	IsRedirect              bool

	SegmentID                    int
	Segment                      Segment
	VisitDuration                int
	IncrementedOmniboxTypedScore bool
	PubliclyRoutable             bool
}

func (Visit) TableName() string {
	return tablePrefix + "visits"
}

func (visit Visit) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"id":     visit.ID,
		"url_id": visit.URLID,
	}
}

type TransitionType string

const (
	TransitionLink             TransitionType = "LINK"
	TransitionTyped            TransitionType = "TYPED"
	TransitionAutoBookmark     TransitionType = "AUTO_BOOKMARK"
	TransitionAutoSubframe     TransitionType = "AUTO_SUBFRAME"
	TransitionManualSubframe   TransitionType = "MANUAL_SUBFRAME"
	TransitionGenerated        TransitionType = "GENERATED"
	TransitionAutoToplevel     TransitionType = "AUTO_TOPLEVEL"
	TransitionSubmit           TransitionType = "FORM_SUBMIT"
	TransitionReload           TransitionType = "RELOAD"
	TransitionKeyword          TransitionType = "KEYWORD"
	TransitionKeywordGenerated TransitionType = "KEYWORD_GENERATED"
)

type TransitionQualifierType string

const (
	TransitionQualifierForwardBack    TransitionQualifierType = "FORWARD_BACK"
	TransitionQualifierFromAddressBar TransitionQualifierType = "FROM_ADDRESS_BAR"
	TransitionQualifierHomePage       TransitionQualifierType = "HOME_PAGE"
	TransitionQualifierChainStart     TransitionQualifierType = "CHAIN_START"
	TransitionQualifierChainEnd       TransitionQualifierType = "CHAIN_END"
	TransitionQualifierClientRedirect TransitionQualifierType = "CLIENT_REDIRECT"
	TransitionQualifierServerRedirect TransitionQualifierType = "SERVER_REDIRECT"
)

type VisitAlias Visit
type SourceVisit struct {
	VisitAlias
	Transition int64
}

func (p *chrome) importVisits(db *gorm.DB) error {
	selection := "id, url as url_id, datetime((visit_time/1000000)-11644473600, 'unixepoch') as time, " +
		"from_visit as visit_id, transition, segment_id, visit_duration, incremented_omnibox_typed_score, publicly_routable"

	var sourceVisits []SourceVisit

	err := db.
		Raw("select "+selection+" from visits order by id limit ?", dbRowSelectLimit).
		Scan(&sourceVisits).
		Error
	if err != nil {
		return err
	}
	if err := p.saveVisits(sourceVisits); err != nil {
		return err
	}
	for len(sourceVisits) == dbRowSelectLimit {
		lastVisit := sourceVisits[len(sourceVisits)-1]
		err = db.
			Raw("select "+selection+" from visits where id > ? order by id limit ?", lastVisit.ID, dbRowSelectLimit).
			Scan(&sourceVisits).
			Error
		if err != nil {
			return err
		}
		if err := p.saveVisits(sourceVisits); err != nil {
			return err
		}
	}

	return nil
}

func (p *chrome) saveVisits(sourceVisits []SourceVisit) error {
	for _, sourceVisit := range sourceVisits {
		visit := Visit(sourceVisit.VisitAlias)
		visit.TransitionType, visit.TransitionQualifierType, visit.IsRedirect = transitionInfoBySourceValue(sourceVisit.Transition)
		err := p.DB().
			FirstOrCreate(&visit, visit.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func transitionInfoBySourceValue(transition int64) (tt TransitionType, tqt TransitionQualifierType, isRedirect bool) {
	switch transition & 0xFF {
	case 0:
		tt = TransitionLink
	case 1:
		tt = TransitionTyped
	case 2:
		tt = TransitionAutoBookmark
	case 3:
		tt = TransitionAutoSubframe
	case 4:
		tt = TransitionManualSubframe
	case 5:
		tt = TransitionGenerated
	case 6:
		tt = TransitionAutoToplevel
	case 7:
		tt = TransitionSubmit
	case 8:
		tt = TransitionReload
	case 9:
		tt = TransitionKeyword
	case 10:
		tt = TransitionKeywordGenerated
	}

	switch transition & 0xFFFFFF00 {
	case 0x01000000:
		tqt = TransitionQualifierForwardBack
	case 0x02000000:
		tqt = TransitionQualifierFromAddressBar
	case 0x04000000:
		tqt = TransitionQualifierHomePage
	case 0x10000000:
		tqt = TransitionQualifierChainStart
	case 0x20000000:
		tqt = TransitionQualifierChainEnd
	case 0x40000000:
		tqt = TransitionQualifierClientRedirect
	case 0x80000000:
		tqt = TransitionQualifierServerRedirect
	}

	isRedirect = transition&0xC0000000 != 0

	return tt, tqt, isRedirect
}
