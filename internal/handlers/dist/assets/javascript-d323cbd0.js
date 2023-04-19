import{g as k}from"./index-b041cc2f.js";function G(i,u){for(var c=0;c<u.length;c++){const t=u[c];if(typeof t!="string"&&!Array.isArray(t)){for(const a in t)if(a!=="default"&&!(a in i)){const o=Object.getOwnPropertyDescriptor(t,a);o&&Object.defineProperty(i,a,o.get?o:{enumerable:!0,get:()=>t[a]})}}}return Object.freeze(Object.defineProperty(i,Symbol.toStringTag,{value:"Module"}))}var A,L;function K(){if(L)return A;L=1;const i="[A-Za-z$_][0-9A-Za-z$_]*",u=["as","in","of","if","for","while","finally","var","new","function","do","return","void","else","break","catch","instanceof","with","throw","case","default","try","switch","continue","typeof","delete","let","yield","const","class","debugger","async","await","static","import","from","export","extends"],c=["true","false","null","undefined","NaN","Infinity"],t=["Intl","DataView","Number","Math","Date","String","RegExp","Object","Function","Boolean","Error","Symbol","Set","Map","WeakSet","WeakMap","Proxy","Reflect","JSON","Promise","Float64Array","Int16Array","Int32Array","Int8Array","Uint16Array","Uint32Array","Float32Array","Array","Uint8Array","Uint8ClampedArray","ArrayBuffer","BigInt64Array","BigUint64Array","BigInt"],a=["EvalError","InternalError","RangeError","ReferenceError","SyntaxError","TypeError","URIError"],o=["setInterval","setTimeout","clearInterval","clearTimeout","require","exports","eval","isFinite","isNaN","parseFloat","parseInt","decodeURI","decodeURIComponent","encodeURI","encodeURIComponent","escape","unescape"],B=["arguments","this","super","console","window","document","localStorage","module","global"],U=[].concat(o,B,t,a);function P(e){return e?typeof e=="string"?e:e.source:null}function p(e){return f("(?=",e,")")}function f(...e){return e.map(n=>P(n)).join("")}function $(e){const y=(s,{after:d})=>{const _="</"+s[0].slice(1);return s.input.indexOf(_,d)!==-1},n=i,m={begin:"<>",end:"</>"},g={begin:/<[A-Za-z0-9\\._:-]+/,end:/\/[A-Za-z0-9\\._:-]+>|\/>/,isTrulyOpeningTag:(s,d)=>{const _=s[0].length+s.index,S=s.input[_];if(S==="<"){d.ignoreMatch();return}S===">"&&(y(s,{after:_})||d.ignoreMatch())}},r={$pattern:i,keyword:u,literal:c,built_in:U},O="[0-9](_?[0-9])*",E=`\\.(${O})`,R="0|[1-9](_?[0-9])*|0[0-7]*[89][0-9]*",I={className:"number",variants:[{begin:`(\\b(${R})((${E})|\\.)?|(${E}))[eE][+-]?(${O})\\b`},{begin:`\\b(${R})\\b((${E})\\b|\\.)?|(${E})\\b`},{begin:"\\b(0|[1-9](_?[0-9])*)n\\b"},{begin:"\\b0[xX][0-9a-fA-F](_?[0-9a-fA-F])*n?\\b"},{begin:"\\b0[bB][0-1](_?[0-1])*n?\\b"},{begin:"\\b0[oO][0-7](_?[0-7])*n?\\b"},{begin:"\\b0[0-7]+n?\\b"}],relevance:0},l={className:"subst",begin:"\\$\\{",end:"\\}",keywords:r,contains:[]},M={begin:"html`",end:"",starts:{end:"`",returnEnd:!1,contains:[e.BACKSLASH_ESCAPE,l],subLanguage:"xml"}},v={begin:"css`",end:"",starts:{end:"`",returnEnd:!1,contains:[e.BACKSLASH_ESCAPE,l],subLanguage:"css"}},D={className:"string",begin:"`",end:"`",contains:[e.BACKSLASH_ESCAPE,l]},T={className:"comment",variants:[e.COMMENT(/\/\*\*(?!\/)/,"\\*/",{relevance:0,contains:[{className:"doctag",begin:"@[A-Za-z]+",contains:[{className:"type",begin:"\\{",end:"\\}",relevance:0},{className:"variable",begin:n+"(?=\\s*(-)|$)",endsParent:!0,relevance:0},{begin:/(?=[^\n])\s/,relevance:0}]}]}),e.C_BLOCK_COMMENT_MODE,e.C_LINE_COMMENT_MODE]},C=[e.APOS_STRING_MODE,e.QUOTE_STRING_MODE,M,v,D,I,e.REGEXP_MODE];l.contains=C.concat({begin:/\{/,end:/\}/,keywords:r,contains:["self"].concat(C)});const w=[].concat(T,l.contains),N=w.concat([{begin:/\(/,end:/\)/,keywords:r,contains:["self"].concat(w)}]),b={className:"params",begin:/\(/,end:/\)/,excludeBegin:!0,excludeEnd:!0,keywords:r,contains:N};return{name:"Javascript",aliases:["js","jsx","mjs","cjs"],keywords:r,exports:{PARAMS_CONTAINS:N},illegal:/#(?![$_A-z])/,contains:[e.SHEBANG({label:"shebang",binary:"node",relevance:5}),{label:"use_strict",className:"meta",relevance:10,begin:/^\s*['"]use (strict|asm)['"]/},e.APOS_STRING_MODE,e.QUOTE_STRING_MODE,M,v,D,T,I,{begin:f(/[{,\n]\s*/,p(f(/(((\/\/.*$)|(\/\*(\*[^/]|[^*])*\*\/))\s*)*/,n+"\\s*:"))),relevance:0,contains:[{className:"attr",begin:n+p("\\s*:"),relevance:0}]},{begin:"("+e.RE_STARTERS_RE+"|\\b(case|return|throw)\\b)\\s*",keywords:"return throw case",contains:[T,e.REGEXP_MODE,{className:"function",begin:"(\\([^()]*(\\([^()]*(\\([^()]*\\)[^()]*)*\\)[^()]*)*\\)|"+e.UNDERSCORE_IDENT_RE+")\\s*=>",returnBegin:!0,end:"\\s*=>",contains:[{className:"params",variants:[{begin:e.UNDERSCORE_IDENT_RE,relevance:0},{className:null,begin:/\(\s*\)/,skip:!0},{begin:/\(/,end:/\)/,excludeBegin:!0,excludeEnd:!0,keywords:r,contains:N}]}]},{begin:/,/,relevance:0},{className:"",begin:/\s/,end:/\s*/,skip:!0},{variants:[{begin:m.begin,end:m.end},{begin:g.begin,"on:begin":g.isTrulyOpeningTag,end:g.end}],subLanguage:"xml",contains:[{begin:g.begin,end:g.end,skip:!0,contains:["self"]}]}],relevance:0},{className:"function",beginKeywords:"function",end:/[{;]/,excludeEnd:!0,keywords:r,contains:["self",e.inherit(e.TITLE_MODE,{begin:n}),b],illegal:/%/},{beginKeywords:"while if switch catch for"},{className:"function",begin:e.UNDERSCORE_IDENT_RE+"\\([^()]*(\\([^()]*(\\([^()]*\\)[^()]*)*\\)[^()]*)*\\)\\s*\\{",returnBegin:!0,contains:[b,e.inherit(e.TITLE_MODE,{begin:n})]},{variants:[{begin:"\\."+n},{begin:"\\$"+n}],relevance:0},{className:"class",beginKeywords:"class",end:/[{;=]/,excludeEnd:!0,illegal:/[:"[\]]/,contains:[{beginKeywords:"extends"},e.UNDERSCORE_TITLE_MODE]},{begin:/\b(?=constructor)/,end:/[{;]/,excludeEnd:!0,contains:[e.inherit(e.TITLE_MODE,{begin:n}),"self",b]},{begin:"(get|set)\\s+(?="+n+"\\()",end:/\{/,keywords:"get set",contains:[e.inherit(e.TITLE_MODE,{begin:n}),{begin:/\(\)/},b]},{begin:/\$[(.]/}]}}return A=$,A}var x=K();const F=k(x),H=G({__proto__:null,default:F},[x]);export{H as j};
