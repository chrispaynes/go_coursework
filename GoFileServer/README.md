#GoFileServer
Command Line File Server

##Todo
* Hot Module Reloading
* Markdown Parser for README.md

##Features
* Accepting Command Line Arguments
* Flag Package
* HTTP Server

##Usage
Compile and run code
`$ go run main.go -port=1337 -path="../"`
1. Click to accept incoming network connects on the port.
2. Open a browser and visit http://localhost:1337

Output help for flags
```
$ go build
$ ./GoFileServer -help
  
  -path string
        filepath to serve
  -port string
        communication port used to serve HTTP on (default: "3000")
```

Provide alternative name for build file
`$ go build -o fs`

Install globally executable binary to Go/Bin 
`$ go install`

Create Linux Binary
`$ GOOS=linux GOARCH=arm go build`

Create Windows Binary
`$ GOOS=windows GOARCH=386 go build`

##Running Tests
