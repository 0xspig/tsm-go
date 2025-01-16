// rollup.config.mjs
import commonjs from "@rollup/plugin-commonjs";
import nodeResolve from "@rollup/plugin-node-resolve";

export default {
	input: 'ui/src/main.js',
	output: {
		file: 'ui/static/gen/main.js',
		format: 'es'
	},
	plugins:[
		commonjs(),
		nodeResolve(),
	]
};