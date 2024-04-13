# templater

```golang
package main

import (
	"html/template"
	"strings"

	"github.com/PengShaw/GoUtilsKit/templater"
)

func main() {
	funcs := template.FuncMap{"ToUpper": strings.ToUpper}
	r, err := templater.RenderText("_", "hello, {{ . | ToUpper }}!", "world", funcs)
	if err != nil {
		panic(err)
	}
	println(string(r))
}
```

