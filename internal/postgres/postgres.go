package postgres

import (
	"cm_enrich/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"strings"
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
	city, street_type, street, service, house_numbers, house_ranges,
	full_address, kladr, region_kladr, region_name, region_name_with_type, region_type, region_type_full,
	street_kladr, street_name, street_name_with_type, street_type_full
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
)`

	_, err = conn.Exec(context.Background(), query,
		msg.MP, msg.ID, msg.Organization, msg.ShortDescription, msg.Event, msg.EventStart, msg.EventStop,
		msg.City, msg.StreetType, msg.Street, msg.Service, houseNumbers, houseRanges,
		msg.FullAddress, msg.Kladr, msg.RegionKladr, msg.RegionName, msg.RegionNameWithType, msg.RegionType, msg.RegionTypeFull,
		msg.StreetKladr, msg.StreetName, msg.StreetNameWithType, msg.StreetTypeFull,
	)

	if err != nil {
		return fmt.Errorf("error executing insert: %v", err)
	}

	log.Printf("Inserted message %s into PostgreSQL", msg.ID)
	return nil
}
