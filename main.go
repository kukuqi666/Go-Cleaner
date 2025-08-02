package main

import (
	"encoding/json"
	"os"

	"Go-Cleaner/cleaner"
	"Go-Cleaner/gui"
)

/*
  Go-Cleaner
  Author: kukuqi666
  Email: kukuqi666@gmail.com
  Version: 1.0.0
*/

func main() {
	file, _ := os.ReadFile("rules/default_rules.json")
	var rules []cleaner.Rule
	json.Unmarshal(file, &rules)

	gui.RunUI(rules)
}
