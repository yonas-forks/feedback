package util

import (
	"strings"
)

func SpecialToSingle() map[string]string {
	return map[string]string{"series": "series",
		"bases":    "base",
		"coaches":  "coach",
		"dies":     "die",
		"children": "child",
		"species":  "species"}
}
func SpecialToPlural() map[string]string {
	return map[string]string{"series": "series",
		"child":   "children",
		"coach":   "coaches",
		"species": "species"}
}

func Plural(s string) string {
	tokens := strings.Split(s, "_")
	if len(tokens) == 1 {
		return PluralLogic(s)
	} else {
		front := strings.Join(tokens[0:len(tokens)-1], "_")
		return front + "_" + PluralLogic(tokens[len(tokens)-1])
	}
}

func PluralLogic(s string) string {
	m := SpecialToPlural()
	if m[s] != "" {
		return m[s]
	}
	if strings.HasSuffix(s, "y") {
		return strings.TrimSuffix(s, "y") + "ies"
	}
	return s + "s"
}

func Unplural(s string) string {
	tokens := strings.Split(s, "_")
	if len(tokens) == 1 {
		return UnpluralLogic(s)
	} else {
		front := strings.Join(tokens[0:len(tokens)-1], "_")
		return front + "_" + UnpluralLogic(tokens[len(tokens)-1])
	}
}

func UnpluralLogic(s string) string {
	m := SpecialToSingle()
	if m[s] != "" {
		return m[s]
	}
	return strings.TrimSuffix(s, "s")
}
