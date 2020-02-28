# modal-logic-theorem-prover

### Authors
Andrea Jegher

### References
* Based on the book by P. Jackson, “Logic-based  knowledge  representation,” MIT Press, 1985

### Thanks
* Thanks to [KaTeX](https://github.com/KaTeX/KaTeX) for a JavaScript library for TeX math rendering on the web.

### Deploy
* [Install golang](https://golang.org/doc/install)
* ```cd $GOPATH/src```
* ```mkdir github.com```
* ```cd github.com```
* ```git clone https://github.com/AndreaJegher/gomoltp.git```
* Install the local command
* ```cd gomoltp/cmd/moltprunner```
* ```go install```
* Install the http server
* ```cd gomoltp/cmd/moltpserver```
* ```go install```

### Examples
* Local command
* ```$GPATH/bin/moltprunner -f '\Box \Box  p \to \Diamond \Diamond p'```
* Http Server
* ```./moltpserver -static $GPATH/src/github.com/gomoltp/cmd/moltpserver/static -templates $GPATH/src/github.com/gomoltp/cmd/moltpserver/templates -v```
* Then visit [http://localhost:4000](http://localhost:4000) from your browser
