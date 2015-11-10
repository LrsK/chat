package main

var greetings map[string]struct{}

func initWords() {
	greetings = make(map[string]struct{})
	greetings["Hello"] = struct{}{}
	greetings["Hi"] = struct{}{}
	greetings["hru"] = struct{}{}
	greetings["How are you?"] = struct{}{}
}

func generateReply(history [][]byte) []byte {
	if _, ok := greetings[string(history[len(history)-1])]; ok {
		return []byte("Hi! How are you?")
	}

	return []byte("Hello")
}
