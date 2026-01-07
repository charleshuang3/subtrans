package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/charleshuang3/subtrans/pkg/sub"
	"github.com/charleshuang3/subtrans/pkg/translator"
)

func parseFromIndex(fromIndex string) (int, int, int, error) {
	parts := strings.Split(fromIndex, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("--from must be in format item,line,seg")
	}
	var fromItem, fromLine, fromSeg int
	var err error
	fromItem, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error parsing item index: %w", err)
	}
	fromLine, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error parsing line index: %w", err)
	}
	fromSeg, err = strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error parsing seg index: %w", err)
	}
	return fromItem, fromLine, fromSeg, nil
}

func main() {
	inputFile := flag.String("i", "", "input file path (required)")
	outputFile := flag.String("o", "", "output file path (required)")
	targetLang := flag.String("target-lang", "", "target language (optional)")
	configPath := flag.String("c", "", "config file path (optional)")
	fromIndex := flag.String("from", "", "resume from index (item,line,seg)")
	promptKey := flag.String("prompt", "default", "prompt key from config (optional)")
	llmProvider := flag.String("llm", "default", "LLM provider to use (optional)")
	dryRun := flag.Bool("dry-run", false, "dry run without making API calls (optional)")
	flag.Parse()

	if *inputFile == "" {
		log.Fatalf("Error: -i (input file) is required")
	}
	if *outputFile == "" {
		log.Fatalf("Error: -o (output file) is required")
	}

	log.Printf("input file: %s", *inputFile)
	log.Printf("output file: %s", *outputFile)

	confPath, err := config.FindConfig(*configPath)
	if err != nil {
		log.Fatalf("Error finding config file: %v", err)
	}
	log.Printf("config file: %s", confPath)

	cfg, err := config.Read(confPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if *targetLang != "" {
		// overwrite target language
		cfg.TargetLang = *targetLang
	}
	log.Printf("target lang: %s", cfg.TargetLang)
	log.Printf("LLM provider: %s", *llmProvider)

	var fromItem, fromLine, fromSeg int
	if *fromIndex != "" {
		var err error
		fromItem, fromLine, fromSeg, err = parseFromIndex(*fromIndex)
		if err != nil {
			log.Fatalf("Error parsing from index: %v", err)
		}
		log.Printf("resuming from index: %d,%d,%d", fromItem, fromLine, fromSeg)
	}

	translator, err := translator.NewLLMTranslator(cfg, *promptKey, *llmProvider, *dryRun)
	if err != nil {
		log.Fatalf("Error creating translator: %v", err)
	}

	log.Printf("dry run: %t", *dryRun)

	if *fromIndex != "" {
		err = sub.TranslateFileFromIndex(*inputFile, *outputFile, translator, fromItem, fromLine, fromSeg)
	} else {
		err = sub.TranslateFile(*inputFile, *outputFile, translator)
	}
	if err != nil {
		log.Fatalf("Error translating file: %v", err)
	}
	log.Printf("Translation completed")
}
