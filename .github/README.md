# gPanel ![TravisCI gPanel Build](https://travis-ci.org/Ennovar/gPanel.svg?branch=master)

A web-hosting control panel written in Go.

*__Note:__ This software currently only runs on Linux systems.*

This project may qualify for Mozilla's 2018 Global Sprint, in efforts to satisfy their rules I am linking their Community Participation Guidelines [here](https://www.mozilla.org/en-US/about/governance/policies/participation/) along with their already being a code of conduct within this project located [here](https://github.com/Ennovar/gPanel/blob/master/.github/CODE_OF_CONDUCT.md).

## Table of Contents
1. [Technology Stack](#technology-stack)
2. [Preview Images](#preview-images)
3. [Contribution Set-up & Deployment](#contribution-set-up--deployment)
4. [Installation for Use](#installation-for-use)
    * [System Requirements](#system-requirements)
    * [Installing openssh-server](#installing-openssh-server)
    * [Creating the Host Key-pair](#creating-the-host-key-pair)
    * [.ssh Folder and Files Permissions Reference](#ssh-folder-and-files-permissions-reference)
    * [Getting the Repository and Running](#getting-the-repository-and-running)

## Technology Stack

Backend: __[Go (1.8+)](https://golang.org/)__  
Database: __[Bolt](https://github.com/boltdb/bolt)__  
CSS Toolkit(s): __[Bootstrap 4](http://getbootstrap.com/) & [Font Awesome](http://fontawesome.io/)__  
JS Toolkit(s): __[jQuery](https://jquery.com/)__

## Preview Images

gPanel Structure
![Image of gPanel Structure](https://nextwavesolutions.io/images/gPanelStructure.png)

gPanel Server
![Image of gPanel Server](https://user-images.githubusercontent.com/30050545/36277136-9d0cdffc-1255-11e8-8a33-b503087a32f8.png)

gPanel Account
![Image of gPanel Account](https://user-images.githubusercontent.com/30050545/36277135-9cf4feaa-1255-11e8-8957-9f02a9cfb7e4.png)

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

## Installation for Use

#### System Requirements

- Linux
	- adduser command (already installed on most debian-based Linux systems)
	- deluser command (already installed on most debian-based linux systems)
	- ssh-keygen command (already installed on most debian-based linux systems)
	- openssh-server installed (installation guide below)
	- golang (installation guide below)
	- [php-cgi](http://us3.php.net/downloads.php) IF you want to be able to serve .php files
- OSX
	- Currently there is no support for OSX, but it is planned for the future.
-----------------------------------


#### Installing Golang

1. sudo apt-get purge golang*
2. Download latest version from https://www.golang.org/dl/
3. sudo tar -C /usr/local -xzf go[VERSION].[OS]-[ARCH].tar.gz
4. For system-wide installation (reccommended)
	a. vim /etc/profile
	b. Add "export PATH=$PATH:/usr/local/go/bin"
5. For local installation
	a. If ~/.profile doesn't exist then create it (touch ~/.profile)
	b. Add "export PATH=$PATH:/usr/local/go/bin" to said file
6. Logout and login for changes to /etc/profile or ~/.profile to take effect
7. mkdir ~/go && mkdir ~/go/bin && mkdir ~/go/src && mkdir ~/go/pkg
8. GOROOT=~/go


#### Installing openssh-server

1. sudo apt-get install openssh-server
2. sudo cp /etc/ssh/sshd_config /etc/ssh/sshd_config.default
3. sudo vim /etc/ssh/sshd_config
4. Uncomment (no # before line) and/or set the following lines in /etc/ssh/sshd_config
	- PermitRootLogin no
	- AuthorizedKeysFile %h/.ssh/authorized_keys
	- PasswordAuthentication no
	- PermitEmptyPasswords no
	- RSAAuthentication yes
	- PubkeyAuthentication yes
5. sudo systemctl restart ssh


#### Creating the Host Key-pair

1. ssh-keygen -t rsa -N "PASSWORD" -f ~/.ssh/id_rsa [change PASSWORD]


#### .ssh Folder and Files Permissions Reference

1. To check permissions
	a. cd ~/.ssh
	b. ls -l -a (look at ls -l -a dump below)
2. To change permissions
	a. chmod [PERMISSIONS NUMBER] [FILE]
3. To change ownership
	a. chown [USER]: [FILE]

.ssh [700 && owned by correct user]  
&nbsp;&nbsp;&nbsp;&nbsp;id_rsa [600 && owned by correct user]  
&nbsp;&nbsp;&nbsp;&nbsp;id_rsa.pub [644 && owned by correct user]  
&nbsp;&nbsp;&nbsp;&nbsp;authorized_keys [644 && owned by correct user]  
&nbsp;&nbsp;&nbsp;&nbsp;known_hosts [644 && owned by the correct user]  

Correct output of ls -l -a of ~/.ssh
```shell
drwx------ 2 root root 4096 Jan 17 14:49 .  
drwx------ 7 root root 4096 Jan 17 14:42 ..  
-rw-r--r-- 1 root root    0 Jan 17 14:49 authorized_keys  
-rw------- 1 root root 1766 Jan 17 14:43 id_rsa  
-rw-r--r-- 1 root root  401 Jan 17 14:43 id_rsa.pub  
-rw-r--r-- 1 root root  444 Oct 10  2016 known_hosts
```


#### Getting the Repository and Running

1. go get github.com/Ennovar/gPanel
2. cd ~/go/src/github.com/Ennovar/gPanel
3. go build gpanel.go
4. sudo ./gpanel
