import{g as u}from"./index-b37b5403.js";function g(e,n){for(var a=0;a<n.length;a++){const t=n[a];if(typeof t!="string"&&!Array.isArray(t)){for(const r in t)if(r!=="default"&&!(r in e)){const o=Object.getOwnPropertyDescriptor(t,r);o&&Object.defineProperty(e,r,o.get?o:{enumerable:!0,get:()=>t[r]})}}}return Object.freeze(Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}))}var s,i;function p(){if(i)return s;i=1;function e(n){return{name:"Test Anything Protocol",case_insensitive:!0,contains:[n.HASH_COMMENT_MODE,{className:"meta",variants:[{begin:"^TAP version (\\d+)$"},{begin:"^1\\.\\.(\\d+)$"}]},{begin:/---$/,end:"\\.\\.\\.$",subLanguage:"yaml",relevance:0},{className:"number",begin:" (\\d+) "},{className:"symbol",variants:[{begin:"^ok"},{begin:"^not ok"}]}]}}return s=e,s}var c=p();const f=u(c),b=g({__proto__:null,default:f},[c]);export{b as t};
