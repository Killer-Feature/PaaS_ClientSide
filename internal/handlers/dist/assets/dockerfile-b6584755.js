import{g as l}from"./index-b37b5403.js";function f(r,e){for(var n=0;n<e.length;n++){const t=e[n];if(typeof t!="string"&&!Array.isArray(t)){for(const o in t)if(o!=="default"&&!(o in r)){const i=Object.getOwnPropertyDescriptor(t,o);i&&Object.defineProperty(r,o,i.get?i:{enumerable:!0,get:()=>t[o]})}}}return Object.freeze(Object.defineProperty(r,Symbol.toStringTag,{value:"Module"}))}var a,c;function u(){if(c)return a;c=1;function r(e){return{name:"Dockerfile",aliases:["docker"],case_insensitive:!0,keywords:"from maintainer expose env arg user onbuild stopsignal",contains:[e.HASH_COMMENT_MODE,e.APOS_STRING_MODE,e.QUOTE_STRING_MODE,e.NUMBER_MODE,{beginKeywords:"run cmd entrypoint volume add copy workdir label healthcheck shell",starts:{end:/[^\\]$/,subLanguage:"bash"}}],illegal:"</"}}return a=r,a}var s=u();const d=l(s),p=f({__proto__:null,default:d},[s]);export{p as d};
