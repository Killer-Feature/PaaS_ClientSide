import{g as b}from"./index-b37b5403.js";function c(e,n){for(var o=0;o<n.length;o++){const r=n[o];if(typeof r!="string"&&!Array.isArray(r)){for(const t in r)if(t!=="default"&&!(t in e)){const u=Object.getOwnPropertyDescriptor(r,t);u&&Object.defineProperty(e,t,u.get?u:{enumerable:!0,get:()=>r[t]})}}}return Object.freeze(Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}))}var a,i;function f(){if(i)return a;i=1;function e(n){return{name:"ERB",subLanguage:"xml",contains:[n.COMMENT("<%#","%>"),{begin:"<%[%=-]?",end:"[%-]?%>",subLanguage:"ruby",excludeBegin:!0,excludeEnd:!0}]}}return a=e,a}var s=f();const g=b(s),d=c({__proto__:null,default:g},[s]);export{d as e};
