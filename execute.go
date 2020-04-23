package main

import (
	"log"

	"example.com/amidist/command"
)

func execute(dsc string) {
	executeResources(dsc)
}

func executeResources(dsc string) {
	c := prepareDataStore(dsc)
	resources := command.Resources
	for _, rawResource := range c.RawResources {
		resource, ok := resources[rawResource.Type]
		if !ok {
			log.Fatalf(" InvalidResource: Resource %s Not Found", rawResource.Type)
		}

		if err := resource.Prepare(rawResource.Config); err != nil {
			log.Fatalln(err)
		}

		if err := resource.Run(); err != nil {
			log.Fatalln(err)
		}
	}
}

func prepareDataStore(dsc string) config {
	log.Println(" ****************** Start: Template Parsing *********************")

	parsedTemplateFileName := parseTemplate(dsc)
	config := unmarshalYaml(parsedTemplateFileName)

	log.Println(" ****************** End: Template Parsing *********************")

	return config
}
