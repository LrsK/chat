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

func makeSentences(input string) (sentences []string) {
	snt := strings.FieldsFunc(input, func(r rune) bool {
		switch r {
		case '.', '?', '!':
			return true
		}
		return false
	})

	for _, s := range snt {
		sentences = append(sentences, strings.Trim(s, " \t\r\n"))
	}

	return
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
	if len(sentences) == 0 {
		// This is probably one long sentence
		sentences = append(sentences, lastEntry)
	}

	if _, ok := greetings[lastEntry]; ok {
		return []byte("Hi! How are you?")
	}

	return []byte("Hello")
}
