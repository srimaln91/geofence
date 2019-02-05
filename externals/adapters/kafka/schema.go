package kafka

import (
	"log"

	schemaregistry "github.com/landoop/schema-registry"
)

type Schema struct {
	ID             int
	Schema         string
	Subject        string
	Version        int
	CurrentVersion int
}

type SchemaRegistry struct {
	Client          *schemaregistry.Client
	Schemas         map[string]map[int]Schema
	Versions        map[string][]int
	DefaultVersions map[string]int
}

// GetSchemas initializes all schemas
func GetSchemas(schemas map[string]int) (*SchemaRegistry, error) {

	client, _ := schemaregistry.NewClient(schemaregistry.DefaultURL)

	subjectSchemas := make(map[string]map[int]Schema)
	defaultVersionSchemas := make(map[string]int)

	for subject, defaultv := range schemas {

		log.Println("Default schema ID for topic ", subject, " ID: ", defaultv)

		versions, _ := client.Versions(subject)
		log.Println("Initial Schema Versions for ", subject, ": ", versions)

		for _, v := range versions {

			sch, _ := client.GetSchemaBySubject(subject, v)

			Fetchedschema := Schema{
				ID:      sch.ID,
				Schema:  sch.Schema,
				Subject: sch.Subject,
				Version: sch.Version,
			}

			log.Println("Registering Schema for subject: ", subject, " Version:", v)

			versionMap := make(map[int]Schema)
			versionMap[sch.Version] = Fetchedschema
			subjectSchemas[subject] = versionMap

		}

		defaultVersionSchemas[subject] = defaultv
	}

	schemaRegistry := SchemaRegistry{
		Client:          client,
		Schemas:         subjectSchemas,
		DefaultVersions: defaultVersionSchemas,
	}

	return &schemaRegistry, nil

}
