package memory

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/spf13/viper"
	"github.com/tingold/boring/config"
	"github.com/tingold/boring/model"
	"go.uber.org/zap"
	"strings"
)

//MemoryCatalog is an in-memory implementation of the data catalog
type MemoryCatalog struct {
	db *pgxpool.Pool
	TableMap map[string]Table
}

func (mc *MemoryCatalog) Load() error {
	// connect to db
	if mc.db == nil {
		err := mc.connectPostgres()
		if err != nil {
			return err
		}
	}
	//find spatial tables
	sql := "SELECT f_table_schema, f_table_name, f_geometry_column, coord_dimension, srid, type from geometry_columns"
	rows, err := mc.db.Query(context.Background(), sql)
	if err != nil{
		zap.L().Error("unable to search geometry_columns table")
		return err
	}
	tempMap := make(map[string]Table)
	for rows.Next(){

		table := Table{}
		err := rows.Scan(&table.Schema, &table.Name, &table.GeometryColumnName, &table.Dimension, &table.SRID, &table.GeometryType)
		if err != nil{
			zap.S().Errorf("unable to scan row: %s", err.Error())
			continue
		}
		//potentially have multiple geometry cols per table...only use the first for now
		//skip if it exists
		if _, ok := tempMap[table.Name]; ok {
			continue
		}
		//are we only exposing some tables?
		if viper.GetBool(config.POSTGRES_USE_TABLE_PREFIX){
			//check to see if it matches the pattern before adding it
			if strings.HasPrefix(table.Name, viper.GetString(config.POSTGRES_TABLE_PREFIX)){

				tempMap[table.Name] = table
				mc.TableMap = tempMap
			}
		} else {
			tempMap[table.Name] = table
			mc.TableMap = tempMap
		}

	}
	rows.Close()

	sql = "SELECT St_AsBinary(St_SetSRID(ST_Extent($1),4326)) as bbox from $2"
	good := 0

	for _, table := range tempMap{
		//can't use table names as parameters :/
		tempSql := strings.Replace(sql, "$2",table.Schema+"."+table.Name, 1)
		tempSql = strings.Replace(tempSql, "$1", table.GeometryColumnName,1)
		row := mc.db.QueryRow(context.Background(),tempSql)
		var bbox []byte
		err = row.Scan(&bbox)
		if err != nil {
			zap.S().Errorf("unable to determine extent of %s: %s",table.Name, err.Error())
			continue
		}
		scanner := wkb.Scanner(&table.Extent)
		err = scanner.Scan(bbox)
		if err != nil{
			zap.S().Errorf("unable to scan geometry: %s", err.Error())
			continue
		}
		good ++

	}
	zap.S().Infof("Got %d good bboxes out of %d",good, len(tempMap) )
	zap.S().Infof("found %d spatial tables", len(tempMap))

	return nil
}

func (mc *MemoryCatalog) GetCollections(ctx context.Context) (*model.Collections, error) {

	allCollections := make([]model.Collection,0)

	for _, table := range mc.TableMap {
		allCollections = append(allCollections, mc.tableToCollection(table, fmt.Sprintf("%s",ctx.Value("boring_base"))))
	}
	collectionResponse := model.Collections{
		Links:       nil,
		Collections: &allCollections,
	}
	return &collectionResponse, nil
}

func (mc *MemoryCatalog) GetCollection(ctx context.Context,id string) (*model.Collection, error) {




	return nil, nil
}

func (mc *MemoryCatalog) connectPostgres() error {

	connstring := "database=" + viper.GetString(config.POSTGRES_DATABASE) + " host=" + viper.GetString(config.POSTGRES_HOST) + " user=" + viper.GetString(config.POSTGRES_USER) + " password=" + viper.GetString(config.POSTGRES_PASSWORD) + " port=" + viper.GetString(config.POSTGRES_PORT)

	dbLogger := zapadapter.NewLogger(zap.L())

	poolConfig, err := pgxpool.ParseConfig(connstring)

	if err != nil {
		zap.L().Fatal("Unable to parse connection string")
	}

	poolConfig.ConnConfig.Logger = dbLogger

	mc.db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)

	if err != nil {

		if viper.GetBool(config.POSTGRES_FAILFAST){
			zap.S().Fatalf("failed to connect to database: %s", err.Error())
		} else {
			zap.S().Errorf("failed to connect to database (will reattempt): %s", err.Error())
		}
	}

	return nil

}

func (mc *MemoryCatalog) tableToCollection(table Table, baseUrl string)(model.Collection){

	itemLink := model.Link{
		Href:  baseUrl+"/collections/"+table.Name+"/items",
		Rel:   "items",
		Type:  "application/json",
		Title: "items",
	}
	bounds := table.Extent.Bound()
	extent := model.Extent{
		SpatialExtent:  &model.SpatialExtent{
			CRS:  "http://www.opengis.net/def/crs/OGC/1.3/CRS84",
			BBox: &[]float64{bounds.Min.X(), bounds.Min.Y(), bounds.Max.X(), bounds.Max.Y()},
		},
		TemporalExtent: nil,
	}
	collection := model.Collection{

		Id:          table.Name,
		Title:       table.Name,
		Description: table.Name,
		Links:       &[]model.Link{itemLink},
		Extent:      &extent,
		ItemType:    "feature",
		CRS:         &[]string{"http://www.opengis.net/def/crs/OGC/1.3/CRS84"},
	}

	return collection
}