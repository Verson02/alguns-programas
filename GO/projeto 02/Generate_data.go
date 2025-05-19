package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Dados fictícios
var firstNames = []string{"Alice", "Bob", "Charlie", "David", "Eve", "Patolino", "Nicolas ", "Gustavo "}
var lastNames = []string{"Smith", "Johnson", "Williams", "Jones", "Brown", "Librellato ", "Verson "}
var streets = []string{"Main St", "Elm St", "Maple Ave", "Oak St", "Pine St"}
var cities = []string{"Springfield", "Riverside", "Greenville", "Franklin", "Clinton", "chique chique ", "Parisnagua"}

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

func writeCSV(w http.ResponseWriter, data []map[string]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=data.csv")
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"Name", "Email", "Phone", "Address"}
	if err := writer.Write(header); err != nil {
		http.Error(w, "Error writing CSV", http.StatusInternalServerError)
		return
	}

	for _, record := range data {
		if err := writer.Write([]string{record["Name"], record["Email"], record["Phone"], record["Address"]}); err != nil {
			http.Error(w, "Error writing CSV", http.StatusInternalServerError)
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, data []map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment;filename=data.json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		http.Error(w, "Error writing JSON", http.StatusInternalServerError)
		return
	}
}

func writeXML(w http.ResponseWriter, data []map[string]string) {
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Content-Disposition", "attachment;filename=data.xml")
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(data); err != nil {
		http.Error(w, "Error writing XML", http.StatusInternalServerError)
		return
	}
}

func generateDataHandler(w http.ResponseWriter, r *http.Request) {
	// Recupera os parâmetros da URL
	numRecordsStr := r.URL.Query().Get("numRecords")
	format := r.URL.Query().Get("format")
	filename := r.URL.Query().Get("filename")

	// Valida o parâmetro numRecords
	numRecords, err := strconv.Atoi(numRecordsStr)
	if err != nil || numRecords < 1 {
		http.Error(w, "Invalid number of records", http.StatusBadRequest)
		return
	}

	// Valida o parâmetro format
	format = strings.ToLower(format)
	if format != "csv" && format != "json" && format != "xml" {
		http.Error(w, "Unsupported format", http.StatusBadRequest)
		return
	}

	// Valida e sanitiza o parâmetro filename
	if filename == "" {
		filename = "data"
	}
	filename = strings.ReplaceAll(filename, " ", "_") // Substitui espaços por underlines para evitar problemas de nome de arquivo

	// Gera os dados
	data := generateUserData(numRecords)

	// Configura o cabeçalho e escreve o arquivo
	switch format {
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
		writeCSV(w, data)
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", filename))
		writeJSON(w, data)
	case "xml":
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xml", filename))
		writeXML(w, data)
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/generate", generateDataHandler)

	port := ":8080"
	fmt.Printf("Starting server on port%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}