package suggest

import (
	"strconv"
	"testing"
)

func TestBasic(t *testing.T) {

	suggest := NewSuggest()
	suggest.AddDocument("bitbucket repository")
	suggest.AddDocument("git repo")
	suggest.AddDocument("software development")

	if result := suggest.Search("repo"); len(result) != 2 {
		t.Error("FAIL: expected 2, actual: " + strconv.Itoa(len(result)))
	}

	if result := suggest.Search("repo"); len(result) != 2 {
		t.Error("FAIL: expected 2, actual: " + strconv.Itoa(len(result)))
	}

	if result := suggest.Search("dev"); len(result) != 1 {
		t.Error("FAIL: suggest.Search(\"dev\"). actual result: " + strconv.Itoa(len(result)))
	}

	if result := suggest.Search("bi"); len(result) != 1 {
		t.Error("FAIL")
	}
}

func TestRankedResult(t *testing.T) {
	suggest := NewSuggest()
	suggest.AddDocument("bitbucket repository")
	suggest.AddDocument("repo")
	suggest.AddDocument("git repo")

	result := suggest.Search("repo")

	if len(result) != 3 {
		t.Error("FAIL")
	}

	if result[0] != "repo" {
		t.Error("FAIL")
	}

	if result[1] != "git repo" {
		t.Error("FAIL")
	}

	if result[2] != "bitbucket repository" {
		t.Error("FAIL")
	}

}

func TesttokenizeSymbol(t *testing.T) {
	symbol := "__consumer_offsets"
	sanitized := tokenizeSymbol(symbol)
	if sanitized != "consumer offsets" {
		t.Error("FAIL santizeSymbols. Expected 'consumer offsets', Actual: " + sanitized)
	}

	symbol = "__samza_checkpoint_samza-hello-world-i001"
	sanitized = tokenizeSymbol(symbol)
	if sanitized != "samza checkpoint samza hello world i001" {
		t.Error("Failed")
	}

	symbol = "OneWord"
	sanitized = tokenizeSymbol(symbol)
	if sanitized != "one word" {
		t.Error("FAILED")
	}
}
