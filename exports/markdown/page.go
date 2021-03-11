package markdown

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
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
	default:
		return ""
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
		return p.Children[i].Type < p.Children[j].Type
	})

	var previousChildren ChildType

	for _, children := range p.Children {
		if children.Type != previousChildren {
			if _, err := file.WriteString(fmt.Sprintf("# %s\n", children.Type)); err != nil {
				return err
			}
			previousChildren = children.Type
		}

		escaped := strings.Replace(children.String, "/", "\\", -1)
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
