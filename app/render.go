package app

import "html/template"

// A simple template for server-side rendering.
//
// This is used because I chose to go with a minimalist frontend, low on js.
// Therefore, this is not to be considered a maintenable solution for frontend.
//
// For a real-life API, this should be reverted to the more natural json message.
var list *template.Template

func init() {
	list = template.Must(template.New("list").Parse(listTemplate))
}

const listTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
	<style>
	 .col {
	   float: left;
	   width: 25%;
	 }
	 .row {
	   display: flex;
	 }
	 .row:after {
	   content: "";
	   display: table;
	   clear: both;
	 }
	</style>
	<script>
	function delImage(key) {
	  window.alert("deleting: " + key);

	  var xhr = new XMLHttpRequest();
	  xhr.open('DELETE', '/images/' + key, true);
	  xhr.send();
	}

	function patchImage(key) {
	  window.alert("updating: " + key);

	  var formElem = document.getElementById("patch-" + key);
	  var body = new FormData(formElem);
	  var xhr = new XMLHttpRequest();

	  xhr.open('PATCH', '/images/' + key, true);
	  xhr.send(body);
	}
	</script>
</head>
<body>
<div>
    <a href="/">Back</a>
    <h1>Uploaded images</h1>

    <li>
	  {{ range . }}
        <ul>
		<div class="row">
		  <div class="col">
			<img src="data:image/png;base64,{{ .Thumb }}" alt="thumbnail {{ .Key }}">
		  </div>
		  <div class="col">
			<a href="/images/{{ .Key }}" target="_blank">
		    {{ .Key }}
			</a>
		  </div>
		  <div class="col">
			  <form id="patch-{{ .Key }}" method="POST" action="javascript:void(0);" enctype="multipart/form-data" onsubmit="patchImage('{{ .Key }}');">
			    <input type="file" name="file">
			    <input type="submit" value="Submit">
			  </form>
		  </div>
		  <div class="col">
			  <button onclick="delImage('{{ .Key }}');">Delete</button>
		  </div>
		</div>
		</ul>
	  {{- end }}
    </li>
</div>
</body>
</html>
`
