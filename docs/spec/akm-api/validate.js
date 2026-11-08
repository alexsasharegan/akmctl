const Ajv2020 = require("ajv/dist/2020");
const fs = require("fs");
const path = require("path");
const root = process.argv[2];
const sd = path.join(root, "schemas"),
  cd = path.join(root, ".fixtures", "cases");
const ajv = new Ajv2020({ strict: true, allErrors: true });
const files = fs.readdirSync(sd).filter((f) => f.endsWith(".json"));
for (const f of files)
  ajv.addSchema(JSON.parse(fs.readFileSync(path.join(sd, f))));
let bad = 0;

// Matchers (ADR 261108-15dw): the expected response of a case may hold $id(..),
// $hash and $any. The first two become sample values that fit the schema's
// patterns; each $any becomes the first sample value the schema accepts at
// that position (the value is not part of the contract).
const samples = ["0123456789abcdef", [], {}, "x", 0, true, "0".repeat(40)];
function fillMarkers(v, anys, at) {
  if (typeof v === "string") {
    if (v === "$any") {
      anys.push(at);
      return v;
    }
    const ids = {};
    return v
      .replace(/\$id\((\d{6})(?:,(\w+))?\)/g, (_, d, l) => {
        const k = d + "," + (l || "");
        ids[k] = ids[k] || String(Object.keys(ids).length).padStart(4, "0");
        return d + "-" + ids[k];
      })
      .replace(/\$hash/g, "0".repeat(40));
  }
  if (Array.isArray(v))
    return v.map((x, i) => fillMarkers(x, anys, at + "/" + i));
  if (v && typeof v === "object") {
    const o = {};
    for (const k of Object.keys(v))
      o[k] = fillMarkers(v[k], anys, at + "/" + k);
    return o;
  }
  return v;
}
function setAt(root, pointer, value) {
  const parts = pointer.split("/").slice(1);
  let o = root;
  for (const p of parts.slice(0, -1)) o = o[p];
  o[parts[parts.length - 1]] = value;
}
const bindingNames = ["api", "http", "cli"];
for (const f of fs.readdirSync(cd).sort()) {
  const c = JSON.parse(fs.readFileSync(path.join(cd, f)));
  // the operation is the file name's prefix (set-15.json is a set case): the
  // response of a refusal that never reached the operation names "http"
  const op = f.split("-")[0];
  for (const b of c.bindings || [])
    if (!bindingNames.includes(b)) {
      bad++;
      console.log(f, "bindings: unknown binding", b);
    }
  for (const g of c.git || []) {
    if (
      typeof g.message !== "string" ||
      typeof g.date !== "string" ||
      !g.files ||
      typeof g.files !== "object"
    ) {
      bad++;
      console.log(f, "git: each commit needs message, date and files");
    }
  }
  const base = "https://akm.invalid/api/1.0/" + op + ".json";
  const vreq = ajv.compile({ $ref: base + "#/$defs/request" }),
    vres = ajv.compile({ $ref: base + "#/$defs/response" });
  const anys = [];
  const expected = fillMarkers(c.expected, anys, "");
  for (const at of anys) {
    let ok = false;
    for (const s of samples) {
      setAt(expected, at, s);
      vres(expected);
      const here = (vres.errors || []).filter(
        (e) => e.instancePath === at || e.instancePath.startsWith(at + "/"),
      );
      if (here.length === 0) {
        ok = true;
        break;
      }
    }
    if (!ok) {
      bad++;
      console.log(f, "response", at, "accepts none of the sample values");
    }
  }
  const a = c.expected.ok === false ? true : vreq(c.request),
    b = vres(expected);
  if (!a) {
    bad++;
    console.log(f, "request", JSON.stringify(vreq.errors));
  }
  if (!b) {
    bad++;
    console.log(f, "response", JSON.stringify(vres.errors).slice(0, 600));
  }
}
console.log(
  files.length,
  "schemas compiled;",
  fs.readdirSync(cd).length,
  "cases;",
  bad,
  "failures",
);
