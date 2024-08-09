package setting

import (
	"github.com/bionicosmos/aegle/model"
)

var X model.Setting

func Init() {
	var err error
	X, err = model.LoadSettings()
	if err != nil {
		panic(err)
	}
}
