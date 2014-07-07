// Copyright (c) 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/conformal/btcjson"
	"io/ioutil"
	"log"
	"os"
	//"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	//"github.com/davecgh/go-spew/spew"
)

func Serialize(item btcjson.ListTransactionsResult) string {
	amt := fmt.Sprintf("%.8f", item.Amount)
	fee := fmt.Sprintf("%.8f", item.Fee)
	return item.Address + "," + amt + "," + fee
}

func main() {
	//var cmd *exec.Cmd
	connected := make(chan struct{})
	var firstConn bool
	// based off of btcwebsocket example for reference
	ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnBtcdConnected: func(conn bool) {
			if conn && !firstConn {
				firstConn = true
				connected <- struct{}{}
			}
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
	log.Println("Waiting for BTCD")	
    
	// Unlock wallet.
	log.Println("Unlock Wallet")
	if err := client.WalletPassphrase("walletpass", 3600); err != nil {
		log.Printf("%s: Cannot unlock wallet: %v", "localhost:18556", err)

	}

	// count transactions
	client.NotifyNewTransactions(true)
	log.Println("Monitoring TXns")
	for {

		StartTime := time.Now().Unix()
		data, err := client.ListTransactionsCount("", 100)

		// Client
		if err != nil {
			log.Printf("ListTransactionsCount RPC Error: %s", err)
			break
		} else {
			txn := data[0]
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
					log.Printf("Transactions per second: %f", tps)
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
