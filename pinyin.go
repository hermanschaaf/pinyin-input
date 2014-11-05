package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hermanschaaf/cedict"
)

// extractTone splits the tone number and the pinyin syllable returning a string
// and an integer, e.g., dong1 => dong, 1
func extractTone(p string) (string, int) {
	tone := int(p[len(p)-1]) - 48

	if tone < 0 || tone > 5 {
		return p, 0
	}
	return p[0 : len(p)-1], tone
}

func loadCEDict(file io.Reader) (map[string]string, error) {
	pronDict := map[string]string{}

	c := cedict.New(file)
	for {
		err := c.NextEntry()
		if err == cedict.NoMoreEntries {
			return pronDict, nil
		} else if err != nil {
			return pronDict, err
		}

		r := []rune(c.Entry().Simplified)
		pinyin := c.Entry().Pinyin
		pinParts := strings.Split(pinyin, " ")
		// skip entries where pinyin does not match chars
		if len(pinParts) != len(r) {
			continue
		}
		for i := range pinParts {
			s, _ := extractTone(pinParts[i])
			s = strings.ToLower(s)
			if !strings.Contains(pronDict[s], string(r[i])) {
				pronDict[s] += string(r[i])
			}
		}
	}
	return pronDict, nil
}

func main() {
	file, err := os.Open("cedict_ts.u8")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	d, err := loadCEDict(file)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("sounds.json", b, 0644)
	if err != nil {
		panic(err)
	}
}
