import{g as c}from"./index-b37b5403.js";function s(r,e){for(var o=0;o<e.length;o++){const t=e[o];if(typeof t!="string"&&!Array.isArray(t)){for(const n in t)if(n!=="default"&&!(n in r)){const a=Object.getOwnPropertyDescriptor(t,n);a&&Object.defineProperty(r,n,a.get?a:{enumerable:!0,get:()=>t[n]})}}}return Object.freeze(Object.defineProperty(r,Symbol.toStringTag,{value:"Module"}))}var f,i;function b(){if(i)return f;i=1;function r(e){return{name:"Backus–Naur Form",contains:[{className:"attribute",begin:/</,end:/>/},{begin:/::=/,end:/$/,contains:[{begin:/</,end:/>/},e.C_LINE_COMMENT_MODE,e.C_BLOCK_COMMENT_MODE,e.APOS_STRING_MODE,e.QUOTE_STRING_MODE]}]}}return f=r,f}var u=b();const _=c(u),g=s({__proto__:null,default:_},[u]);export{g as b};
