import ForceGraph3D from "3d-force-graph";

const Graph = new ForceGraph3D(document.getElementById('3d-graph'));
var xmlhttp = new XMLHttpRequest;
var parsedJson;

Graph.backgroundColor("#0d1e1f");


xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
		console.log(this.responseText);
		Graph.graphData(JSON.parse(this.responseText));
    Graph.nodeColor(node => {
      console.log(node);
      //md files
      if (node.data_type == 1){
        return "#fff380";
      }
      // tags
      if (node.data_type == 2){
        return "#6b93c6";
      }
    });
    Graph.nodeVal(node => {
      //md files
      if (node.data_type == 1){
        return 1;
      }
      // tags
      if (node.data_type == 2){
        return 3;
      }
    })
    }
};
xmlhttp.open("GET", "/graph-json", true);
xmlhttp.send();
