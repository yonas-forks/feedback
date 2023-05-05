package router

import (
	"encoding/json"
	"os"

	"github.com/andrewarrow/feedback/files"
)

func (c *Context) ReadJsonBodyIntoParams() {
	c.Params = map[string]any{}
	body := c.BodyAsString()
	//fmt.Println(body)
	json.Unmarshal([]byte(body), &c.Params)
}

func (c *Context) ReadJsonBodyAsArray() []any {
	var list []any
	body := c.BodyAsString()
	json.Unmarshal([]byte(body), &list)
	return list
}

func (c *Context) ExecuteTemplate(filename string, vars any) {
	c.router.Template.ExecuteTemplate(c.Writer, filename, vars)
}

func (c *Context) ReadJsonBodyIntoParamsWithLog(file string) {
	home := files.UserHomeDir()
	filename := home + "/" + file
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	c.Params = map[string]any{}
	body := c.BodyAsString()

	defer f.Close()
	f.WriteString(body)
	json.Unmarshal([]byte(body), &c.Params)
}
