# logger

```go
package main

import (
	"errors"

	"github.com/PengShaw/GoUtilsKit/logger"
)

func main() {
	logger.SetLevel(logger.LevelInfo)
	logger.Debug("This will not print")
	logger.Infoln("Print info")

	log := logger.New(logger.LevelError)
	log.Info("This will not print")
	logger.Errorln("Print error")
	log.Fatalf("Print error(%s) and followed by a call to os.Exit(1)", errors.New("some error"))
}
```