import{g as c}from"./index-b37b5403.js";function u(r,e){for(var a=0;a<e.length;a++){const n=e[a];if(typeof n!="string"&&!Array.isArray(n)){for(const t in n)if(t!=="default"&&!(t in r)){const i=Object.getOwnPropertyDescriptor(n,t);i&&Object.defineProperty(r,t,i.get?i:{enumerable:!0,get:()=>n[t]})}}}return Object.freeze(Object.defineProperty(r,Symbol.toStringTag,{value:"Module"}))}var l,s;function f(){if(s)return l;s=1;function r(e){return{name:"SML (Standard ML)",aliases:["ml"],keywords:{$pattern:"[a-z_]\\w*!?",keyword:"abstype and andalso as case datatype do else end eqtype exception fn fun functor handle if in include infix infixr let local nonfix of op open orelse raise rec sharing sig signature struct structure then type val with withtype where while",built_in:"array bool char exn int list option order real ref string substring vector unit word",literal:"true false NONE SOME LESS EQUAL GREATER nil"},illegal:/\/\/|>>/,contains:[{className:"literal",begin:/\[(\|\|)?\]|\(\)/,relevance:0},e.COMMENT("\\(\\*","\\*\\)",{contains:["self"]}),{className:"symbol",begin:"'[A-Za-z_](?!')[\\w']*"},{className:"type",begin:"`[A-Z][\\w']*"},{className:"type",begin:"\\b[A-Z][\\w']*",relevance:0},{begin:"[a-z_]\\w*'[\\w']*"},e.inherit(e.APOS_STRING_MODE,{className:"string",relevance:0}),e.inherit(e.QUOTE_STRING_MODE,{illegal:null}),{className:"number",begin:"\\b(0[xX][a-fA-F0-9_]+[Lln]?|0[oO][0-7_]+[Lln]?|0[bB][01_]+[Lln]?|[0-9][0-9_]*([Lln]|(\\.[0-9_]*)?([eE][-+]?[0-9_]+)?)?)",relevance:0},{begin:/[-=]>/}]}}return l=r,l}var o=f();const g=c(o),p=u({__proto__:null,default:g},[o]);export{p as s};
