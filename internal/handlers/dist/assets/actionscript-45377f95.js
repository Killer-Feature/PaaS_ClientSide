import{g as _}from"./index-b041cc2f.js";function d(r,a){for(var i=0;i<a.length;i++){const e=a[i];if(typeof e!="string"&&!Array.isArray(e)){for(const t in e)if(t!=="default"&&!(t in r)){const n=Object.getOwnPropertyDescriptor(e,t);n&&Object.defineProperty(r,t,n.get?n:{enumerable:!0,get:()=>e[t]})}}}return Object.freeze(Object.defineProperty(r,Symbol.toStringTag,{value:"Module"}))}var c,o;function l(){if(o)return c;o=1;function r(e){return e?typeof e=="string"?e:e.source:null}function a(...e){return e.map(n=>r(n)).join("")}function i(e){const t=/[a-zA-Z_$][a-zA-Z0-9_$]*/,n=/([*]|[a-zA-Z_$][a-zA-Z0-9_$]*)/,u={className:"rest_arg",begin:/[.]{3}/,end:t,relevance:10};return{name:"ActionScript",aliases:["as"],keywords:{keyword:"as break case catch class const continue default delete do dynamic each else extends final finally for function get if implements import in include instanceof interface internal is namespace native new override package private protected public return set static super switch this throw try typeof use var void while with",literal:"true false null undefined"},contains:[e.APOS_STRING_MODE,e.QUOTE_STRING_MODE,e.C_LINE_COMMENT_MODE,e.C_BLOCK_COMMENT_MODE,e.C_NUMBER_MODE,{className:"class",beginKeywords:"package",end:/\{/,contains:[e.TITLE_MODE]},{className:"class",beginKeywords:"class interface",end:/\{/,excludeEnd:!0,contains:[{beginKeywords:"extends implements"},e.TITLE_MODE]},{className:"meta",beginKeywords:"import include",end:/;/,keywords:{"meta-keyword":"import include"}},{className:"function",beginKeywords:"function",end:/[{;]/,excludeEnd:!0,illegal:/\S/,contains:[e.TITLE_MODE,{className:"params",begin:/\(/,end:/\)/,contains:[e.APOS_STRING_MODE,e.QUOTE_STRING_MODE,e.C_LINE_COMMENT_MODE,e.C_BLOCK_COMMENT_MODE,u]},{begin:a(/:\s*/,n)}]},e.METHOD_GUARD],illegal:/#/}}return c=i,c}var s=l();const E=_(s),p=d({__proto__:null,default:E},[s]);export{p as a};
