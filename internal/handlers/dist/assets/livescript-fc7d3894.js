import{g as v}from"./index-b37b5403.js";function B(a,o){for(var i=0;i<o.length;i++){const n=o[i];if(typeof n!="string"&&!Array.isArray(n)){for(const t in n)if(t!=="default"&&!(t in a)){const s=Object.getOwnPropertyDescriptor(n,t);s&&Object.defineProperty(a,t,s.get?s:{enumerable:!0,get:()=>n[t]})}}}return Object.freeze(Object.defineProperty(a,Symbol.toStringTag,{value:"Module"}))}var g,E;function R(){if(E)return g;E=1;const a=["as","in","of","if","for","while","finally","var","new","function","do","return","void","else","break","catch","instanceof","with","throw","case","default","try","switch","continue","typeof","delete","let","yield","const","class","debugger","async","await","static","import","from","export","extends"],o=["true","false","null","undefined","NaN","Infinity"],i=["Intl","DataView","Number","Math","Date","String","RegExp","Object","Function","Boolean","Error","Symbol","Set","Map","WeakSet","WeakMap","Proxy","Reflect","JSON","Promise","Float64Array","Int16Array","Int32Array","Int8Array","Uint16Array","Uint32Array","Float32Array","Array","Uint8Array","Uint8ClampedArray","ArrayBuffer","BigInt64Array","BigUint64Array","BigInt"],n=["EvalError","InternalError","RangeError","ReferenceError","SyntaxError","TypeError","URIError"],t=["setInterval","setTimeout","clearInterval","clearTimeout","require","exports","eval","isFinite","isNaN","parseFloat","parseInt","decodeURI","decodeURIComponent","encodeURI","encodeURIComponent","escape","unescape"],s=["arguments","this","super","console","window","document","localStorage","module","global"],A=[].concat(t,s,i,n);function b(e){const p=["npm","print"],f=["yes","no","on","off","it","that","void"],y=["then","unless","until","loop","of","by","when","and","or","is","isnt","not","it","that","otherwise","from","to","til","fallthrough","case","enum","native","list","map","__hasProp","__extends","__slice","__bind","__indexOf"],c={keyword:a.concat(y),literal:o.concat(f),built_in:A.concat(p)},r="[A-Za-z$_](?:-[0-9A-Za-z$_]|[0-9A-Za-z$_])*",d=e.inherit(e.TITLE_MODE,{begin:r}),l={className:"subst",begin:/#\{/,end:/\}/,keywords:c},_={className:"subst",begin:/#[A-Za-z$_]/,end:/(?:-[0-9A-Za-z$_]|[0-9A-Za-z$_])*/,keywords:c},u=[e.BINARY_NUMBER_MODE,{className:"number",begin:"(\\b0[xX][a-fA-F0-9_]+)|(\\b\\d(\\d|_\\d)*(\\.(\\d(\\d|_\\d)*)?)?(_*[eE]([-+]\\d(_\\d|\\d)*)?)?[_a-z]*)",relevance:0,starts:{end:"(\\s*/)?",relevance:0}},{className:"string",variants:[{begin:/'''/,end:/'''/,contains:[e.BACKSLASH_ESCAPE]},{begin:/'/,end:/'/,contains:[e.BACKSLASH_ESCAPE]},{begin:/"""/,end:/"""/,contains:[e.BACKSLASH_ESCAPE,l,_]},{begin:/"/,end:/"/,contains:[e.BACKSLASH_ESCAPE,l,_]},{begin:/\\/,end:/(\s|$)/,excludeEnd:!0}]},{className:"regexp",variants:[{begin:"//",end:"//[gim]*",contains:[l,e.HASH_COMMENT_MODE]},{begin:/\/(?![ *])(\\.|[^\\\n])*?\/[gim]*(?=\W)/}]},{begin:"@"+r},{begin:"``",end:"``",excludeBegin:!0,excludeEnd:!0,subLanguage:"javascript"}];l.contains=u;const I={className:"params",begin:"\\(",returnBegin:!0,contains:[{begin:/\(/,end:/\)/,keywords:c,contains:["self"].concat(u)}]},m={begin:"(#=>|=>|\\|>>|-?->|!->)"};return{name:"LiveScript",aliases:["ls"],keywords:c,illegal:/\/\*/,contains:u.concat([e.COMMENT("\\/\\*","\\*\\/"),e.HASH_COMMENT_MODE,m,{className:"function",contains:[d,I],returnBegin:!0,variants:[{begin:"("+r+"\\s*(?:=|:=)\\s*)?(\\(.*\\)\\s*)?\\B->\\*?",end:"->\\*?"},{begin:"("+r+"\\s*(?:=|:=)\\s*)?!?(\\(.*\\)\\s*)?\\B[-~]{1,2}>\\*?",end:"[-~]{1,2}>\\*?"},{begin:"("+r+"\\s*(?:=|:=)\\s*)?(\\(.*\\)\\s*)?\\B!?[-~]{1,2}>\\*?",end:"!?[-~]{1,2}>\\*?"}]},{className:"class",beginKeywords:"class",end:"$",illegal:/[:="\[\]]/,contains:[{beginKeywords:"extends",endsWithParent:!0,illegal:/[:="\[\]]/,contains:[d]},d]},{begin:r+":",end:":",returnBegin:!0,returnEnd:!0,relevance:0}])}}return g=b,g}var S=R();const L=v(S),T=B({__proto__:null,default:L},[S]);export{T as l};
