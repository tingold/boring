package catalog

import (
	"github.com/tingold/boring/catalog/memory"
)

var catalog *memory.MemoryCatalog

func GetCatalog() (*memory.MemoryCatalog, error){

	//only supporting memory catalog right now -- in the future we can write to postgres
	//if viper.GetBool(config.POSTGRES_READONLY){}
	if catalog == nil{
		cat := memory.MemoryCatalog{}
		err := cat.Load()
		if err != nil {
			return nil,err
		}
		catalog = &cat
	}
	return catalog, nil
}


