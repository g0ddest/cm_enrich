package handlers

import (
	"bytes"
	"cm_enrich/internal/config"
	"cm_enrich/internal/models"
	"cm_enrich/internal/openai"
	"cm_enrich/internal/postgres"
	localsqs "cm_enrich/internal/sqs" // Переименован для избежания конфликта имен
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type AddressResponse struct {
	Country     string `json:"country"`
	FullAddress string `json:"full_address"`
	Kladr       string `json:"kladr"`
	Region      Region `json:"region"`
	Street      Street `json:"street"`
	City        City   `json:"city"`
}

type City struct {
	Kladr        string `json:"kladr"`
	Name         string `json:"name"`
	NameWithType string `json:"name_with_type"`
	Type         string `json:"type"`
	TypeFull     string `json:"type_full"`
}

type Region struct {
	Kladr        string `json:"kladr"`
	Name         string `json:"name"`
	NameWithType string `json:"name_with_type"`
	Type         string `json:"type"`
	TypeFull     string `json:"type_full"`
}

type Street struct {
	Kladr        string `json:"kladr"`
	Name         string `json:"name"`
	NameWithType string `json:"name_with_type"`
	Type         string `json:"type"`
	TypeFull     string `json:"type_full"`
}

func ProcessMessages(cfg *config.Config) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.SQSRegion),
	}))

	sqsSvc := sqs.New(sess)

	for {
		messages, err := localsqs.ReceiveMessages(sqsSvc, cfg.SQSQueueURL)
		if err != nil {
			log.Printf("Error receiving messages: %v", err)
			continue
		}

		for _, msg := range messages {
			var enrichmentMsg models.EnrichmentMsg
			if err := json.Unmarshal([]byte(*msg.Body), &enrichmentMsg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			if err := enrichAndSave(cfg, &enrichmentMsg); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}

			if err := localsqs.DeleteMessage(sqsSvc, cfg.SQSQueueURL, msg); err != nil {
				log.Printf("Error deleting message from SQS: %v", err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func enrichAndSave(cfg *config.Config, msg *models.EnrichmentMsg) error {
	streetType := ""
	if msg.StreetType != nil {
		streetType = *msg.StreetType
	}
	msg.StreetTypeRaw = streetType
	msg.StreetType = nil

	address := fmt.Sprintf("%s %s", msg.City, msg.Street)
	enrichedData, err := callAddressAPI(cfg.AddressAPIURL, address)
	if err != nil {
		log.Printf("Error calling address API: %v", err)
	}

	if enrichedData != nil && len(enrichedData) == 1 {
		addr := enrichedData[0]
		msg.RegionKladr = addr.Region.Kladr
		msg.RegionName = addr.Region.Name
		msg.RegionType = addr.Region.Type
		msg.StreetKladr = addr.Street.Kladr
		msg.StreetName = addr.Street.Name
		msg.StreetTypeFull = addr.Street.TypeFull
		msg.CityKladr = addr.City.Kladr
		msg.CityName = addr.City.Name
		msg.CityType = addr.City.Type
	} else {
		if enrichedData == nil {
			log.Printf("No enrichment data was found %s", address)
		} else {
			log.Printf("%d enrichment data was found %s, asking AI", len(enrichedData), address)

			var fullAddresses []string
			for _, addr := range enrichedData {
				fullAddresses = append(fullAddresses, addr.FullAddress)
			}

			selectedVariant, err := openai.CallOpenAI(
				cfg.OpenAIAPIKey,
				fmt.Sprintf("%s %s %s", msg.City, func() string {
					if msg.StreetType != nil {
						return *msg.StreetType
					}
					return ""
				}(), msg.Street),
				fullAddresses,
			)
			if err != nil {
				log.Printf("Error calling OpenAI: %v", err)
			} else {
				log.Printf("AI selected variant: %s", fullAddresses[selectedVariant])
				addr := enrichedData[selectedVariant]
				msg.RegionKladr = addr.Region.Kladr
				msg.RegionName = addr.Region.Name
				msg.RegionType = addr.Region.Type
				msg.StreetKladr = addr.Street.Kladr
				msg.StreetName = addr.Street.Name
				msg.StreetTypeFull = addr.Street.TypeFull
				msg.CityKladr = addr.City.Kladr
				msg.CityName = addr.City.Name
				msg.CityType = addr.City.Type
			}
		}
	}

	if err := postgres.SaveToPostgres(cfg.PostgresConnStr, msg); err != nil {
		return fmt.Errorf("error saving to Postgres: %v", err)
	}

	return nil
}

func callAddressAPI(apiURL, address string) ([]AddressResponse, error) {
	form := url.Values{}
	form.Add("address", address)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	var addressResp []AddressResponse
	if err := json.NewDecoder(resp.Body).Decode(&addressResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return addressResp, nil
}
