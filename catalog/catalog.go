package catalog

import (

	"github.com/tingold/boring/model"
)

var catalog Catalog

type Catalog interface {

	Load() error
	GetCollections() (*[]model.Collection, error)
	GetCollection(id string)(*model.Collection, error)
}

func GetCatalog() (*Catalog, error){

	//only supporting memory catalog right now -- in the future we can write to postgres
	//if viper.GetBool(config.POSTGRES_READONLY){}
	catalog = MemoryCatalog{}


}