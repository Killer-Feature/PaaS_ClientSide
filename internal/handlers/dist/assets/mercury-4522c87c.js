import{g as u}from"./index-b37b5403.js";function m(i,e){for(var n=0;n<e.length;n++){const r=e[n];if(typeof r!="string"&&!Array.isArray(r)){for(const t in r)if(t!=="default"&&!(t in i)){const a=Object.getOwnPropertyDescriptor(r,t);a&&Object.defineProperty(i,t,a.get?a:{enumerable:!0,get:()=>r[t]})}}}return Object.freeze(Object.defineProperty(i,Symbol.toStringTag,{value:"Module"}))}var _,s;function f(){if(s)return _;s=1;function i(e){const n={keyword:"module use_module import_module include_module end_module initialise mutable initialize finalize finalise interface implementation pred mode func type inst solver any_pred any_func is semidet det nondet multi erroneous failure cc_nondet cc_multi typeclass instance where pragma promise external trace atomic or_else require_complete_switch require_det require_semidet require_multi require_nondet require_cc_multi require_cc_nondet require_erroneous require_failure",meta:"inline no_inline type_spec source_file fact_table obsolete memo loop_check minimal_model terminates does_not_terminate check_termination promise_equivalent_clauses foreign_proc foreign_decl foreign_code foreign_type foreign_import_module foreign_export_enum foreign_export foreign_enum may_call_mercury will_not_call_mercury thread_safe not_thread_safe maybe_thread_safe promise_pure promise_semipure tabled_for_io local untrailed trailed attach_to_io_state can_pass_as_mercury_type stable will_not_throw_exception may_modify_trail will_not_modify_trail may_duplicate may_not_duplicate affects_liveness does_not_affect_liveness doesnt_affect_liveness no_sharing unknown_sharing sharing",built_in:"some all not if then else true fail false try catch catch_any semidet_true semidet_false semidet_fail impure_true impure semipure"},r=e.COMMENT("%","$"),t={className:"number",begin:"0'.\\|0[box][0-9a-fA-F]*"},a=e.inherit(e.APOS_STRING_MODE,{relevance:0}),o=e.inherit(e.QUOTE_STRING_MODE,{relevance:0}),l={className:"subst",begin:"\\\\[abfnrtv]\\|\\\\x[0-9a-fA-F]*\\\\\\|%[-+# *.0-9]*[dioxXucsfeEgGp]",relevance:0};return o.contains=o.contains.slice(),o.contains.push(l),{name:"Mercury",aliases:["m","moo"],keywords:n,contains:[{className:"built_in",variants:[{begin:"<=>"},{begin:"<=",relevance:0},{begin:"=>",relevance:0},{begin:"/\\\\"},{begin:"\\\\/"}]},{className:"built_in",variants:[{begin:":-\\|-->"},{begin:"=",relevance:0}]},r,e.C_BLOCK_COMMENT_MODE,t,e.NUMBER_MODE,a,o,{begin:/:-/},{begin:/\.$/}]}}return _=i,_}var c=f();const d=u(c),b=m({__proto__:null,default:d},[c]);export{b as m};
