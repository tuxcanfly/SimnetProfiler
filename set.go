package main

import "github.com/conformal/btcjson"
import "fmt"
import "strconv"

func Serialize(item btcjson.ListTransactionsResult) string {
	amt := strconv.FormatFloat(item.Amount)
	fee := strconv.FormatFloat(item.Fee)
	return item.Address+','+amt+','+fee
}

func SerializeTransactions(data []btcjson.ListTransactionsResult} []string {
	result := []string
	for _, i := data {
		result.Add(Serialize(i))
	}
	return result
}
	

type Set struct {
  data []btcjson.ListTransactionsResult
}

func (this *Set) Add(element btcjson.ListTransactionsResult) {
  this.data = append(this.data, element)
}





func (this *Set) Remove(element btcjson.ListTransactionsResult) (bool) {
  for _, elem:= range this.data {
    if elem== element {
      return true
    }
  }
  return false
}


func (this *Set) IsMember(element btcjson.ListTransactionsResult) (bool) {
  for _, elem:= range this.data {
    if elem== element {
      return true
    }
  }
  return false
}

func (this *Set) Length() (btcjson.ListTransactionsResult) {
  return len(this.data)
}
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


func MakeNewSet() (Set) {
  return &Set{make([]btcjson.ListTransactionsResult, 0, 100)}
}
