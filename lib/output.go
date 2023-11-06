package lib

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jszwec/csvutil"
	"gopkg.in/yaml.v3"
	"os"
)

func OutputYAML(d interface{}) error {
	yamlEnc := yaml.NewEncoder(os.Stdout)
	return yamlEnc.Encode(d)
}

func OutputJSON(d interface{}) error {
	jsonEnc := json.NewEncoder(os.Stdout)
	jsonEnc.SetIndent("", "  ")
	return jsonEnc.Encode(d)
}

func OutputCSV(v interface{}) error {
	csvWriter := csv.NewWriter(os.Stdout)
	csvEnc := csvutil.NewEncoder(csvWriter)

	if err := csvEnc.Encode(v); err != nil {
		return err
	}
	csvWriter.Flush()

	return nil
}
