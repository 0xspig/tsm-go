import * as THREE from 'three';
import WebGL from 'three/addons/capabilities/WebGL.js';

const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(75,window.innerWidth/window.innerHeight, .1, 1000);
const orthCam = new THREE.OrthographicCamera(-1,1,1,-1,0,1);
const renderer = new THREE.WebGLRenderer();

renderer.setSize( window.innerWidth, window.innerHeight )
document.body.appendChild( renderer.domElement)

const geometry = new THREE.BoxGeometry(1,1,1);
const material = new THREE.MeshBasicMaterial({color: 0x00ff00});
const cube = new THREE.Mesh(geometry, material);


let squareUniforms = {
    time: {value: 1.0 }
};

const squareGeometry = new THREE.PlaneGeometry(2,2);
const squareMaterial = new THREE.ShaderMaterial( {
	uniforms: squareUniforms,
	vertexShader: document.getElementById( 'vertexShader' ).textContent,
	fragmentShader: document.getElementById( 'fragmentShader' ).textContent

} );

var squareMesh = new THREE.Mesh(squareGeometry, squareMaterial)


scene.add(squareMesh);



function animate(){
    squareUniforms[ 'time' ].value = performance.now()/1000;
    renderer.render(scene, orthCam);
}

if ( WebGL.isWebGL2Available() ) {

	// Initiate function or other initializations here
    renderer.setAnimationLoop(animate)

} else {

	const warning = WebGL.getWebGL2ErrorMessage();
	document.getElementById( 'container' ).appendChild( warning );

}