SimnetProfiler
==============

Network Data and Metrics from the Conformal Simnet

Status:

  * Finishing Set Struct and Methods for Unique Transactions
  
Theory of Operation:

  * Process will listen into the Simnetwork via btcd at arbitrary times (network "mile markers")
  count transactions over a window of time and report results.
  Number of transactions, fee size and blockchain bloat under consideration for measurement 
  and reporting.
  


Metrics [TODO]

	1. transactions per second
	2. max fee size per second
	3. block chain size
	3. create sse server for dashboard
	4. add simnet control (knobs)
	    a. create small world network sim [near completion]
	    b. create scale free network  sim [tentative]
	    c. temporal network tracking and comparisons
 
       
