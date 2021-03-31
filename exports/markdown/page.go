package markdown

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type ChildType int

const (
	ChildSpotify ChildType = iota + 1
	ChildRescueTime
	ChildGooglePlaceVisit
)

func (ct ChildType) String() string {
	switch ct {
	case ChildSpotify:
		return "Spotify"
	case ChildRescueTime:
		return "RescueTime"
	case ChildGooglePlaceVisit:
		return "Google Places"
	default:
		panic("unknown child type")
	}
}

type Page struct {
	Title    string
	Children []Child
	Tag      string
}

type Child struct {
	String string
	Type   ChildType
	Time   time.Time
}

func (p *Page) Write(outputPath string) error {
	filename := fmt.Sprintf("%s.md", strings.Replace(p.Title, "/", "\\", -1))

	file, err := os.OpenFile(path.Join(outputPath, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	sort.Slice(p.Children, func(i, j int) bool {
		return p.Children[i].Time.Before(p.Children[j].Time)
	})

	for _, children := range p.Children {
		line := children.Time.Format("3:04 pm") + ": " + children.String
		escaped := strings.Replace(line, "/", "\\", -1)
		if _, err := file.WriteString(fmt.Sprintf("- %s\n", escaped)); err != nil {
			return err
		}
	}

	if p.Tag != "" {
		if _, err := file.WriteString(fmt.Sprintf("\n#%s", p.Tag)); err != nil {
			return err
		}
	}

	return nil
}
