// simnet network profiler

package main

import (
	"fmt"
	"github.com/conformal/btcjson"
	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func Serialize(item btcjson.ListTransactionsResult) string {
	amt := fmt.Sprintf("%.8f", item.Amount)
	fee := fmt.Sprintf("%.8f", item.Fee)
	return item.Address + "," + amt + "," + fee
}

func main() {

	// based off of btcwebsocket example for reference
	ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance btcutil.Amount, confirmed bool) {
			log.Printf("New balance for account %s: %v", account,
				balance)
		},
	}

	// Connect to local btcwallet RPC server using websockets.
	certHomeDir := btcutil.AppDataDir("btcwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:18554",
		Endpoint:     "ws",
		User:         "rpcuser",
		Pass:         "rpcpass",
		Certificates: certs,
	}
	client, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	client.NotifyNewTransactions(true)

	for {
		StartTime := time.Now().Unix()
		data, err := client.ListTransactionsCount("", 100)
		txn := data[0]
		// Client

		if err != nil {
			log.Printf("ListTransactionsCount RPC Error: %s", err)
			break
		} else {

			for {
				TxnArray, err := client.ListTransactionsCount("", 100)
				if err != nil {
					log.Printf("ListTransactionsCount RPC Error: %s", err)
					break
				}

				log.Println("Checking Transactions")
				var TxnFound bool = true

				for i := range TxnArray {
					if Serialize(TxnArray[i]) == Serialize(txn) {
						TxnFound = true
						break
					} else {

						TxnFound = false

					}

				}

				if TxnFound == false {
					EndTime := time.Now().Unix()
					delta := float64((StartTime - EndTime))

					if delta == 0 {
						delta = float64(0.00000001) // at least a satoshi
					}
					tps := 100.0 / delta
					log.Printf("Transactions per second: %d", tps)
					break
				}
			}
		}
	}

	// shutdown with ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Printf("Client shutdown: %v, stopping proc and exiting..", sig)

			time.AfterFunc(time.Second*1, func() {
				log.Println("Client shutting down...")
				client.Shutdown()
				log.Println("Client shutdown complete.")
			})
		}
		os.Exit(1)
	}()

	client.WaitForShutdown()
}
