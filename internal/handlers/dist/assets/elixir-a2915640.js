import{g as f}from"./index-b37b5403.js";function m(a,e){for(var r=0;r<e.length;r++){const t=e[r];if(typeof t!="string"&&!Array.isArray(t)){for(const n in t)if(n!=="default"&&!(n in a)){const i=Object.getOwnPropertyDescriptor(t,n);i&&Object.defineProperty(a,n,i.get?i:{enumerable:!0,get:()=>t[n]})}}}return Object.freeze(Object.defineProperty(a,Symbol.toStringTag,{value:"Module"}))}var s,l;function A(){if(l)return s;l=1;function a(e){const r="[a-zA-Z_][a-zA-Z0-9_.]*(!|\\?)?",t="[a-zA-Z_]\\w*[!?=]?|[-+~]@|<<|>>|=~|===?|<=>|[<>]=?|\\*\\*|[-/+%^&*~`|]|\\[\\]=?",n={$pattern:r,keyword:"and false then defined module in return redo retry end for true self when next until do begin unless nil break not case cond alias while ensure or include use alias fn quote require import with|0"},i={className:"subst",begin:/#\{/,end:/\}/,keywords:n},o={className:"number",begin:"(\\b0o[0-7_]+)|(\\b0b[01_]+)|(\\b0x[0-9a-fA-F_]+)|(-?\\b[1-9][0-9_]*(\\.[0-9_]+([eE][-+]?[0-9]+)?)?)",relevance:0},c=`[/|([{<"']`,_={className:"string",begin:"~[a-z](?="+c+")",contains:[{endsParent:!0,contains:[{contains:[e.BACKSLASH_ESCAPE,i],variants:[{begin:/"/,end:/"/},{begin:/'/,end:/'/},{begin:/\//,end:/\//},{begin:/\|/,end:/\|/},{begin:/\(/,end:/\)/},{begin:/\[/,end:/\]/},{begin:/\{/,end:/\}/},{begin:/</,end:/>/}]}]}]},u={className:"string",begin:"~[A-Z](?="+c+")",contains:[{begin:/"/,end:/"/},{begin:/'/,end:/'/},{begin:/\//,end:/\//},{begin:/\|/,end:/\|/},{begin:/\(/,end:/\)/},{begin:/\[/,end:/\]/},{begin:/\{/,end:/\}/},{begin:/</,end:/>/}]},d={className:"string",contains:[e.BACKSLASH_ESCAPE,i],variants:[{begin:/"""/,end:/"""/},{begin:/'''/,end:/'''/},{begin:/~S"""/,end:/"""/,contains:[]},{begin:/~S"/,end:/"/,contains:[]},{begin:/~S'''/,end:/'''/,contains:[]},{begin:/~S'/,end:/'/,contains:[]},{begin:/'/,end:/'/},{begin:/"/,end:/"/}]},b={className:"function",beginKeywords:"def defp defmacro",end:/\B\b/,contains:[e.inherit(e.TITLE_MODE,{begin:r,endsParent:!0})]},S=e.inherit(b,{className:"class",beginKeywords:"defimpl defmodule defprotocol defrecord",end:/\bdo\b|$|;/}),g=[d,u,_,e.HASH_COMMENT_MODE,S,b,{begin:"::"},{className:"symbol",begin:":(?![\\s:])",contains:[d,{begin:t}],relevance:0},{className:"symbol",begin:r+":(?!:)",relevance:0},o,{className:"variable",begin:"(\\$\\W)|((\\$|@@?)(\\w+))"},{begin:"->"},{begin:"("+e.RE_STARTERS_RE+")\\s*",contains:[e.HASH_COMMENT_MODE,{begin:/\/: (?=\d+\s*[,\]])/,relevance:0,contains:[o]},{className:"regexp",illegal:"\\n",contains:[e.BACKSLASH_ESCAPE,i],variants:[{begin:"/",end:"/[a-z]*"},{begin:"%r\\[",end:"\\][a-z]*"}]}],relevance:0}];return i.contains=g,{name:"Elixir",keywords:n,contains:g}}return s=a,s}var E=A();const I=f(E),p=m({__proto__:null,default:I},[E]);export{p as e};
