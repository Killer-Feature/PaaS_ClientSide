import{g as h}from"./index-b37b5403.js";function x(o,d){for(var c=0;c<d.length;c++){const n=d[c];if(typeof n!="string"&&!Array.isArray(n)){for(const t in n)if(t!=="default"&&!(t in o)){const i=Object.getOwnPropertyDescriptor(n,t);i&&Object.defineProperty(o,t,i.get?i:{enumerable:!0,get:()=>n[t]})}}}return Object.freeze(Object.defineProperty(o,Symbol.toStringTag,{value:"Module"}))}var b,N;function L(){if(N)return b;N=1;function o(n){return n?typeof n=="string"?n:n.source:null}function d(...n){return"("+n.map(i=>o(i)).join("|")+")"}function c(n){const t=d(...["(?:NeedsTeXFormat|RequirePackage|GetIdInfo)","Provides(?:Expl)?(?:Package|Class|File)","(?:DeclareOption|ProcessOptions)","(?:documentclass|usepackage|input|include)","makeat(?:letter|other)","ExplSyntax(?:On|Off)","(?:new|renew|provide)?command","(?:re)newenvironment","(?:New|Renew|Provide|Declare)(?:Expandable)?DocumentCommand","(?:New|Renew|Provide|Declare)DocumentEnvironment","(?:(?:e|g|x)?def|let)","(?:begin|end)","(?:part|chapter|(?:sub){0,2}section|(?:sub)?paragraph)","caption","(?:label|(?:eq|page|name)?ref|(?:paren|foot|super)?cite)","(?:alpha|beta|[Gg]amma|[Dd]elta|(?:var)?epsilon|zeta|eta|[Tt]heta|vartheta)","(?:iota|(?:var)?kappa|[Ll]ambda|mu|nu|[Xx]i|[Pp]i|varpi|(?:var)rho)","(?:[Ss]igma|varsigma|tau|[Uu]psilon|[Pp]hi|varphi|chi|[Pp]si|[Oo]mega)","(?:frac|sum|prod|lim|infty|times|sqrt|leq|geq|left|right|middle|[bB]igg?)","(?:[lr]angle|q?quad|[lcvdi]?dots|d?dot|hat|tilde|bar)"].map(e=>e+"(?![a-zA-Z@:_])")),i=new RegExp(["(?:__)?[a-zA-Z]{2,}_[a-zA-Z](?:_?[a-zA-Z])+:[a-zA-Z]*","[lgc]__?[a-zA-Z](?:_?[a-zA-Z])*_[a-zA-Z]{2,}","[qs]__?[a-zA-Z](?:_?[a-zA-Z])+","use(?:_i)?:[a-zA-Z]*","(?:else|fi|or):","(?:if|cs|exp):w","(?:hbox|vbox):n","::[a-zA-Z]_unbraced","::[a-zA-Z:]"].map(e=>e+"(?![a-zA-Z:_])").join("|")),T=[{begin:/[a-zA-Z@]+/},{begin:/[^a-zA-Z@]?/}],v=[{begin:/\^{6}[0-9a-f]{6}/},{begin:/\^{5}[0-9a-f]{5}/},{begin:/\^{4}[0-9a-f]{4}/},{begin:/\^{3}[0-9a-f]{3}/},{begin:/\^{2}[0-9a-f]{2}/},{begin:/\^{2}[\u0000-\u007f]/}],M={className:"keyword",begin:/\\/,relevance:0,contains:[{endsParent:!0,begin:t},{endsParent:!0,begin:i},{endsParent:!0,variants:v},{endsParent:!0,relevance:0,variants:T}]},O={className:"params",relevance:0,begin:/#+\d?/},P={variants:v},C={className:"built_in",relevance:0,begin:/[$&^_]/},z={className:"meta",begin:"% !TeX",end:"$",relevance:10},I=n.COMMENT("%","$",{relevance:0}),_=[M,O,P,C,z,I],A={begin:/\{/,end:/\}/,relevance:0,contains:["self",..._]},B=n.inherit(A,{relevance:0,endsParent:!0,contains:[A,..._]}),D={begin:/\[/,end:/\]/,endsParent:!0,relevance:0,contains:[A,..._]},f={begin:/\s+/,relevance:0},l=[B],m=[D],a=function(e,r){return{contains:[f],starts:{relevance:0,contains:e,starts:r}}},s=function(e,r){return{begin:"\\\\"+e+"(?![a-zA-Z@:_])",keywords:{$pattern:/\\[a-zA-Z]+/,keyword:"\\"+e},relevance:0,contains:[f],starts:r}},g=function(e,r){return n.inherit({begin:"\\\\begin(?=[ 	]*(\\r?\\n[ 	]*)?\\{"+e+"\\})",keywords:{$pattern:/\\[a-zA-Z]+/,keyword:"\\begin"},relevance:0},a(l,r))},p=(e="string")=>n.END_SAME_AS_BEGIN({className:e,begin:/(.|\r?\n)/,end:/(.|\r?\n)/,excludeBegin:!0,excludeEnd:!0,endsParent:!0}),E=function(e){return{className:"string",end:"(?=\\\\end\\{"+e+"\\})"}},u=(e="string")=>({relevance:0,begin:/\{/,starts:{endsParent:!0,contains:[{className:e,end:/(?=\})/,endsParent:!0,contains:[{begin:/\{/,end:/\}/,relevance:0,contains:["self"]}]}]}}),Z=[...["verb","lstinline"].map(e=>s(e,{contains:[p()]})),s("mint",a(l,{contains:[p()]})),s("mintinline",a(l,{contains:[u(),p()]})),s("url",{contains:[u("link"),u("link")]}),s("hyperref",{contains:[u("link")]}),s("href",a(m,{contains:[u("link")]})),...[].concat(...["","\\*"].map(e=>[g("verbatim"+e,E("verbatim"+e)),g("filecontents"+e,a(l,E("filecontents"+e))),...["","B","L"].map(r=>g(r+"Verbatim"+e,a(m,E(r+"Verbatim"+e))))])),g("minted",a(m,a(l,E("minted"))))];return{name:"LaTeX",aliases:["tex"],contains:[...Z,..._]}}return b=c,b}var R=L();const y=h(R),k=x({__proto__:null,default:y},[R]);export{k as l};
