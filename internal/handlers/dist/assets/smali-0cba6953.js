import{g as c}from"./index-b041cc2f.js";function g(r,a){for(var n=0;n<a.length;n++){const e=a[n];if(typeof e!="string"&&!Array.isArray(e)){for(const t in e)if(t!=="default"&&!(t in r)){const s=Object.getOwnPropertyDescriptor(e,t);s&&Object.defineProperty(r,t,s.get?s:{enumerable:!0,get:()=>e[t]})}}}return Object.freeze(Object.defineProperty(r,Symbol.toStringTag,{value:"Module"}))}var i,o;function m(){if(o)return i;o=1;function r(a){const n=["add","and","cmp","cmpg","cmpl","const","div","double","float","goto","if","int","long","move","mul","neg","new","nop","not","or","rem","return","shl","shr","sput","sub","throw","ushr","xor"],e=["aget","aput","array","check","execute","fill","filled","goto/16","goto/32","iget","instance","invoke","iput","monitor","packed","sget","sparse"],t=["transient","constructor","abstract","final","synthetic","public","private","protected","static","bridge","system"];return{name:"Smali",contains:[{className:"string",begin:'"',end:'"',relevance:0},a.COMMENT("#","$",{relevance:0}),{className:"keyword",variants:[{begin:"\\s*\\.end\\s[a-zA-Z0-9]*"},{begin:"^[ ]*\\.[a-zA-Z]*",relevance:0},{begin:"\\s:[a-zA-Z_0-9]*",relevance:0},{begin:"\\s("+t.join("|")+")"}]},{className:"built_in",variants:[{begin:"\\s("+n.join("|")+")\\s"},{begin:"\\s("+n.join("|")+")((-|/)[a-zA-Z0-9]+)+\\s",relevance:10},{begin:"\\s("+e.join("|")+")((-|/)[a-zA-Z0-9]+)*\\s",relevance:10}]},{className:"class",begin:`L[^(;:
]*;`,relevance:0},{begin:"[vp][0-9]+"}]}}return i=r,i}var l=m();const u=c(l),b=g({__proto__:null,default:u},[l]);export{b as s};
