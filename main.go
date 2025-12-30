package main

import (
	"flag"
	"log"

	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/charleshuang3/subtrans/pkg/sub"
	"github.com/charleshuang3/subtrans/pkg/translator"
)

func main() {
	inputFile := flag.String("i", "", "input file path (required)")
	outputFile := flag.String("o", "", "output file path (required)")
	targetLang := flag.String("target-lang", "", "target language (optional)")
	configPath := flag.String("c", "", "config file path (optional)")
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

	translator := translator.NewLLMTranslator(cfg)
	err = sub.TranslateFile(*inputFile, *outputFile, translator)
	if err != nil {
		log.Fatalf("Error translating file: %v", err)
	}
	log.Printf("Translation completed")
}
