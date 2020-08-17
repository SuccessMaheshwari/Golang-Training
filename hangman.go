package main

import (
  "fmt"
  "strings"
  "flag"
  "net/http"
  "encoding/json"
  "io/ioutil"
)

var dev = flag.Bool("dev", false, "Development Mode")

type hangman struct {
  word string
  entries map[string]bool
  correct_entries map[string]bool
  placeholder []string
  chances int8
  input string
  score int
  streak int
  lives int8
  round int
}
//Initializing hangman
func (h *hangman) initialize() {
  if *dev {               //if development mode
    h.word = "elephant"   // word is set as elephant
  } else {
    h.word = getWord()    //else any random word is set
  }
  h.chances = int8(len(h.word)-1)
  h.entries = make(map[string]bool)
  h.correct_entries = make(map[string]bool)
  h.placeholder = make([]string, 0, 0)
  for i:=0; i<len(h.word); i++ {
    h.placeholder = append(h.placeholder, "-")
  }
}
//function to display info 
func (h *hangman) display(){
  fmt.Println("\nLives:", h.lives, "Chances left: ", h.chances, "Score: ",h.score)
  fmt.Println("Word: ", h.placeholder)
  fmt.Println("Guessed letters: ", h.get_entries())
}
//function to get user input
func (h *hangman) getInput(){
  fmt.Print("Enter a letter/word (esc to exit): ")
  fmt.Scanln(&h.input)
  h.input = strings.ToLower(h.input)
}
//function to check whether word/letter is correct or not
func (h *hangman) contains(){
  if len(h.input) == len(h.word) && h.input == h.word {   //if correct word entered
    h.placeholder = strings.Split(h.word,"")              //update placeholder with word letters
  } else {
    if h.entries[h.input] == true || h.correct_entries[h.input] == true { //check if already entered
      fmt.Println("Already guessed...\n")
    } else {
      flag := 0
      for i, v := range h.word {
        if h.input == string(v){            //if letter present in word 
          h.correct_entries[h.input] = true //add in map
          h.placeholder[i] = h.input        //update placeholder
          flag += 1
        }
      }
      if flag == 0 {                        //if letter is not in word
        h.entries[h.input] = true           //enter in map
        h.chances -= 1                      //reduce chances
      }
    }
  }
}
//function to get wrong guesses
func (h *hangman) get_entries() (guesses []string){
  for i, _ := range h.entries {
    guesses = append(guesses, i)
  }
  return 
}
//function to get a random word from internet
func getWord() string {
  result, err := http.Get("https://random-word-api.herokuapp.com/word?number=10") //send GET request
  if err != nil {
    return "anonymous" //if error occured set word to "anonymous"
  }
  defer result.Body.Close()
  
  body, err := ioutil.ReadAll(result.Body) //read all content from response body
  
  var words []string
  err = json.Unmarshal(body,&words) //retrieve []words from json object

  for _,v := range words {
    if len(v) > 5 && len(v) <= 9 { 
      return v          //return word having length 5-9
    } 
  }
  return "anonymous"    //if no word found return "anonymous"
}
//check whether player won 
func (h *hangman) is_won() bool{
  if h.word == strings.Join(h.placeholder,"") {   //compare word with placeholder
    h.streak += 1                       //increment streak if letter is in word
    h.score += h.streak*100              //add points as per streak
    fmt.Println("Correct...")
    fmt.Println("The word was", h.word)
    return true
  }
  return false
}
//check whether player lost
func (h *hangman) is_lost() bool{
  if h.chances == 0 {             //check chances becomes 0 or not
    h.streak = 0                  //restore streak to 0
    h.score -= 50                 //cut score
    h.lives -= 1                  //reduce life      
    fmt.Println("Oops!!!")
    fmt.Println("The word was", h.word)
    return true
  }  
  return false
}

func main() {
  var h hangman
  exit := false
  flag.Parse()
  h.lives = 3
  h.round = 1
  for h.lives > 0 {
    h.initialize()
    fmt.Println("\n----------- Round ",h.round,"-----------")  
    for{
      h.display()                    //display info
      h.getInput()                   //take user input
      if byte(h.input[0]) == 27 {    //if esc is entered then exit
        exit = true
        break
      } else {
        h.contains()        
      }
      if h.is_won() {
        break 
      }    
      if h.is_lost() {    
        break
      }
    }
    if exit { 
      break
    }
    h.round += 1    
  }
  fmt.Println("Your final score is", h.score)
}