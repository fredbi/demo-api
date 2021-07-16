package assets

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Static6a6ccdbaecac83ad55e9e6c9566e6188ea0600d6 = "<!DOCTYPE html>\n<html>\n<head>\n    <meta charset=\"utf-8\">\n</head>\n<body>\n<div>\n    <h1>Image uploader demo</h1>\n\n    <li>\n        <ul><a href=\"/upload.html\">Upload new image</a></ul>\n        <ul><a href=\"/images\">List images</a></ul>\n    </li>\n</div>\n</body>\n</html>\n"
var _Static542bf49e159060dabddc0576bba7281534863a2b = "<!DOCTYPE html> \n<html>  \n<head> \n  <meta charset=\\\"utf-8\\\"> \n</head> \n<body> \n  <div> \n      <a href=\"/\">Back</a>\n  </div> \n  <div> \n    <h1>Image upload</h1> \n    <form action=\"/images\" method=\"post\" enctype=\"multipart/form-data\"> \n    <p> \n      <input type=\"file\" name=\"file\"> \n    </p> \n    <p> \n      <button type=\"submit\">Submit</button> \n    </p> \n    </form> \n  </div> \n</body> \n</html>\n"

// Static returns go-assets FileSystem
var Static = assets.NewFileSystem(map[string][]string{"/": []string{"index.html", "upload.html"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1626439931, 1626439931751552079),
		Data:     nil,
	}, "/index.html": &assets.File{
		Path:     "/index.html",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1626441193, 1626441193153924563),
		Data:     []byte(_Static6a6ccdbaecac83ad55e9e6c9566e6188ea0600d6),
	}, "/upload.html": &assets.File{
		Path:     "/upload.html",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1626441245, 1626441245410518508),
		Data:     []byte(_Static542bf49e159060dabddc0576bba7281534863a2b),
	}}, "")
