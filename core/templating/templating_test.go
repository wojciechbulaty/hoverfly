package templating

import (
	"testing"
	. "github.com/onsi/gomega"
)


func Test_ShouldApplyTemplateToQueryParams(t *testing.T) {
	RegisterTestingT(t)

	Expect(ApplyTemplate()).To(Equal(`

	`))
}