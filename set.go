package main

import (
       "github.com/conformal/btcjson"
       "fmt"

       )

// converts struct to a hopefully unique comma deliminated string
func Serialize(item btcjson.ListTransactionsResult) string {
	amt := fmt.Sprintf("%.8f", item.Amount)
	fee := fmt.Sprintf("%.8f", item.Fee)
	return item.Address+","+amt+","+fee
}

// converts 100 transactions struct array to string slice
func SerializeTransactions(data []btcjson.ListTransactionsResult) []string {
	result := make([]string, 0, 100)

	for i := 0; i < 100; i++ {
		result[i] = Serialize(data[i])
	}
	return result
}
	

// basic Set data structure
// works with strings for reasons of comparison operations
type Set struct {
  data []string
}

// add element to set
func (this *Set) Add(element string) {
  this.data = append(this.data, element)
}

// remove element from set
func (this *Set) Remove(element string) (bool) {
  for _, elem:= range this.data {
    if elem == element {
      return true
    }
  }
  return false
}

// test if element is member of the set
func (this *Set) IsMember(element string) (bool) {
  for _, elem:= range this.data {
    if elem== element {
      return true
    }
  }
  return false
}

// report length of the set
func (this *Set) Length() (int) {
  return len(this.data)
}

// remove duplicates in the set
func (this *Set) Deduplicate() {
  length := len(this.data) - 1
  for i := 0; i < length; i++ {
    for j := i + 1; j <= length; j++ {
      if (this.data[i] == this.data[j]) {
        this.data[j] = this.data[length]
        this.data = this.data[0:length]
        length--
        j--
      }
    }
  }
}


// create new set of initial size zero and capacity 100
// adjust for proper network sampling conditions
func MakeNewSet() Set {
  return Set{make([]string, 0, 100)}
}

func main(){
	fmt.Println("Testing 123")
}
	
