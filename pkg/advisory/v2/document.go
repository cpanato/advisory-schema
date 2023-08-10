package v2

import (
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/go-version"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const SchemaVersion = "2"

type Document struct {
	SchemaVersion string     `yaml:"schema-version"`
	Package       Package    `yaml:"package"`
	Advisories    Advisories `yaml:"advisories,omitempty"`
}

func (doc Document) Name() string {
	return doc.Package.Name
}

func (doc Document) Validate() error {
	return labelError(doc.Name(),
		errors.Join(
			doc.ValidateSchemaVersion(),
			doc.Package.Validate(),
			doc.Advisories.Validate(),
		),
	)
}

func (doc Document) ValidateSchemaVersion() error {
	docSchemaVersion, err := version.NewVersion(doc.SchemaVersion)
	if err != nil {
		return err
	}

	currentSchemaVersion, err := version.NewVersion(SchemaVersion)
	if err != nil {
		return err
	}

	if docSchemaVersion.GreaterThan(currentSchemaVersion) {
		return fmt.Errorf("document schema version %q is newer than the latest known schema version %q; if %q is supported by a later version of wolfictl, please update wolfictl and try this again", doc.SchemaVersion, SchemaVersion, doc.SchemaVersion)
	}

	// Document schema version also can't be earlier than the current schema version's MAJOR number.
	currentMajorNumber := currentSchemaVersion.Segments()[0]
	docMajorNumber := docSchemaVersion.Segments()[0]
	if docMajorNumber < currentMajorNumber {
		return fmt.Errorf("document schema version %q is too old to operate on with this version of wolfictl, document must use at least schema version \"%d\"", doc.SchemaVersion, currentMajorNumber)
	}

	return nil
}

func decodeDocument(r io.Reader) (*Document, error) {
	doc := &Document{}
	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)
	err := decoder.Decode(doc)
	if err != nil {
		return nil, err
	}

	if doc.SchemaVersion == "" {
		doc.SchemaVersion = "1"
	}

	return doc, nil
}

type Package struct {
	Name string `yaml:"name"`
}

func (p Package) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("package name must not be empty")
	}

	return nil
}

type Advisories []Advisory

func (advs Advisories) Validate() error {
	if len(advs) == 0 {
		return fmt.Errorf("this file should not exist if there are no advisories recorded")
	}

	return labelError("advisories",
		errors.Join(lo.Map(advs, func(adv Advisory, _ int) error {
			return adv.Validate()
		})...),
	)
}

func (advs Advisories) Contains(advisoryID string) bool {
	for _, adv := range advs {
		if adv.ID == advisoryID {
			return true
		}
	}

	return false
}

func (advs Advisories) Get(id string) (Advisory, bool) {
	for _, adv := range advs {
		if adv.ID == id {
			return adv, true
		}
	}

	return Advisory{}, false
}

// GetByVulnerabilityID returns the advisory that references the given ID as its
// advisory ID or as an alias. If such an advisory does not exist, the second
// return value will be false, otherwise it will be true.
func (advs Advisories) GetByVulnerabilityID(id string) (Advisory, bool) {
	for _, adv := range advs {
		if adv.ID == id {
			return adv, true
		}

		for _, alias := range adv.Aliases {
			if alias == id {
				return adv, true
			}
		}
	}

	return Advisory{}, false
}

func (advs Advisories) Update(id string, advisory Advisory) Advisories {
	for i, adv := range advs {
		if adv.ID == id {
			advs[i] = advisory
			return advs
		}
	}

	return advs
}
