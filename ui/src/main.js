import ForceGraph3D from "3d-force-graph";

const Graph = new ForceGraph3D(document.getElementById('view'));
var xmlhttp = new XMLHttpRequest;
var parsedJson;

Graph.backgroundColor("#0d1e1f");

Graph.onNodeClick(node => {
  // Aim at node from outside it
  const distance = 600;
  const distRatio = 1 + distance/Math.hypot(node.x, node.y, node.z);

  const newPos = node.x || node.y || node.z
    ? { x: node.x * distRatio, y: node.y * distRatio, z: node.z * distRatio }
    : { x: 0, y: 0, z: distance }; // special case if node is in (0,0,0)

  Graph.cameraPosition(
    newPos, // new position
    node, // lookAt ({ x, y, z })
    1800  // ms transition duration
  );
});


var canvas = Graph.renderer().domElement
canvas.id = "scene"
function resizeWindow(){
  Graph.camera().aspect = canvas.clientWidth / canvas.clientHeight;
  Graph.camera().updateProjectionMatrix()
}

addEventListener("resize", resizeWindow)

xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
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
