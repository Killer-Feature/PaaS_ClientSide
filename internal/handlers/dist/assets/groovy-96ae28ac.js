import{g as E}from"./index-b37b5403.js";function v(a,i){for(var o=0;o<i.length;o++){const n=i[o];if(typeof n!="string"&&!Array.isArray(n)){for(const t in n)if(t!=="default"&&!(t in a)){const e=Object.getOwnPropertyDescriptor(n,t);e&&Object.defineProperty(a,t,e.get?e:{enumerable:!0,get:()=>n[t]})}}}return Object.freeze(Object.defineProperty(a,Symbol.toStringTag,{value:"Module"}))}var s,f;function y(){if(f)return s;f=1;function a(e){return e?typeof e=="string"?e:e.source:null}function i(e){return o("(?=",e,")")}function o(...e){return e.map(c=>a(c)).join("")}function n(e,r={}){return r.variants=e,r}function t(e){const r="[A-Za-z0-9_$]+",c=n([e.C_LINE_COMMENT_MODE,e.C_BLOCK_COMMENT_MODE,e.COMMENT("/\\*\\*","\\*/",{relevance:0,contains:[{begin:/\w+@/,relevance:0},{className:"doctag",begin:"@[A-Za-z]+"}]})]),l={className:"regexp",begin:/~?\/[^\/\n]+\//,contains:[e.BACKSLASH_ESCAPE]},u=n([e.BINARY_NUMBER_MODE,e.C_NUMBER_MODE]),g=n([{begin:/"""/,end:/"""/},{begin:/'''/,end:/'''/},{begin:"\\$/",end:"/\\$",relevance:10},e.APOS_STRING_MODE,e.QUOTE_STRING_MODE],{className:"string"});return{name:"Groovy",keywords:{built_in:"this super",literal:"true false null",keyword:"byte short char int long boolean float double void def as in assert trait abstract static volatile transient public private protected synchronized final class interface enum if else for while switch case break default continue throw throws try catch finally implements extends new import package return instanceof"},contains:[e.SHEBANG({binary:"groovy",relevance:10}),c,g,l,u,{className:"class",beginKeywords:"class interface trait enum",end:/\{/,illegal:":",contains:[{beginKeywords:"extends implements"},e.UNDERSCORE_TITLE_MODE]},{className:"meta",begin:"@[A-Za-z]+",relevance:0},{className:"attr",begin:r+"[ 	]*:",relevance:0},{begin:/\?/,end:/:/,relevance:0,contains:[c,g,l,u,"self"]},{className:"symbol",begin:"^[ 	]*"+i(r+":"),excludeBegin:!0,end:r+":",relevance:0}],illegal:/#|<\//}}return s=t,s}var d=y();const b=E(d),m=v({__proto__:null,default:b},[d]);export{m as g};
