echo "Starting Postgis"
docker run -e POSTGRES_DB=boring -e POSTGRES_USER=postgis -p 5432:5432 -e POSTGRES_PASSWORD=password -d mdillon/postgis
echo "Sleeping for 20 seconds to allow time for Postgis to start..."
sleep 20
echo "Loading Coastlines (linestring)"
ogr2ogr -f PostgreSQL "PG:user=postgis host=localhost password=password dbname=boring" -skipfailures https://d2ad6b4ur7yvpq.cloudfront.net/naturalearth-3.3.0/ne_50m_coastline.geojson
echo "Loading Admin Areas (polygon)"
ogr2ogr -f PostgreSQL "PG:user=postgis host=localhost password=password dbname=boring" -skipfailures https://d2ad6b4ur7yvpq.cloudfront.net/naturalearth-3.3.0/ne_50m_admin_0_map_units.geojson
echo "Loading Airports (points)"
ogr2ogr -f PostgreSQL "PG:user=postgis host=localhost password=password dbname=boring" -skipfailures https://d2ad6b4ur7yvpq.cloudfront.net/naturalearth-3.3.0/ne_10m_airports.geojson
echo "Done!"