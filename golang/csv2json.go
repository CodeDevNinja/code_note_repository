package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// convertCSVtoJSON reads a CSV file, converts it to JSON format, and writes it to a new JSON file.
func convertCSVtoJSON(inputFilePath, outputFilePath string) error {
	// Open the CSV file
	csvFile, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer csvFile.Close()

	// Create the JSON file
	jsonFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer jsonFile.Close()

	writer := bufio.NewWriter(jsonFile)
	defer writer.Flush()

	// Read the CSV content
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %w", err)
	}

	// Assuming the first row is header
	header := records[0]
	for _, row := range records[1:] {
		rowMap := make(map[string]string)
		for colIdx, col := range row {
			rowMap[header[colIdx]] = col
		}

		// Convert to JSON
		jsonData, err := json.Marshal(rowMap)
		if err != nil {
			continue // Skip errors in individual rows
		}

		// Write JSON to file, each object as a new line
		writer.Write(jsonData)
		writer.WriteString("\n")
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program> <input_directory>")
		return
	}

	inputDir := os.Args[1]
	// "/Users/caidongzhu/DevProgram/tt_data_center/temp/data/tiktok/insight/adgroup/stat_day_country/every_30days/"
	outputDir := inputDir + "json/"

	// Read all files in the input directory
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		fmt.Println("Error reading input directory:", err)
		return
	}

	for _, fileDir := range files {
		fmt.Println(fileDir.Name())
		csvFiles, err := ioutil.ReadDir(inputDir + fileDir.Name())
		if err != nil {
			fmt.Println("Error reading input directory:", err)
			return
		}
		if !fileDir.IsDir() {
			continue // Skip directories
		}
		for _, csvFile := range csvFiles {
			if csvFile.IsDir() {
				continue // Skip directories
			}
			if strings.HasSuffix(csvFile.Name(), ".csv") {
				// Construct full paths for input and output
				inputFilePath := filepath.Join(inputDir, fileDir.Name(), csvFile.Name())
				outputFilePath := filepath.Join(outputDir, fileDir.Name(), strings.TrimSuffix(csvFile.Name(), ".csv")+".json")
				// Ensure the output directory exists
				if err := os.MkdirAll(filepath.Join(outputDir, fileDir.Name()), os.ModePerm); err != nil {
					fmt.Println("Error creating output directory:", err)
					return
				}
				// Convert the CSV file to a JSON file
				if err := convertCSVtoJSON(inputFilePath, outputFilePath); err != nil {
					fmt.Println("Error converting file:", err)
					continue
				}
			}
		}

	}
}
