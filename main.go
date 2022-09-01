//1. Get Startweb <Napp.Init> <sClientID> 654887908
//Result URL: http://racing.natsoft.com.au/654887908/object_585327.83w/Result?1

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var drivers = []string{
	"J.Craig",
	"H.Inwood",
	"J.Craig",
	"D.Smith",
	"W.Foot",
	"B.Stinson",
	"C.Manning",
	"G.Dufficy",
	"I.Ashcroft",
	"J.Canellis",
	"S.Tate",
	"B.Woodland",
	"S.Thompson",
	"J.Lawrence",
	"M.Miller",
	"C.Fraser",
	"B.Sheedy",
	"R.Albronda",
	"S.Tidyman",
	"B.Pronesti",
	"J.Pirozzi",
	"P.Halfpenny",
	"G.Manning",
	"C.Butterfield ",
	"A.Leacy",
	"C.Viola",
	"T.Maynard",
	"D.Agathos",
	"R.Baskus",
	"G.Reynolds",
	"K.Alderton",
	"I.Joyce",
	"M.Ricketts",
	"M.Early",
	"D.Barbaro",
	"A.Cortes",
	"T.Solly",
	"B.Wilson",
	"K.Avramidis",
	"J.Cox",
	"G.Oliver",
	"T.Worton",
	"L.Lichtenberger",
	"M.Boylan",
	"M.Kiss",
	"A.Mcdonald",
	"S.Doorey",
	"K.Albornda",
	"S.Kendrick",
	"C.Fraser",
	"P.Cannon",
	"M.Crutcher",
	"A.Lawrence",
	"J.Broad",
}

func main() {
	for {
		fmt.Print("Result Page URL:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		url := input.Text()

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(""))
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		body := string(b)

		eventHeaderStart := strings.Index(body, "<H4>")
		eventHeaderEnd := strings.Index(body[eventHeaderStart+4:], "</H4>")

		title := ""
		eventDetails := strings.Split(body[eventHeaderStart+4:eventHeaderStart+eventHeaderEnd-1], "<BR>")
		for _, detail := range eventDetails {
			detail = strings.TrimSpace(detail)
			if len(detail) <= 0 {
				continue
			}

			title = fmt.Sprintf("%s %s", title, detail)
		}
		title = strings.TrimSpace(title)

		headerLineIdx := strings.Index(body, "Pos Car")
		body = body[headerLineIdx:]

		standings := map[string]int{}

		lines := strings.Split(body, "\n")
		for i, line := range lines {
			if i < 2 {
				continue
			}

			if len(line) <= 0 {
				break
			}

			pos, _ := strconv.Atoi(strings.TrimSpace(line[:4]))
			//car, _ := strconv.Atoi(strings.TrimSpace(line[5:9]))
			driver := strings.TrimSpace(line[9:34])
			comp := strings.Split(driver, " ")

			id := fmt.Sprintf("%s.%s", string(comp[0][0]), comp[1])

			standings[id] = pos
		}

		fmt.Printf("%v drivers found\n", len(standings))

		filename := fmt.Sprintf("%s.csv", title)

		f, err := os.Create(filename)
		defer f.Close()

		for _, driver := range drivers {
			pos := standings[driver]
			if pos <= 0 {
				f.WriteString(fmt.Sprintf("%s\tDNF\n", driver))
			} else {
				f.WriteString(fmt.Sprintf("%s\t%v\n", driver, pos))
			}

			delete(standings, driver)
		}

		if len(standings) > 0 {
			fmt.Printf("%v new drivers\n", len(standings))
			for driver, pos := range standings {
				if pos <= 0 {
					f.WriteString(fmt.Sprintf("%s\tDNF\n", driver))
				} else {
					f.WriteString(fmt.Sprintf("%s\t%v\n", driver, pos))
				}
			}
		}

		fmt.Printf("Results in %s\n", filename)

		fmt.Printf("================================================================================\n")
	}
}
