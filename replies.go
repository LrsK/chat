package main

import "strings"

var greetings map[string]struct{}

func initDictionary() {
	greetings = make(map[string]struct{})
	greetings["hello"] = struct{}{}
	greetings["hi"] = struct{}{}
	greetings["hru"] = struct{}{}
	greetings["how are you"] = struct{}{}
}

func makeSentences(input string) []string {
	var prev rune
	snt := strings.FieldsFunc(input, func(r rune) bool {
		if r != ' ' {
			prev = r
			return false
		}
		switch prev {
		case '.', '?', '!':
			prev = r
			return true
		}

		prev = r
		return false
	})

	var sentences []string
	for _, s := range snt {
		sentences = append(sentences, strings.Trim(s, " \t\r\n"))
	}

	return sentences
}

func generateReply(history [][]byte) []byte {
	// Analyze input by
	// Look for question mark at last byte of last input
	// split on "." then " " remove other chars
	// for each word in each sentence, look for interesting words to save for later
	// If analysis is too hard, reply with a question.

	// Begin with last entry, and remove uppercase chars
	lastEntry := strings.ToLower(string(history[len(history)-1]))

	// Try to make sentences
	sentences := makeSentences(lastEntry)
	for _, sentence := range sentences {
		end := rune(sentence[len(sentence)-1])
		if end == '!' {
			return []byte("An imperative.")
		} else if end == '?' {
			return []byte("A question.")
		} else {
			return []byte("A normal sentence.")
		}
	}
	return []byte("Hello")
}
