import{g as s}from"./index-b37b5403.js";function c(e,a){for(var l=0;l<a.length;l++){const r=a[l];if(typeof r!="string"&&!Array.isArray(r)){for(const t in r)if(t!=="default"&&!(t in e)){const n=Object.getOwnPropertyDescriptor(r,t);n&&Object.defineProperty(e,t,n.get?n:{enumerable:!0,get:()=>r[t]})}}}return Object.freeze(Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}))}var i,u;function p(){if(u)return i;u=1;function e(a){return{name:"Julia REPL",contains:[{className:"meta",begin:/^julia>/,relevance:10,starts:{end:/^(?![ ]{6})/,subLanguage:"julia"},aliases:["jldoctest"]}]}}return i=e,i}var o=p();const f=s(o),g=c({__proto__:null,default:f},[o]);export{g as j};
