package sqlgen

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/andrewarrow/feedback/models"
)

func UpdateRow(model *models.Model) string {
	tableName := model.TableName()
	buffer := []string{"UPDATE "}
	buffer = append(buffer, tableName+" set ")

	cols := []string{}
	for i, field := range model.Fields {
		cols = append(cols, fmt.Sprintf("%s=$%d", field.Name, i+1))
	}
	buffer = append(buffer, strings.Join(cols, ","))
	buffer = append(buffer, fmt.Sprintf(" where guid=$%d", len(model.Fields)+1))

	return strings.Join(buffer, "")
}

func InsertRowNoRandomDefaults(tableName string,
	fields []*models.Field,
	override map[string]any) (string, []any) {
	return insertRow(false, tableName, fields, override)
}

func InsertRowWithRandomDefaults(tableName string,
	fields []*models.Field,
	override map[string]any) (string, []any) {
	return insertRow(true, tableName, fields, override)
}

func insertRow(random bool, tableName string,
	fields []*models.Field,
	override map[string]any) (string, []any) {

	buffer := []string{"INSERT INTO "}
	buffer = append(buffer, tableName+" (")

	cols := []string{}
	for _, field := range fields {
		if field.Name == "id" || field.Name == "created_at" || field.Name == "updated_at" {
			continue
		}
		cols = append(cols, field.Name)
	}
	buffer = append(buffer, strings.Join(cols, ","))
	buffer = append(buffer, ") values (")
	cols = []string{}
	params := []any{}
	count := 1
	for _, field := range fields {
		if field.Name == "id" || field.Name == "created_at" || field.Name == "updated_at" {
			continue
		}
		cols = append(cols, fmt.Sprintf("$%d", count))
		count++
		val := override[field.Name]
		if val == nil {
			if random {
				val = field.RandomValue()
			} else {
				val = field.Default()
			}
		}
		if field.Flavor == "list" && val != nil {
			list := []string{}
			thing1, isArrayAny := val.([]any)
			thing2, isArrayString := val.([]string)
			thing3, isString := val.(string)

			if isString {
				list = append(list, strings.ToLower(thing3))
			} else if isArrayAny {
				for _, s := range thing1 {
					list = append(list, strings.ToLower(s.(string)))
				}
			} else if isArrayString {
				for _, s := range thing2 {
					list = append(list, strings.ToLower(s))
				}
			}
			val = strings.Join(list, ",")
		} else if field.Flavor == "json" {
			asBytes, _ := json.Marshal(val)
			val = string(asBytes)
		} else if field.Flavor == "json_list" {
			asBytes, _ := json.Marshal(val)
			val = string(asBytes)
		}
		params = append(params, val)
	}
	buffer = append(buffer, strings.Join(cols, ","))
	buffer = append(buffer, ")")

	return strings.Join(buffer, ""), params
}

func UpdateRowFromParams(tableName string,
	fields []*models.Field,
	override map[string]any, where string) (string, []any) {

	params := []any{}
	buffer := []string{"UPDATE "}
	buffer = append(buffer, tableName+" set ")

	cols := []string{}
	count := 1
	for _, field := range fields {
		if field.Name == "id" || field.Name == "created_at" {
			continue
		}
		cols = append(cols, fmt.Sprintf("%s=$%d", field.Name, count))
		count++
		val := override[field.Name]
		if field.Flavor == "list" {
			list := fixListItems(val)
			val = strings.Join(list, ",")
		} else if field.Flavor == "json" {
			asBytes, _ := json.Marshal(val)
			val = string(asBytes)
		} else if field.Flavor == "json_list" {
			asBytes, _ := json.Marshal(val)
			val = string(asBytes)
		}
		params = append(params, val)
	}
	cols = append(cols, fmt.Sprintf("updated_at=$%d", count))
	params = append(params, time.Now())

	buffer = append(buffer, strings.Join(cols, ","))
	buffer = append(buffer, fmt.Sprintf(" %s$%d", where, count+1))
	return strings.Join(buffer, ""), params
}

func fixListItems(val any) []string {
	s, ok := val.(string)
	if ok {
		return []string{s}
	}
	list := []string{}
	items, ok := val.([]any)
	if ok {
		for _, s := range items {
			list = append(list, strings.ToLower(s.(string)))
		}
		return list
	}
	return val.([]string)
}
