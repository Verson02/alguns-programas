package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Dados fictícios
var firstNames = []string{"Alice", "Nicolas", "Charlie", "David", "Eve"}
var lastNames = []string{"Smith", "Johnson", "Williams", "Jones", "Brown"}
var streets = []string{"Main St", "Elm St", "Maple Ave", "Oak St", "Pine St"}
var cities = []string{"Springfield", "Riverside", "Greenville", "Franklin", "Clinton"}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateFullName() string {
	return fmt.Sprintf("%s %s", firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))])
}

func generateEmail(name string) string {
	domains := []string{"example.com", "test.com", "sample.com"}
	return fmt.Sprintf("%s@%s", name, domains[rand.Intn(len(domains))])
}

func generatePhoneNumber() string {
	return fmt.Sprintf("+1-%03d-%03d-%04d", rand.Intn(1000), rand.Intn(1000), rand.Intn(10000))
}

func generateAddress() string {
	return fmt.Sprintf("%d %s, %s, USA", rand.Intn(1000), streets[rand.Intn(len(streets))], cities[rand.Intn(len(cities))])
}

// Atualizar para remover a data de nascimento
func generateUserData(numRecords int) []map[string]string {
	var data []map[string]string
	for i := 0; i < numRecords; i++ {
		name := generateFullName()
		email := generateEmail(name)
		phone := generatePhoneNumber()
		address := generateAddress()
		data = append(data, map[string]string{
			"Name":    name,
			"Email":   email,
			"Phone":   phone,
			"Address": address,
		})
	}
	return data
}

func writeCSV(filename string, data []map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Remover 'DOB' do cabeçalho
	header := []string{"Name", "Email", "Phone", "Address"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range data {
		if err := writer.Write([]string{record["Name"], record["Email"], record["Phone"], record["Address"]}); err != nil {
			return err
		}
	}

	return nil
}

func writeJSON(filename string, data []map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func writeXML(filename string, data []map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	return encoder.Encode(data)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numRecords := flag.Int("n", 1000, "Number of records to generate")
	format := flag.String("format", "csv", "Output file format (csv, json, xml)")
	filename := flag.String("f", "test_data", "Output file name without extension")
	flag.Parse()

	data := generateUserData(*numRecords)

	var err error
	switch *format {
	case "csv":
		err = writeCSV(*filename+".csv", data)
	case "json":
		err = writeJSON(*filename+".json", data)
	case "xml":
		err = writeXML(*filename+".xml", data)
	default:
		fmt.Println("Unsupported format. Please use csv, json, or xml.")
		return
	}

	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Data successfully written to %s.%s\n", *filename, *format)
}
