package main

import (
	"personae-fasti/api"
	"personae-fasti/data"
	"personae-fasti/opt"
)

var Config *opt.Conf
var Storage *data.Storage
var Api *api.APIServer

func main() {

	Config = opt.InitConfig()
	Storage = data.NewStorage(Config)
	Api = api.InitServer(Config, Storage)

}
