package gounistudium

import (
	"fmt"
	"strings"
	"testing"

	"github.com/antchfx/htmlquery"
)

func TestRequest(t *testing.T) {
	baseUrl := "https://unistudium.unipg.it/cercacorso.php?%d"

	roomType := ClassRoom
	query := "ingegneria"

	formattedUrl := fmt.Sprintf(baseUrl, roomType)
	payload := strings.NewReader(fmt.Sprintf("query=%s", query))

	_, err := Request(formattedUrl, payload)
	if err != nil {
		t.Fail()
	}
}

func TestParseRoom(t *testing.T) {
	input := "<html><body><table><tbody><tr><td>Corso</td><td>Professore</td><td>Classe</td><td>Link</td></tr></tbody></table></body></html>"

	doc, err := htmlquery.Parse(strings.NewReader(input))
	if err != nil {
		t.Fail()
	}

	trNodes, err := htmlquery.QueryAll(doc, "//tr")
	if err != nil {
		t.Fail()
	}

	for i := 0; i < len(trNodes); i++ {
		_, err := HtmlToRoom(trNodes[i], ClassRoom)
		if err != nil {
			t.Fail()
		}
	}
}
