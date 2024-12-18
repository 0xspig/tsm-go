import ForceGraph3D from "3d-force-graph";

const Graph = new ForceGraph3D(document.getElementById('view'))
  .width(document.getElementById("view").getBoundingClientRect().width - 1)
  .height(document.getElementById("view").getBoundingClientRect().height-100);
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

  // ui stuff
  document.getElementById("data-content").innerHTML ='<h1 class="nodeTitle">Node id: '+node.name+'</h1><br>'
});


var canvas = Graph.renderer().domElement
canvas.id = "scene"
function resizeWindow(){
  Graph.width(document.getElementById("view").getBoundingClientRect().width - 1)
  Graph.height(document.getElementById("view").getBoundingClientRect().height);
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


document.getElementById("data-tab").addEventListener("click", dataToggle)
var dataViewing = true;
function dataToggle(){
  if (dataViewing){
    document.getElementById("data-content").style.flexBasis="0%";
    document.getElementById("data-content").style.display ="none";
    var elemChildren = document.getElementById("data-content").children
    for (var i = 0; i < elemChildren.length; i++){
      elemChildren[i].style.display = "none";
    }
    dataViewing=false
  }else{
    document.getElementById("data-content").style.flexBasis="45%";
    document.getElementById("data-content").style.display ="inline";
    var elemChildren = document.getElementById("data-content").children
    for (var i = 0; i < elemChildren.length; i++){
      elemChildren[i].style.display = "inline";
    }
    dataViewing=true
  }
    resizeWindow()
}