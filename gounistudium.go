package gounistudium

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

type roomType int

const (
	ClassRoom      roomType = 0
	ExamRoom       roomType = 1
	GraduationRoom roomType = 245
)

const baseUrl = "https://unistudium.unipg.it/cercacorso.php?p=%d"

type Room struct {
	CourseName     string
	Professor      string
	Degree         string
	MeetingLink    string
	Time           string
	GraduationCode string
	Department     string
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
		room, err := HtmlToRoom(trNodes[i], roomType)
		if err != nil {
			fmt.Println("Errore HtmlToRoom")
			return nil, err
		}

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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("errore ioutil.Readall")
		return "", err
	}

	return string(body), nil
}

// HtmlToRoom parses the HTML code into a struct of type Room
func HtmlToRoom(node *html.Node, roomType roomType) (Room, error) {
	tdNodes, err := htmlquery.QueryAll(node, "//td")
	if err != nil {
		return Room{}, err
	}

	if roomType == ExamRoom {
		return Room{
			CourseName:  htmlquery.InnerText(tdNodes[0]),
			Professor:   htmlquery.InnerText(tdNodes[1]),
			Degree:      htmlquery.InnerText(tdNodes[2]),
			Time:        htmlquery.InnerText(tdNodes[3]),
			MeetingLink: htmlquery.SelectAttr(tdNodes[4].FirstChild, "href"),
		}, nil
	}

	if roomType == GraduationRoom {
		a, err := xpath.Compile("/form/input[7]")
		if err != nil {
			fmt.Println("Errore xpath.Compile")
			return Room{}, err
		}

		s := htmlquery.SelectAttr(htmlquery.QuerySelector(tdNodes[3], a), "value")
		return Room{
			Department:     htmlquery.InnerText(tdNodes[0]),
			Degree:         htmlquery.InnerText(tdNodes[1]),
			Time:           htmlquery.InnerText(tdNodes[2]),
			GraduationCode: strings.Split(s, " ")[2],
		}, nil
	}

	return Room{
		CourseName:  htmlquery.InnerText(tdNodes[0]),
		Professor:   htmlquery.InnerText(tdNodes[1]),
		Degree:      htmlquery.InnerText(tdNodes[2]),
		MeetingLink: htmlquery.SelectAttr(tdNodes[3].FirstChild, "href"),
	}, nil
}
