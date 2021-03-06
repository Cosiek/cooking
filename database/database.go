package database

import (
	"errors"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("pantry.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Produce{})

	return db
}

// ============================================================================
// MODELS =====================================================================
// ============================================================================

var Mesures = map[int8]string{
	1: "szt.",
	2: "kg",
	3: "l",
}

// PRODUCE ============================

type Produce struct {
	gorm.Model
	Name   string
	Mesure int8
}

func (p *Produce) setName(name string) error {
	if len(name) > 30 {
		return errors.New("name to long")
	}
	p.Name = name
	return nil
}

func (p *Produce) setMesure(mesureIdStr string) error {
	unknownMesureMsg := "unknown mesure - choose one from the list."
	// try to convert to int
	mesure64, err := strconv.ParseInt(mesureIdStr, 10, 8)
	if err != nil {
		return errors.New(unknownMesureMsg)
	}
	mesure := int8(mesure64)
	if _, ok := Mesures[mesure]; !ok {
		return errors.New(unknownMesureMsg)
	}
	p.Mesure = mesure
	return nil
}

func (p *Produce) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = p.ID
	m["name"] = p.Name
	m["mesure_id"] = p.Mesure
	return m
}

func GetProduceFromMap(input map[string]interface{}) (*Produce, map[string]string) {
	var produce Produce
	errors := map[string]string{}
	if produceId, ok := input["id"]; ok {
		produce.ID = produceId.(uint)
	}

	if _, ok := input["name"]; ok {
		err := produce.setName(input["name"].(string))
		if err != nil {
			errors["name"] = err.Error()
		}
	} else {
		errors["name"] = "Parameter 'name' is required"
	}

	if _, ok := input["mesure"]; ok {
		err := produce.setMesure(input["mesure"].(string))
		if err != nil {
			errors["mesure"] = err.Error()
		}
	} else {
		errors["mesure"] = "Parameter 'mesure' is required"
	}

	if len(errors) > 0 {
		return nil, errors
	} else {
		return &produce, errors
	}
}
