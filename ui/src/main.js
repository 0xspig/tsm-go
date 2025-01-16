import ForceGraph3D from "3d-force-graph";

const Graph = new ForceGraph3D(document.getElementById('view'))
  .width(document.getElementById("view").getBoundingClientRect().width - 1)
  .height(document.getElementById("view").getBoundingClientRect().height-100);
var parsedJson;

Graph.backgroundColor("#0d1e1f");

Graph.onNodeClick(node => {
  targetNode(node.id);
});


var canvas = Graph.renderer().domElement
canvas.id = "scene"
  Graph.width(document.getElementById("view").getBoundingClientRect().width - 1)
  Graph.height(document.getElementById("view").getBoundingClientRect().height);
  Graph.camera().aspect = canvas.clientWidth / canvas.clientHeight;
  Graph.camera().updateProjectionMatrix()
function resizeWindow(){
  Graph.width(document.getElementById("view").getBoundingClientRect().width - 1)
  Graph.height(document.getElementById("view").getBoundingClientRect().height);
  Graph.camera().aspect = canvas.clientWidth / canvas.clientHeight;
  Graph.camera().updateProjectionMatrix()
}

addEventListener("resize", resizeWindow)

var xmlhttp = new XMLHttpRequest;
xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
		Graph.graphData(JSON.parse(this.responseText));
    console.log("graph data updated");
    pushGraphParams();
    }
    // target home node
    targetNode("home.md");
};

function pushGraphParams(){

    Graph.nodeColor(node => {
      if (node.targeted == true){
        return "#ff33ff";
      }
      //md files
      if (node.data_type == 1){
        return "#fff380";
      }
      // tags
      if (node.data_type == 2){
        return "#6b93c6";
      }
      // categories
      if (node.data_type == 3){
        return "#63335c";
      }
      // external
      if (node.data_type == 4){
        return "#68d588";
      }
    });

    Graph.nodeVal(node => {
      if (node.targeted == true){
        return 8;
      }
      //md files
      if (node.data_type == 1){
        return 2;
      }
      // tags
      if (node.data_type == 2){
        return 1;
      }
      // categories
      if (node.data_type == 3){
        return 4;
      }
      // external
      if (node.data_type == 4){
        return 1;
      }
    })
}

export function targetNode(nodeID){

    var targetNode;

    console.log("targeting Node: "+nodeID);
    Graph.graphData().nodes.forEach(node => {
      if (node.id == nodeID){
        targetNode = node;
        // if external
        if (targetNode.data_type == 4) {
          window.open(node.source, '_blank').focus();
          return;
        }
        targetNode.targeted = true;
      }else{
        node.targeted = false;
      }
    });
    console.log(targetNode);
    if (targetNode.data_type == 4){
      return;
    } 
    // Aim at targetNode from outside it
    const distance = 600;
    const distRatio = 1 + distance/Math.hypot(targetNode.x, targetNode.y, targetNode.z);

    const newPos = targetNode.x || targetNode.y || targetNode.z
      ? { x: targetNode.x * distRatio, y: targetNode.y * distRatio, z: targetNode.z * distRatio }
      : { x: 0, y: 0, z: distance }; // special case if targetNode is in (0,0,0)

    Graph.cameraPosition(
      newPos, // new position
      targetNode, // lookAt ({ x, y, z })
      1800  // ms transition duration
    );

    // ui stuff
    getNodeData(targetNode);
    pushGraphParams();
}

xmlhttp.open("GET", "/graph-json", true);
xmlhttp.send();

/* !!!!! TAB FUNCTION !!!
  I'm removing this for the time being since its ugly and doesn't add anything to the UX

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

*/

function getNodeData(node){
  var node_data_request = new XMLHttpRequest;
  node_data_request.onreadystatechange = function() {
    document.getElementById("data").innerHTML = node_data_request.responseText;
  }
  node_data_request.open("GET", "/node-data/"+node.id, true);
  node_data_request.send();
}