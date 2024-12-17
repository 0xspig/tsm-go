import ForceGraph3D from "3d-force-graph";

const Graph = new ForceGraph3D(document.getElementById('3d-graph'));
var xmlhttp = new XMLHttpRequest;
var parsedJson;


xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
		console.log(this.responseText);
		Graph.graphData(JSON.parse(this.responseText));
    }
};
xmlhttp.open("GET", "/graph-json", true);
xmlhttp.send();
