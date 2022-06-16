package gounistudium

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
)

type roomType int

const (
	ClassRoom      roomType = 0
	ExamRoom       roomType = 1
	GraduationRoom roomType = 245
)

const baseUrl = "https://unistudium.unipg.it/cercacorso.php?%d"

type Room struct {
	CourseName  string
	Professor   string
	Degree      string
	MeetingLink string
}

// FindRooms returns an array of type Room containing information about the rooms
// each room is defined by a Room struct
// it needs a roomType and a query such as the course name, the professor's name, the degree
// or any keyword related to the room that you are looking for.
func FindRooms(roomType roomType, query string) ([]Room, error) {
	var rooms []Room

	formattedUrl := fmt.Sprintf(baseUrl, roomType)
	payload := strings.NewReader(fmt.Sprintf("query=%s", query))

	body, err := Request(formattedUrl, payload)
	if err != nil {
		return []Room{}, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return []Room{}, err
	}

	trNodes, err := htmlquery.QueryAll(doc, "//tr")
	if err != nil {
		return []Room{}, err
	}

	// the first tr element in the table contains the headers, so we can ignore it
	for i := 1; i < len(trNodes); i++ {
		room, _ := HtmlToRoom(trNodes[i])
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// Request retrieves the data from the University's website
func Request(requestUrl string, payload *strings.Reader) (string, error) {
	req, err := http.NewRequest(http.MethodPost, requestUrl, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return string(body), nil
}

// HtmlToRoom parses the HTML code into a struct of type Room
func HtmlToRoom(node *html.Node) (Room, error) {
	tdNodes, err := htmlquery.QueryAll(node, "//td")
	if err != nil {
		return Room{}, err
	}

	return Room{
		CourseName:  htmlquery.InnerText(tdNodes[0]),
		Professor:   htmlquery.InnerText(tdNodes[1]),
		Degree:      htmlquery.InnerText(tdNodes[2]),
		MeetingLink: htmlquery.SelectAttr(tdNodes[3].FirstChild, "href"),
	}, nil
}
