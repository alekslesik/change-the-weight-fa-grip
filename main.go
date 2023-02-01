package main

import (
	"bufio"
	// "flag"
	"fmt"
	"strings"

	// "fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"gopkg.in/ini.v1"
)

type inputs struct {
	delta     int
	operation string
}

func main() {

	var inputs inputs

	fmt.Println("Input OPERATION SYMBOL string (+. -, /, *)")
	fmt.Scan(&inputs.operation)

	fmt.Println("Input DELTA integer (not float)")
	fmt.Scan(&inputs.delta)

	// Cyclo2.ini
	cycloFile, err := ini.Load("Cyclo2.ini")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	// Create backup file
	// err = cycloFile.SaveTo("Cyclo2_backup.ini")
	// if err != nil {
	// 	log.Fatalf("Error when save backup file: %s", err)
	// }

	err = doOperationDelta(cycloFile, inputs.operation, inputs.delta)
	if err != nil {
		log.Fatal("ERR: Operation error")
	}

	// Write file
	err = cycloFile.SaveToIndent("buffer.ini", "")
	if err != nil {
		log.Fatalf("Error when save file: %s", err)
	}

	// Open the buffer file
	file, err := os.Open("buffer.ini")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	fileTosave, err := os.Create("Cyclo2_result.ini")
	if err != nil {
		log.Fatalf("Error when create file: %s", err)
	}
	defer fileTosave.Close()

	// read line by line
	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.Replace(line, "  ", "", -1)
		line = strings.Replace(line, "= ", "=", -1)
		line = strings.Replace(line, " =", "=", -1)
		line = strings.Replace(line, "=", " = ", -1)
		fileTosave.WriteString(line + "\n")
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	// close buffer file
	file.Close()

	//delete buffer file
	err = os.Remove("buffer.ini")
	if err != nil {
		log.Fatalf("Error when remove file: %s", err)
	}

	fmt.Println("DONE! See Cyclo2_result.ini")
}

// Do operation whith all tension keys "m_*=" in FAGRIP sections Cyclo2.ini file
func doOperationDelta(cycloFile *ini.File, operation string, delta int) error {
	// Loop every diagramm [section]
	for _, section := range cycloFile.Sections() {
		// Take only if diagramm name include "FAGRIP"
		re, err := regexp.Compile(`.*FAGRIP.*`)
		if err != nil {
			log.Fatalf("Open regexp section compile: %s", err)
			return err
		}

		if !re.Match([]byte(section.Name())) {
			continue
		}

		// Loop every tension key m_*=
		for _, sectionKey := range section.Keys() {
			// Take only if tension key name include m_*
			re, err := regexp.Compile(`m_\d*`)
			if err != nil {
				log.Fatalf("Open regexp key compile: %s", err)
				return err
			}

			if re.Match([]byte(sectionKey.Name())) || sectionKey.Name() == "StopTension" {
				// Get tension key value
				keyInt, err := sectionKey.Int()
				if err != nil {
					log.Fatalf("Error get key value: %s", err)
					return err
				}

				switch operation {
				case "+":
					// And plus to delta
					sectionKey.SetValue(strconv.Itoa(keyInt + delta))
				case "-":
					// And subtract to delta
					sectionKey.SetValue(strconv.Itoa(keyInt - delta))
				case "/":
					// And devide to delta
					sectionKey.SetValue(strconv.Itoa(keyInt / delta))
				case "*":
					// And multiple to delta
					sectionKey.SetValue(strconv.Itoa(keyInt * delta))
				}

			} else {
				continue
			}
		}
	}

	return nil
}