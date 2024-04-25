package container

import (
	"fmt"

	"github.com/sarulabs/di/v2"
)

var (
	container  di.Container
	defs       []di.Def
	defsByTags = map[string][]string{}
)

const App = di.App

func Add(def di.Def) {
	defs = append(defs, def)

	if len(def.Tags) > 0 {
		registerTags(def.Tags, def.Name)
	}
}

func registerTags(tags []di.Tag, name string) {
	for _, tag := range tags {
		defsByTags[tag.Name] = append(defsByTags[tag.Name], name)
	}
}

func Build(scopes ...string) error {
	var (
		builder *di.Builder
		err     error
	)
	if builder, err = di.NewBuilder(scopes...); err != nil {
		return fmt.Errorf("can't create builder: %w", err)
	}

	if err = builder.Add(defs...); err != nil {
		return fmt.Errorf("can't add definitions: %w", err)
	}

	container = builder.Build()

	return nil
}

func Fill(name string, dst interface{}) error {
	return container.Fill(name, dst)
}

func GetByTag(tagName string) []string {
	return defsByTags[tagName]
}
