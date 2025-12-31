package sub

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTranslator struct {
	translations map[string]string
	maxLength    int
}

func (m *mockTranslator) Length(text string) int {
	return 1
}

func (m *mockTranslator) MaxLength() int {
	return m.maxLength
}

func (m *mockTranslator) Translate(texts []string) ([]string, error) {
	result := make([]string, len(texts))
	for i, text := range texts {
		if trans, ok := m.translations[text]; ok {
			result[i] = trans
		} else {
			result[i] = text
		}
	}
	return result, nil
}

func TestTranslateFile(t *testing.T) {
	inputContent := `1
00:00:01,000 --> 00:00:04,000
Hello world

2
00:00:05,000 --> 00:00:08,000
How are you?
`

	expected := "\ufeff" + `1
00:00:01,000 --> 00:00:04,000
Hola mundo

2
00:00:05,000 --> 00:00:08,000
¿Cómo estás?
`

	tmpDir := t.TempDir()
	tmpInput := filepath.Join(tmpDir, "input.srt")
	tmpOutput := filepath.Join(tmpDir, "output.srt")

	err := os.WriteFile(tmpInput, []byte(inputContent), 0644)
	assert.NoError(t, err)

	translator := &mockTranslator{
		translations: map[string]string{
			"Hello world":  "Hola mundo",
			"How are you?": "¿Cómo estás?",
		},
		maxLength: 10,
	}

	err = TranslateFile(tmpInput, tmpOutput, translator)
	assert.NoError(t, err)

	outputContent, err := os.ReadFile(tmpOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(outputContent))
}

func TestTranslateFileBatch(t *testing.T) {
	inputContent := `1
00:00:01,000 --> 00:00:04,000
Line 1

2
00:00:05,000 --> 00:00:08,000
Line 2

3
00:00:09,000 --> 00:00:12,000
Line 3

4
00:00:13,000 --> 00:00:16,000
Line 4
`

	expected := "\ufeff" + `1
00:00:01,000 --> 00:00:04,000
Trans 1

2
00:00:05,000 --> 00:00:08,000
Trans 2

3
00:00:09,000 --> 00:00:12,000
Trans 3

4
00:00:13,000 --> 00:00:16,000
Trans 4
`

	tmpDir := t.TempDir()
	tmpInput := filepath.Join(tmpDir, "input.srt")
	tmpOutput := filepath.Join(tmpDir, "output.srt")

	err := os.WriteFile(tmpInput, []byte(inputContent), 0644)
	assert.NoError(t, err)

	translator := &mockTranslator{
		translations: map[string]string{
			"Line 1": "Trans 1",
			"Line 2": "Trans 2",
			"Line 3": "Trans 3",
			"Line 4": "Trans 4",
		},
		// so translator will called twice.
		maxLength: 2,
	}

	err = TranslateFile(tmpInput, tmpOutput, translator)
	assert.NoError(t, err)

	outputContent, err := os.ReadFile(tmpOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(outputContent))
}

func TestTranslateFileMultiSegments(t *testing.T) {
	inputContent := `1
00:00:01,000 --> 00:00:04,000
<b>Hello</b> <i>world</i>

2
00:00:05,000 --> 00:00:08,000
<font color="red">How</font> <font color="blue">are</font> you?
`

	expected := "\ufeff" + `1
00:00:01,000 --> 00:00:04,000
<b>Hola</b><i>mundo</i>

2
00:00:05,000 --> 00:00:08,000
<font color="red">Cómo</font><font color="blue">estás</font>¿tú?
`

	tmpDir := t.TempDir()
	tmpInput := filepath.Join(tmpDir, "input.srt")
	tmpOutput := filepath.Join(tmpDir, "output.srt")

	err := os.WriteFile(tmpInput, []byte(inputContent), 0644)
	assert.NoError(t, err)

	translator := &mockTranslator{
		translations: map[string]string{
			"Hello": "Hola",
			"world": "mundo",
			"How":   "Cómo",
			"are":   "estás",
			" you?": "¿tú?",
		},
		maxLength: 10,
	}

	err = TranslateFile(tmpInput, tmpOutput, translator)
	assert.NoError(t, err)

	outputContent, err := os.ReadFile(tmpOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(outputContent))
}
