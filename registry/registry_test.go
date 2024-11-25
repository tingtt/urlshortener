package registry

import "github.com/lithammer/dedent"

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
		rawData: dedent.Dedent(`
			/short1,https://example.com
			/short2,https://example.com
			/github,https://github.com/tingtt/urlshortener
		`),
		invalidRawData: false,
		data: map[string]string{
			"/short1": "https://example.com",
			"/short2": "https://example.com",
			"/github": "https://github.com/tingtt/urlshortener",
		},
	},
	{
		caseName: "invalid raw data",
		rawData: dedent.Dedent(`
			invalid,https://example.com
		`),
		invalidRawData: true,
		data:           nil,
	},
}
