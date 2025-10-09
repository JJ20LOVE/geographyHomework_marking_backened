package main

import (
	"dbdemo/api"
	"dbdemo/model"
	"dbdemo/routers"
	"dbdemo/utils"
)

func main() {
	model.InitDb()
	utils.InitMinIO()
	api.InitWorker()
	routers.InitRouter()
}
