package jpugdoc

import (
	"fmt"
	"log"
	"strings"

	"github.com/jwalton/gchalk"
)

// CheckWord は英単語に対して日本語の単語が対になっているかをチェックします。
func CheckWord(en string, ja string, fileNames []string) {
	vTag, err := versionTag()
	if err != nil {
		log.Fatal(err)
	}
	for _, fileName := range fileNames {
		src := getDiff(vTag, fileName)
		pairs := Extraction(src)
		for _, pair := range pairs {
			if strings.Contains(pair.en, en) {
				if !strings.Contains(pair.ja, ja) {
					fmt.Printf("%s not include word)\n", fileName)
					fmt.Println(gchalk.Green(pair.en))
					fmt.Println(pair.ja)
				}
			}
		}

	}
}
