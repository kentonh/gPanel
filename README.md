# gPanel ![TravisCI gPanel Build](https://travis-ci.org/Ennovar/gPanel.svg?branch=master)

A web-hosting control panel written in Go.

## Stack

Backend: __[Go (1.8+)](https://golang.org/)__  
Database: __[Bolt](https://github.com/boltdb/bolt)__  
CSS Toolkit(s): __[Bootstrap 4](http://getbootstrap.com/) & [Font Awesome](http://fontawesome.io/)__  
JS Toolkit(s): __[jQuery](https://jquery.com/)__

## Contribution Set-up & Deployment

To get the repo...

```shell
# Go get the repo and append it to your $GOPATH
go get github.com/Ennovar/gPanel

# Navigate to the directory (replace $GOPATH with your actual $GOPATH)
cd $GOPATH/github.com/Ennovar/gPanel
```

To set your repo up to contribute...

```shell
# Fork the repo and add it to the list of remotes (replace your-username with your github username)
git remote add fork https://github.com/your-username/gPanel.git

# OPTIONAL: Change the names of the remotes
git remote rename origin upstream
git remote rename fork origin
```

To deploy...

```shell
# Build the binary
go build gpanel.go

# Execute binary as root (root access is needed for functions within the system package)
sudo ./gpanel
```
