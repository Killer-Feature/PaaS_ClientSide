import{g}from"./index-b37b5403.js";function p(t,o){for(var n=0;n<o.length;n++){const e=o[n];if(typeof e!="string"&&!Array.isArray(e)){for(const r in e)if(r!=="default"&&!(r in t)){const a=Object.getOwnPropertyDescriptor(e,r);a&&Object.defineProperty(t,r,a.get?a:{enumerable:!0,get:()=>e[r]})}}}return Object.freeze(Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}))}var s,c;function l(){if(c)return s;c=1;function t(o){return{name:"Tagger Script",contains:[{className:"comment",begin:/\$noop\(/,end:/\)/,contains:[{begin:/\(/,end:/\)/,contains:["self",{begin:/\\./}]}],relevance:10},{className:"keyword",begin:/\$(?!noop)[a-zA-Z][_a-zA-Z0-9]*/,end:/\(/,excludeEnd:!0},{className:"variable",begin:/%[_a-zA-Z0-9:]*/,end:"%"},{className:"symbol",begin:/\\./}]}}return s=t,s}var i=l();const u=g(i),E=p({__proto__:null,default:u},[i]);export{E as t};
