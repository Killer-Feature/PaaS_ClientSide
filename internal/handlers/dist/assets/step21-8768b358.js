import{g as _}from"./index-b041cc2f.js";function c(n,e){for(var a=0;a<e.length;a++){const t=e[a];if(typeof t!="string"&&!Array.isArray(t)){for(const r in t)if(r!=="default"&&!(r in n)){const s=Object.getOwnPropertyDescriptor(t,r);s&&Object.defineProperty(n,r,s.get?s:{enumerable:!0,get:()=>t[r]})}}}return Object.freeze(Object.defineProperty(n,Symbol.toStringTag,{value:"Module"}))}var E,i;function S(){if(i)return E;i=1;function n(e){return{name:"STEP Part 21",aliases:["p21","step","stp"],case_insensitive:!0,keywords:{$pattern:"[A-Z_][A-Z0-9_.]*",keyword:"HEADER ENDSEC DATA"},contains:[{className:"meta",begin:"ISO-10303-21;",relevance:10},{className:"meta",begin:"END-ISO-10303-21;",relevance:10},e.C_LINE_COMMENT_MODE,e.C_BLOCK_COMMENT_MODE,e.COMMENT("/\\*\\*!","\\*/"),e.C_NUMBER_MODE,e.inherit(e.APOS_STRING_MODE,{illegal:null}),e.inherit(e.QUOTE_STRING_MODE,{illegal:null}),{className:"string",begin:"'",end:"'"},{className:"symbol",variants:[{begin:"#",end:"\\d+",illegal:"\\W"}]}]}}return E=n,E}var o=S();const l=_(o),p=c({__proto__:null,default:l},[o]);export{p as s};
