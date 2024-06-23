package setting

import (
	"github.com/bionicosmos/aegle/model"
)

var X model.Setting

func Init() {
	var err error
	X, err = model.FindSetting()
	if err != nil {
		panic(err)
	}
}
