package sub

import (
	"fmt"
	"log"

	"github.com/asticode/go-astisub"
)

type TranslationError struct {
	BatchNumber    int
	CompletedItems int
	FirstFailed    textInfo
	Err            error
}

func (e *TranslationError) Error() string {
	return fmt.Sprintf("batch %d failed: %v (completed %d items, first failed at %d,%d,%d)",
		e.BatchNumber, e.Err, e.CompletedItems, e.FirstFailed.itemIndex, e.FirstFailed.lineIndex, e.FirstFailed.segIndex)
}

func (e *TranslationError) Unwrap() error {
	return e.Err
}

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

func extractInfos(subs *astisub.Subtitles, translator Translator) []textInfo {
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
	return infos
}

func findOffset(infos []textInfo, fromItem, fromLine, fromSeg int) (int, error) {
	for i, info := range infos {
		if info.itemIndex == fromItem && info.lineIndex == fromLine && info.segIndex == fromSeg {
			return i, nil
		}
	}
	return 0, fmt.Errorf("specified index %d,%d,%d not found in input file", fromItem, fromLine, fromSeg)
}

func processBatches(subs *astisub.Subtitles, infos []textInfo, startingOffset int, globalCompleted int, translator Translator, outputPath string, partialLogMsg string) error {
	infosToProcess := infos[startingOffset:]
	batches := createBatches(infosToProcess, translator.MaxLength())

	log.Printf("total batches %d, limit length %d", len(batches), translator.MaxLength())
	currentOffset := 0
	for i, batch := range batches {
		log.Printf("Translating batch %d (items %d, length %d)", i+1, len(batch), getBatchLength(batch))
		translations, err := translator.Translate(batch)
		if err != nil {
			if currentOffset > 0 {
				writeErr := subs.Write(outputPath)
				if writeErr != nil {
					log.Printf("Warning: failed to write partial translation: %v", writeErr)
				} else {
					log.Printf(partialLogMsg, currentOffset)
				}
			}
			return &TranslationError{
				BatchNumber:    i + 1,
				CompletedItems: globalCompleted + currentOffset,
				FirstFailed:    infos[startingOffset+currentOffset],
				Err:            err,
			}
		}
		for j := 0; j < len(batch); j++ {
			info := infosToProcess[currentOffset+j]
			subs.Items[info.itemIndex].Lines[info.lineIndex].Items[info.segIndex].Text = translations[j]
		}
		currentOffset += len(batch)
	}
	log.Printf("Translation completed: %d items translated", len(infosToProcess))

	return subs.Write(outputPath)
}

func TranslateFile(inputPath, outputPath string, translator Translator) error {
	subs, err := astisub.OpenFile(inputPath)
	if err != nil {
		return err
	}

	infos := extractInfos(subs, translator)
	return processBatches(subs, infos, 0, 0, translator, outputPath, "Wrote partial translation with %d completed items")
}

func TranslateFileFromIndex(inputPath, outputPath string, translator Translator, fromItem, fromLine, fromSeg int) error {
	subs, err := astisub.OpenFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file for resuming: %w", err)
	}

	inputSubs, err := astisub.OpenFile(inputPath)
	if err != nil {
		return err
	}

	infos := extractInfos(inputSubs, translator)

	offset, err := findOffset(infos, fromItem, fromLine, fromSeg)
	if err != nil {
		return err
	}

	log.Printf("Resuming translation from item %d (offset %d)", offset, offset)

	return processBatches(subs, infos, offset, offset, translator, outputPath, "Wrote partial translation with %d additional completed items")
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
