# gPanel

A web-hosting control panel written in Go.

## Stack

Version Control: __[git](https://git-scm.com/)__  
Backend: __[Go](https://golang.org/)__  
Database: __[Bolt](https://github.com/boltdb/bolt)__  
Frontend Framework: __[Angular 4](https://angular.io/)__  
CSS Preprocessor: __[SASS](http://sass-lang.com/)__  
CSS Toolkit: __[Bootstrap 4](http://getbootstrap.com/)__  
Package Manager: __[npm](https://www.npmjs.com/)__

## Contribution Set-up & Deployment

```shell
# Go get the repo and append it to your $GOPATH
go get github.com/Ennovar/gPanel

# Navigate to the directory (replace $GOPATH with your actual $GOPATH)
cd $GOPATH/github.com/Ennovar/gPanel
```

To set your repo up to contribute...

```shell
# Fork the repo and add it to the list of remotes
git remote add fork https://github.com/Ennovar/gPanel.git

# OPTIONAL: Change the names of the remotes
git remote rename origin upstream
git remote rename fork origin
```

To deploy...

```shell
# Starting gPanel
go run main.go

# OPTIONAL: Create binary to run
go build main.go
./main
```
