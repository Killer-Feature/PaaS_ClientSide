import{g as T}from"./index-b37b5403.js";function N(t,a){for(var e=0;e<a.length;e++){const n=a[e];if(typeof n!="string"&&!Array.isArray(n)){for(const r in n)if(r!=="default"&&!(r in t)){const i=Object.getOwnPropertyDescriptor(n,r);i&&Object.defineProperty(t,r,i.get?i:{enumerable:!0,get:()=>n[r]})}}}return Object.freeze(Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}))}var b,p;function L(){if(p)return b;p=1;function t(a){var e="[a-zA-Z_\\-+\\*\\/<=>&#][a-zA-Z0-9_\\-+*\\/<=>&#!]*",n="\\|[^]*?\\|",r="(-|\\+)?\\d+(\\.\\d+|\\/\\d+)?((d|e|f|l|s|D|E|F|L|S)(\\+|-)?\\d+)?",i={className:"literal",begin:"\\b(t{1}|nil)\\b"},s={className:"number",variants:[{begin:r,relevance:0},{begin:"#(b|B)[0-1]+(/[0-1]+)?"},{begin:"#(o|O)[0-7]+(/[0-7]+)?"},{begin:"#(x|X)[0-9a-fA-F]+(/[0-9a-fA-F]+)?"},{begin:"#(c|C)\\("+r+" +"+r,end:"\\)"}]},l=a.inherit(a.QUOTE_STRING_MODE,{illegal:null}),v=a.COMMENT(";","$",{relevance:0}),c={begin:"\\*",end:"\\*"},E={className:"symbol",begin:"[:&]"+e},o={begin:e,relevance:0},O={begin:n},m={begin:"\\(",end:"\\)",contains:["self",i,l,s,o]},u={contains:[s,l,c,E,m,o],variants:[{begin:"['`]\\(",end:"\\)"},{begin:"\\(quote ",end:"\\)",keywords:{name:"quote"}},{begin:"'"+n}]},f={variants:[{begin:"'"+e},{begin:"#'"+e+"(::"+e+")*"}]},g={begin:"\\(\\s*",end:"\\)"},d={endsWithParent:!0,relevance:0};return g.contains=[{className:"name",variants:[{begin:e,relevance:0},{begin:n}]},d],d.contains=[u,f,g,i,s,l,v,c,E,O,o],{name:"Lisp",illegal:/\S/,contains:[s,a.SHEBANG(),i,l,v,u,f,g,o]}}return b=t,b}var _=L();const M=T(_),A=N({__proto__:null,default:M},[_]);export{A as l};
