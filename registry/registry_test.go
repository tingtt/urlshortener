package registry

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
)

type RegistryTest struct {
	caseName string

	rawData        string
	invalidRawData bool
	data           map[string]string
}

var registrytests = []RegistryTest{
	{
		caseName:       "no shortened URLs",
		rawData:        "",
		invalidRawData: false,
		data:           map[string]string{},
	},
	{
		caseName: "some shortened URLs",
		rawData: strings.TrimPrefix(dedent.Dedent(`
			/github,https://github.com/tingtt/urlshortener
			/short1,https://example.com
			/short2,https://example.com
		`), "\n"),
		invalidRawData: false,
		data: map[string]string{
			"/github": "https://github.com/tingtt/urlshortener",
			"/short1": "https://example.com",
			"/short2": "https://example.com",
		},
	},
	{
		caseName: "invalid raw data",
		rawData: strings.TrimPrefix(dedent.Dedent(`
			invalid,https://example.com
		`), "\n"),
		invalidRawData: true,
		data:           nil,
	},
}

func Test_registry_savePersistently(t *testing.T) {
	t.Parallel()

	for _, tt := range registrytests {
		if tt.invalidRawData {
			continue
		}

		t.Run("save to file", func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			saveFilePath := path.Join(dir, "save.csv")

			r := &registry{tt.data, saveFilePath}
			err := r.savePersistently()

			assert.NoError(t, err)

			savedRawData, err := os.ReadFile(saveFilePath)
			if err != nil {
				t.Fatal("failed to read file: " + err.Error())
			}
			assert.Equal(t, tt.rawData, string(savedRawData))
		})
	}
}
