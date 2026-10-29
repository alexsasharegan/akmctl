const Ajv2020=require("ajv/dist/2020");const fs=require("fs");const path=require("path");
const root=process.argv[2];const sd=path.join(root,"schemas"),cd=path.join(root,".fixtures","cases");
const ajv=new Ajv2020({strict:true,allErrors:true});
const files=fs.readdirSync(sd).filter(f=>f.endsWith(".json"));
for(const f of files)ajv.addSchema(JSON.parse(fs.readFileSync(path.join(sd,f))));
let bad=0;
for(const f of fs.readdirSync(cd).sort()){
  const c=JSON.parse(fs.readFileSync(path.join(cd,f)));const op=c.expected.operation;
  const base="https://akm.invalid/api/1.0/"+op+".json";
  const vreq=ajv.compile({$ref:base+"#/$defs/request"}),vres=ajv.compile({$ref:base+"#/$defs/response"});
  const a=c.expected.ok===false?true:vreq(c.request),b=vres(c.expected);
  if(!a){bad++;console.log(f,"request",JSON.stringify(vreq.errors))}
  if(!b){bad++;console.log(f,"response",JSON.stringify(vres.errors).slice(0,600))}
}
console.log(files.length,"schemas compiled;",fs.readdirSync(cd).length,"cases;",bad,"failures");
