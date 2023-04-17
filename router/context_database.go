package router

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/andrewarrow/feedback/models"
	"github.com/xeonx/timeago"
)

func (c *Context) FindModel(name string) *models.Model {
	return c.router.FindModel(name)
}

func (r *Router) FindModel(name string) *models.Model {
	return r.Site.FindModel(name)
}

func CastFields(model *models.Model, m map[string]any) {
	if len(m) == 0 {
		return
	}
	for _, field := range model.Fields {
		if field.Flavor == "timestamp" && m[field.Name] != nil {
			tm := m[field.Name].(time.Time)
			ago := timeago.English.Format(tm)
			m[field.Name] = tm.Unix()
			m[field.Name+"_human"] = tm.Format(models.HUMAN)
			m[field.Name+"_ago"] = ago
		} else if field.Flavor == "int" && m[field.Name] != nil {
			m[field.Name] = m[field.Name].(int64)
		} else if field.Flavor == "bool" && m[field.Name] != nil {
			m[field.Name] = m[field.Name].(bool)
		} else if field.Flavor == "json" && m[field.Name] != nil {
			var temp map[string]any
			json.Unmarshal([]byte(m[field.Name].(string)), &temp)
			m[field.Name] = temp
		} else if field.Flavor == "json_list" && m[field.Name] != nil {
			var temp []any
			json.Unmarshal([]byte(m[field.Name].(string)), &temp)
			m[field.Name] = temp
		} else if field.Flavor == "list" {
			s := fmt.Sprintf("%s", m[field.Name])
			tokens := strings.Split(s, ",")
			m[field.Name] = tokens
		} else if m[field.Name] == nil {
			// to nothing, leave it nil
		} else {
			m[field.Name] = fmt.Sprintf("%s", m[field.Name])
		}
	}
}
