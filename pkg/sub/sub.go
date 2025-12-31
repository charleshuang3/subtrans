package sub

import (
	"log"

	"github.com/asticode/go-astisub"
)

type Translator interface {
	Translate(texts []string) ([]string, error)
	Length(text string) int
	MaxLength() int
}

type textInfo struct {
	itemIndex int
	lineIndex int
	segIndex  int
	length    int
	text      string
}

func TranslateFile(inputPath, outputPath string, translator Translator) error {
	subs, err := astisub.OpenFile(inputPath)
	if err != nil {
		return err
	}

	infos := []textInfo{}

	for itemIndex, item := range subs.Items {
		for lineIndex, line := range item.Lines {
			for segIndex, seg := range line.Items {
				if seg.Text == "" {
					continue
				}
				infos = append(infos, textInfo{
					itemIndex: itemIndex,
					lineIndex: lineIndex,
					segIndex:  segIndex,
					text:      seg.Text,
					length:    translator.Length(seg.Text),
				})
			}
		}
	}

	batches := createBatches(infos, translator.MaxLength())

	allTranslations := []string{}
	for i, batch := range batches {
		log.Printf("Translating batch %d (items %d, length %d)", i+1, len(batch), getBatchLength(batch))
		translations, err := translator.Translate(batch)
		if err != nil {
			return err
		}
		allTranslations = append(allTranslations, translations...)
	}
	log.Printf("Translation completed: %d items translated", len(infos))

	for i, info := range infos {
		subs.Items[info.itemIndex].Lines[info.lineIndex].Items[info.segIndex].Text = allTranslations[i]
	}

	return subs.Write(outputPath)
}

func createBatches(infos []textInfo, maxLength int) [][]string {
	batches := [][]string{}
	currentBatch := []string{}
	currentLength := 0

	for _, info := range infos {
		if currentLength+info.length > maxLength && len(currentBatch) > 0 {
			batches = append(batches, currentBatch)
			currentBatch = []string{}
			currentLength = 0
		}
		currentBatch = append(currentBatch, info.text)
		currentLength += info.length
	}

	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batches
}

func getBatchLength(batch []string) int {
	total := 0
	for _, s := range batch {
		total += len(s)
	}
	return total
}
