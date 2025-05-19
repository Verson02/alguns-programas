package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Estruturas base
type FuzzResult struct {
	Input     string
	Response  string
	Status    int
	Error     error
	Timestamp time.Time
	Severity  string
}

type RequestBody struct {
	Data string `json:"data"`
}

type FuzzerConfig struct {
	TargetURL     string
	MaxIterations int
	Timeout       time.Duration
	OutputDir     string
}

type Fuzzer struct {
	config  FuzzerConfig
	client  *http.Client
	results []FuzzResult
}

// Construtor do Fuzzer
func NewFuzzer(config FuzzerConfig) *Fuzzer {
	return &Fuzzer{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		results: make([]FuzzResult, 0),
	}
}

// Gerador de mutações
func (f *Fuzzer) generateMutation(iteration int) string {
	mutations := []string{
		"normal_input",
		"trigger",
		"trigger_crash",
		fmt.Sprintf("test_%d", iteration),
		"<script>alert(1)</script>",
		"' OR '1'='1",
		string(make([]byte, iteration%100)),
		"NULL",
		"undefined",
		"{\"malformed\":\"json\"",
		"../../etc/passwd",
		"SELECT * FROM users",
		"<img src=x onerror=alert(1)>",
		"%00%0A%0D",
	}
	return mutations[iteration%len(mutations)]
}

// Determinação de severidade
func (f *Fuzzer) determineSeverity(result FuzzResult) string {
	if result.Error != nil || result.Status >= 500 {
		return "High"
	}
	if strings.Contains(strings.ToLower(result.Input), "script") ||
		strings.Contains(strings.ToLower(result.Input), "sql") {
		return "Medium"
	}
	return "Low"
}

// Fazer request
func (f *Fuzzer) makeRequest(input string) FuzzResult {
	body := RequestBody{Data: input}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return FuzzResult{
			Input:     input,
			Error:     err,
			Timestamp: time.Now(),
			Severity:  "High",
		}
	}

	resp, err := f.client.Post(
		f.config.TargetURL,
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return FuzzResult{
			Input:     input,
			Error:     err,
			Timestamp: time.Now(),
			Severity:  "High",
		}
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	result := FuzzResult{
		Input:     input,
		Response:  string(respBody),
		Status:    resp.StatusCode,
		Timestamp: time.Now(),
	}
	result.Severity = f.determineSeverity(result)
	return result
}

// Salvar crash individual
func (f *Fuzzer) saveCrash(result FuzzResult) error {
	filename := fmt.Sprintf("%s/crash_%d.txt", f.config.OutputDir, result.Timestamp.Unix())
	content := fmt.Sprintf("Input: %s\nResponse: %s\nStatus: %d\nSeverity: %s\nTimestamp: %s\n",
		result.Input, result.Response, result.Status, result.Severity, result.Timestamp.Format(time.RFC3339))

	return ioutil.WriteFile(filename, []byte(content), 0644)
}

// Executar fuzzing
func (f *Fuzzer) Run() {
	os.MkdirAll(f.config.OutputDir, 0755)

	for i := 0; i < f.config.MaxIterations; i++ {
		input := f.generateMutation(i)
		result := f.makeRequest(input)
		f.results = append(f.results, result)

		fmt.Printf("Iteration %d: Status %d, Severity: %s\n", i, result.Status, result.Severity)

		if result.Status >= 500 || result.Error != nil {
			fmt.Printf("Found potential crash with input: %s\n", input)
			f.saveCrash(result)
		}
	}
}

// Análise de crashes
func (f *Fuzzer) analyzeCrashes() {
	crashTypes := make(map[string]int)
	severityCount := make(map[string]int)

	for _, result := range f.results {
		if result.Status >= 500 || result.Error != nil {
			severityCount[result.Severity]++

			if strings.Contains(result.Input, "script") {
				crashTypes["XSS"]++
			} else if strings.Contains(result.Input, "SQL") {
				crashTypes["SQL Injection"]++
			} else if strings.Contains(result.Input, "trigger") {
				crashTypes["Trigger"]++
			} else {
				crashTypes["Other"]++
			}
		}
	}

	fmt.Println("\n=== Análise de Crashes ===")
	fmt.Printf("Total de testes: %d\n", len(f.results))
	fmt.Printf("Total de crashes: %d\n", len(crashTypes))

	fmt.Println("\nPor Severidade:")
	for sev, count := range severityCount {
		fmt.Printf("%s: %d\n", sev, count)
	}

	fmt.Println("\nPor Tipo:")
	for crashType, count := range crashTypes {
		fmt.Printf("%s: %d\n", crashType, count)
	}
}

// Gerar relatório HTML
func (f *Fuzzer) generateHTMLReport() error {
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Fuzzing Report</title>
        <style>
            body { 
                font-family: Arial, sans-serif; 
                margin: 20px;
                background-color: #f5f5f5;
            }
            .container {
                max-width: 1200px;
                margin: 0 auto;
            }
            .crash { 
                border: 1px solid #ddd; 
                margin: 10px 0;
                padding: 15px;
                border-radius: 5px;
                background-color: white;
            }
            .high { border-left: 5px solid #ff5252; }
            .medium { border-left: 5px solid #ffd740; }
            .low { border-left: 5px solid #69f0ae; }
            .stats {
                display: grid;
                grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
                gap: 20px;
                margin-bottom: 20px;
            }
            .stat-card {
                background: white;
                padding: 15px;
                border-radius: 5px;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Fuzzing Report</h1>
            <div class="stats">
                <div class="stat-card">
                    <h3>Total Tests</h3>
                    <p>%d</p>
                </div>
                <div class="stat-card">
                    <h3>Total Crashes</h3>
                    <p>%d</p>
                </div>
            </div>
            <h2>Crashes Details</h2>
            <div id="crashes">
    `

	crashCount := 0
	for _, result := range f.results {
		if result.Status >= 500 || result.Error != nil {
			crashCount++
			html += fmt.Sprintf(`
                <div class="crash %s">
                    <h3>Crash #%d</h3>
                    <p><strong>Input:</strong> %s</p>
                    <p><strong>Response:</strong> %s</p>
                    <p><strong>Status:</strong> %d</p>
                    <p><strong>Severity:</strong> %s</p>
                    <p><strong>Timestamp:</strong> %s</p>
                </div>
            `, strings.ToLower(result.Severity), crashCount, result.Input, result.Response, result.Status, result.Severity, result.Timestamp.Format(time.RFC3339))
		}
	}

	html += `
            </div>
        </div>
    </body>
    </html>
    `

	finalHTML := fmt.Sprintf(html, len(f.results), crashCount)
	return ioutil.WriteFile("fuzzing_report.html", []byte(finalHTML), 0644)
}

// Exportar para CSV
func (f *Fuzzer) exportToCSV() error {
	csvFile, err := os.Create("fuzzing_results.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Cabeçalho
	writer.Write([]string{"Timestamp", "Input", "Response", "Status", "Severity", "Error"})

	// Dados
	for _, result := range f.results {
		errorStr := ""
		if result.Error != nil {
			errorStr = result.Error.Error()
		}

		writer.Write([]string{
			result.Timestamp.Format(time.RFC3339),
			result.Input,
			result.Response,
			fmt.Sprintf("%d", result.Status),
			result.Severity,
			errorStr,
		})
	}

	return nil
}

func main() {
	config := FuzzerConfig{
		TargetURL:     "http://localhost:3000/parse",
		MaxIterations: 10,
		Timeout:       5 * time.Second,
		OutputDir:     "crashes",
	}

	fuzzer := NewFuzzer(config)

	fmt.Println("Iniciando fuzzing...")
	fuzzer.Run()

	fmt.Println("\nAnalisando resultados...")
	fuzzer.analyzeCrashes()

	fmt.Println("\nGerando relatórios...")
	if err := fuzzer.generateHTMLReport(); err != nil {
		fmt.Printf("Erro ao gerar relatório HTML: %v\n", err)
	}

	if err := fuzzer.exportToCSV(); err != nil {
		fmt.Printf("Erro ao exportar CSV: %v\n", err)
	}

	fmt.Println("\nRelatórios gerados:")
	fmt.Println("1. fuzzing_report.html - Relatório visual detalhado")
	fmt.Println("2. fuzzing_results.csv - Dados para análise em planilha")
	fmt.Println("3. Pasta 'crashes' - Detalhes individuais dos crashes")
}
