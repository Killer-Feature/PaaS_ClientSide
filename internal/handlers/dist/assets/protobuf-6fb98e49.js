import{g as f}from"./index-b041cc2f.js";function c(r,e){for(var n=0;n<e.length;n++){const t=e[n];if(typeof t!="string"&&!Array.isArray(t)){for(const o in t)if(o!=="default"&&!(o in r)){const i=Object.getOwnPropertyDescriptor(t,o);i&&Object.defineProperty(r,o,i.get?i:{enumerable:!0,get:()=>t[o]})}}}return Object.freeze(Object.defineProperty(r,Symbol.toStringTag,{value:"Module"}))}var u,s;function d(){if(s)return u;s=1;function r(e){return{name:"Protocol Buffers",keywords:{keyword:"package import option optional required repeated group oneof",built_in:"double float int32 int64 uint32 uint64 sint32 sint64 fixed32 fixed64 sfixed32 sfixed64 bool string bytes",literal:"true false"},contains:[e.QUOTE_STRING_MODE,e.NUMBER_MODE,e.C_LINE_COMMENT_MODE,e.C_BLOCK_COMMENT_MODE,{className:"class",beginKeywords:"message enum service",end:/\{/,illegal:/\n/,contains:[e.inherit(e.TITLE_MODE,{starts:{endsWithParent:!0,excludeEnd:!0}})]},{className:"function",beginKeywords:"rpc",end:/[{;]/,excludeEnd:!0,keywords:"rpc returns"},{begin:/^\s*[A-Z_]+(?=\s*=[^\n]+;$)/}]}}return u=r,u}var a=d();const p=f(a),b=c({__proto__:null,default:p},[a]);export{b as p};
