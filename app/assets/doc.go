// Package assets provide a quick implementation
// for in-memory static assets for this demo.
//
// This is obviously replaced by a container FS to
// support frontends with larger assets (e.g. react or angular).
//
// In order to regenerate the assets, one needs to install this utility:
// go get github.com/jessevdk/go-assets-builder
package assets

//go:generate go-assets-builder --strip-prefix /assets/html --package assets --variable Static --output=assets.go ../../assets/html/index.html ../../assets/html/upload.html
