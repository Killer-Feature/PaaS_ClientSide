import{g as c}from"./index-b041cc2f.js";function f(n,t){for(var i=0;i<t.length;i++){const e=t[i];if(typeof e!="string"&&!Array.isArray(e)){for(const r in e)if(r!=="default"&&!(r in n)){const a=Object.getOwnPropertyDescriptor(e,r);a&&Object.defineProperty(n,r,a.get?a:{enumerable:!0,get:()=>e[r]})}}}return Object.freeze(Object.defineProperty(n,Symbol.toStringTag,{value:"Module"}))}var s,o;function u(){if(o)return s;o=1;function n(t){const i={className:"string",begin:/'(.|\\[xXuU][a-zA-Z0-9]+)'/},e={className:"string",variants:[{begin:'"',end:'"'}]},a={className:"function",beginKeywords:"def",end:/[:={\[(\n;]/,excludeEnd:!0,contains:[{className:"title",relevance:0,begin:/[^0-9\n\t "'(),.`{}\[\]:;][^\n\t "'(),.`{}\[\]:;]+|[^0-9\n\t "'(),.`{}\[\]:;=]/}]};return{name:"Flix",keywords:{literal:"true false",keyword:"case class def else enum if impl import in lat rel index let match namespace switch type yield with"},contains:[t.C_LINE_COMMENT_MODE,t.C_BLOCK_COMMENT_MODE,i,e,a,t.C_NUMBER_MODE]}}return s=n,s}var l=u();const d=c(l),m=f({__proto__:null,default:d},[l]);export{m as f};
