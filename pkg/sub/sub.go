package sub

import (
	"log"

	"github.com/asticode/go-astisub"
)

const BatchSize = 50

type Translator interface {
	Translate(texts []string) ([]string, error)
}

type textPosition struct {
	itemIndex int
	lineIndex int
	segIndex  int
}

func TranslateFile(inputPath, outputPath string, translator Translator) error {
	subs, err := astisub.OpenFile(inputPath)
	if err != nil {
		return err
	}

	positions := []textPosition{}
	texts := []string{}

	for itemIndex, item := range subs.Items {
		for lineIndex, line := range item.Lines {
			for segIndex, seg := range line.Items {
				if seg.Text != "" {
					texts = append(texts, seg.Text)
					positions = append(positions, textPosition{
						itemIndex: itemIndex,
						lineIndex: lineIndex,
						segIndex:  segIndex,
					})
				}
			}
		}
	}

	allTranslations := []string{}
	totalBatches := (len(texts) + BatchSize - 1) / BatchSize // Round up
	for i := 0; i < len(texts); i += BatchSize {
		end := i + BatchSize
		if end > len(texts) {
			end = len(texts)
		}
		batch := texts[i:end]
		batchNumber := i/BatchSize + 1
		log.Printf("Translating batch %d/%d (items %d-%d of %d)", batchNumber, totalBatches, i+1, end, len(texts))
		translations, err := translator.Translate(batch)
		if err != nil {
			return err
		}
		allTranslations = append(allTranslations, translations...)
	}
	log.Printf("Translation completed: %d items translated", len(texts))

	for i, pos := range positions {
		subs.Items[pos.itemIndex].Lines[pos.lineIndex].Items[pos.segIndex].Text = allTranslations[i]
	}

	return subs.Write(outputPath)
}
