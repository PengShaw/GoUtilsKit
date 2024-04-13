package templater_test

import (
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	"github.com/PengShaw/GoUtilsKit/templater"
)

func TestRenderText(t *testing.T) {
	funcs := template.FuncMap{"ToUpper": strings.ToUpper}
	r, err := templater.RenderText("_", "hello, {{ . | ToUpper }}!", "world", funcs)
	assert.NoError(t, err, "should not be an error")
	assert.Equal(t, []byte("hello, WORLD!"), r, "they should be equal")
}
