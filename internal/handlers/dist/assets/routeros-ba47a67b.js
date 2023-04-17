import{g as d}from"./index-b37b5403.js";function m(o,n){for(var t=0;t<n.length;t++){const e=n[t];if(typeof e!="string"&&!Array.isArray(e)){for(const r in e)if(r!=="default"&&!(r in o)){const i=Object.getOwnPropertyDescriptor(e,r);i&&Object.defineProperty(o,r,i.get?i:{enumerable:!0,get:()=>e[r]})}}}return Object.freeze(Object.defineProperty(o,Symbol.toStringTag,{value:"Module"}))}var a,c;function b(){if(c)return a;c=1;function o(n){const t="foreach do while for if from to step else on-error and or not in",e="global local beep delay put len typeof pick log time set find environment terminal error execute parse resolve toarray tobool toid toip toip6 tonum tostr totime",r="add remove enable disable set get print export edit find run debug error info warning",i="true false yes no nothing nil null",u="traffic-flow traffic-generator firewall scheduler aaa accounting address-list address align area bandwidth-server bfd bgp bridge client clock community config connection console customer default dhcp-client dhcp-server discovery dns e-mail ethernet filter firmware gps graphing group hardware health hotspot identity igmp-proxy incoming instance interface ip ipsec ipv6 irq l2tp-server lcd ldp logging mac-server mac-winbox mangle manual mirror mme mpls nat nd neighbor network note ntp ospf ospf-v3 ovpn-server page peer pim ping policy pool port ppp pppoe-client pptp-server prefix profile proposal proxy queue radius resource rip ripng route routing screen script security-profiles server service service-port settings shares smb sms sniffer snmp snooper socks sstp-server system tool tracking type upgrade upnp user-manager users user vlan secret vrrp watchdog web-access wireless pptp pppoe lan wan layer7-protocol lease simple raw",s={className:"variable",variants:[{begin:/\$[\w\d#@][\w\d_]*/},{begin:/\$\{(.*?)\}/}]},l={className:"string",begin:/"/,end:/"/,contains:[n.BACKSLASH_ESCAPE,s,{className:"variable",begin:/\$\(/,end:/\)/,contains:[n.BACKSLASH_ESCAPE]}]},p={className:"string",begin:/'/,end:/'/};return{name:"Microtik RouterOS script",aliases:["mikrotik"],case_insensitive:!0,keywords:{$pattern:/:?[\w-]+/,literal:i,keyword:t+" :"+t.split(" ").join(" :")+" :"+e.split(" ").join(" :")},contains:[{variants:[{begin:/\/\*/,end:/\*\//},{begin:/\/\//,end:/$/},{begin:/<\//,end:/>/}],illegal:/./},n.COMMENT("^#","$"),l,p,s,{begin:/[\w-]+=([^\s{}[\]()>]+)/,relevance:0,returnBegin:!0,contains:[{className:"attribute",begin:/[^=]+/},{begin:/=/,endsWithParent:!0,relevance:0,contains:[l,p,s,{className:"literal",begin:"\\b("+i.split(" ").join("|")+")\\b"},{begin:/("[^"]*"|[^\s{}[\]]+)/}]}]},{className:"number",begin:/\*[0-9a-fA-F]+/},{begin:"\\b("+r.split(" ").join("|")+")([\\s[(\\]|])",returnBegin:!0,contains:[{className:"builtin-name",begin:/\w+/}]},{className:"built_in",variants:[{begin:"(\\.\\./|/|\\s)(("+u.split(" ").join("|")+");?\\s)+"},{begin:/\.\./,relevance:0}]}]}}return a=o,a}var g=b();const f=d(g),y=m({__proto__:null,default:f},[g]);export{y as r};
