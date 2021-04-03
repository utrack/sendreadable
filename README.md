# sendreadable

A small service that converts web articles to PDFs. It is possible to download them to your PC or send them straight to reMarkable via reMarkable's API.

This is something I've jotted down over one weekend - the code is dirty and awful, there's not much to configure, but it works.

## Prerequisites
go >= v1.14
xelatex (Arch packages: `texlive-bin texlive-core texlive-fontsextra texlive-langchinese texlive-langcyrillic texlive-langextra texlive-latexextra texlive-pstricks texlive-qrcode`)
pkger (https://github.com/markbates/pkger)

## Installation/operation
```
make build
openssl genrsa -out priv.key 4096 # private key for JWT tokens
mkdir /tmp/sendreadable
./sendreadable -key priv.key # add -secure true to send out SSL-only cookies
```
Server will be bound to port 3333.
