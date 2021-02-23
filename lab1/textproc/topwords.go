// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

package textproc

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func topWords(path string, K int) []WordCount {
	//read file
	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
	}

	//declare the struct slice to return and a map to store occurences for each word
	wc := make([]WordCount, 0)
	m := make(map[string]int)

	//split the string into a string array, splitting at whitespace
	sliceOfWords := strings.Fields(string(data))

	//add each key to map where it's value is occurrence
	for i, s := range sliceOfWords {
		fmt.Println(i, s)
		m[s] = m[s] + 1
	}

	//add each map entry to the struct slice
	for key, element := range m {
		fmt.Println(key, element)
		wc = append(wc, WordCount{key, element})
	}
	//sort words
	sortWordCounts(wc)

	//only keep top K
	topK := wc[:K]

	return topK
}

//--------------- DO NOT MODIFY----------------!

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

// Method to convert struct to string format
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.

func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
