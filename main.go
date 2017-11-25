package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"strings"
)

func emailSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	//una linea vuota e "From " dividono i messaggi
	//vedi https://en.wikipedia.org/wiki/Mbox#Family
	if i := strings.Index(string(data), "\n\nFrom "); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return
}

func readEmail(b []byte) {
	// per leggere una mail sono da rimuovere
	// le righe nuoe e il "From "
	const NL = "\n"
	trimmed := strings.TrimLeft(string(b), NL)
	var msgString string
	if strings.Index(trimmed, "From ") == 0 {
		msgString = strings.Join(strings.Split(trimmed, NL)[1:], NL)
	} else {
		msgString = trimmed
	}

	msg, err := mail.ReadMessage(strings.NewReader(msgString))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Subject:", msg.Header.Get("Subject"))
	//stampa su CSV i dati
	var data = []string{msg.Header.Get("Subject")}
	csvWriter(data)
}

func emailScanner(mbox io.Reader) {
	s := bufio.NewScanner(mbox)

	var (
		msg   []byte
		count int
	)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "From ") {
			if msg == nil {

			} else {
				count++
				readEmail(msg)
				msg = nil
			}
		} else {
			msg = append(msg, []byte("\n")...)
			msg = append(msg, s.Bytes()...)
		}
	}
	count++
	readEmail(msg)

	fmt.Println("Total emails:", count)
}

func csvWriter(yourSliceGoesHere []string) {
	f, err := os.OpenFile("data.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	w := csv.NewWriter(f)
	w.Write(yourSliceGoesHere)
	w.Flush()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage:", os.Args[0], "<filename>")
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Unable to open file:", err)
	}
	defer f.Close()

	emailScanner(f)
}
