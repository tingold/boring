package catalog

import "github.com/tingold/boring/model"

//MemoryCatalog is an in-memory implementation of the data catalog
type MemoryCatalog struct {




}

func (mc MemoryCatalog) Load() error {

	return nil
}
func (mc MemoryCatalog) GetCollections() (*[]model.Collection, error) {

	return nil,nil
}

func (mc MemoryCatalog) GetCollection(id string)(*model.Collection, error){
	return nil, nil
}