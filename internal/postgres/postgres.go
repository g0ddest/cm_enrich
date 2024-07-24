package postgres

import (
	"cm_enrich/internal/models"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4"
)

func SaveToPostgres(connStr string, msg *models.EnrichmentMsg) error {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("error connecting to PostgreSQL: %v", err)
	}
	defer conn.Close(context.Background())

	houseNumbers := strings.Join(msg.HouseNumbers, ",")
	houseRanges := strings.Join(msg.HouseRanges, ",")

	query := `INSERT INTO communal_outages (
		message_id, incident_id, organization, short_description, event, event_start, event_stop,
		city, street_type, street_type_raw, street, service, house_numbers, house_ranges,
		region_kladr, region_name, region_type, street_kladr, street_name,
		city_kladr, city_name, city_type
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
	)`

	_, err = conn.Exec(context.Background(), query,
		msg.MP, msg.ID, msg.Organization, msg.ShortDescription, msg.Event, msg.EventStart, msg.EventStop,
		msg.City, msg.StreetType, msg.StreetTypeRaw, msg.Street, msg.Service, houseNumbers, houseRanges,
		msg.RegionKladr, msg.RegionName, msg.RegionType, msg.StreetKladr, msg.StreetName,
		msg.CityKladr, msg.CityName, msg.CityType,
	)

	if err != nil {
		return fmt.Errorf("error executing insert: %v", err)
	}

	log.Printf("Inserted message %s into PostgreSQL", msg.ID)
	return nil
}
