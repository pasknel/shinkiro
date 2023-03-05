# Shinkiro

A tool to search and download files in Sharepoint.

## Build

Clone repository

```
git clone https://github.com/pasknel/shinkiro.git
```

Get required modules

```
go get ./...
```

Build binary

```
go build
```

Check binary

```
➜  ./shinkiro -h
Usage of ./shinkiro:
  -p string
        Password
  -s string
        Server Address
  -t string
        Search Term
  -u string
        Username
```

## Usage

```
./shinkiro -s https://example.sharepoint.com -u USERNAME -p PASSWORD -t SEARCH_TERM 
```