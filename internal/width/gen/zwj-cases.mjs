// Writes testdata/prettier-<version>-zwj.json: emoji ZWJ sequences, standard and not,
// with the width the pinned Prettier gives them (work item 261028-x1b4).
//
//   node internal/width/gen/zwj-cases.mjs   (from the repository root, after npm ci)
//
// Classes: rgi (a line of data/emoji-zwj-sequences.txt), unqualified (the same without
// U+FE0F), pair (two emoji joined by a ZWJ that are not a standard sequence), prefix (a
// standard sequence cut after an element), extended (a standard sequence plus ZWJ and
// one more emoji), tail (a standard sequence followed by a bare ZWJ).
import fs from "node:fs";
import prettier from "prettier";

const w = prettier.util.getStringWidth;
const lines = fs.readFileSync("internal/width/data/emoji-zwj-sequences.txt", "utf8").split("\n");
const seqs = [];
for (const l of lines) {
  const f = l.split("#")[0].split(";")[0].trim();
  if (f) seqs.push(String.fromCodePoint(...f.split(/\s+/).map((h) => parseInt(h, 16))));
}
const rgi = new Set(seqs);
const cases = new Map();
const add = (cls, s) => {
  if (s && !cases.has(s)) cases.set(s, cls);
};
for (const s of seqs) add("rgi", s);
for (const s of seqs) add("unqualified", s.replaceAll("️", ""));
const ZWJ = "‍";
// every element that stands in a standard sequence
const elems = [...new Set(seqs.flatMap((s) => s.split(ZWJ)))].sort();
// a fixed spread of pairs: each element with seven others, by index
for (let i = 0; i < elems.length; i++) {
  for (const k of [1, 31, 211]) {
    const s = elems[i] + ZWJ + elems[(i + k) % elems.length];
    if (!rgi.has(s)) add("pair", s);
  }
}
for (let i = 0; i < seqs.length; i += 5) {
  const parts = seqs[i].split(ZWJ);
  for (let n = 2; n < parts.length; n++) add("prefix", parts.slice(0, n).join(ZWJ));
  add("extended", seqs[i] + ZWJ + elems[(i * 13) % elems.length]);
  add("tail", seqs[i] + ZWJ);
}
const out = [...cases].map(([s, cls], id) => ({ id, cls, hex: Buffer.from(s, "utf8").toString("hex"), w: w(s) }));
fs.mkdirSync("internal/width/testdata", { recursive: true });
fs.writeFileSync(`internal/width/testdata/prettier-${prettier.version}-zwj.json`, "[\n" + out.map((c) => JSON.stringify(c)).join(",\n") + "\n]\n");
const tally = {};
for (const c of out) tally[c.cls + " w=" + c.w] = (tally[c.cls + " w=" + c.w] || 0) + 1;
console.log(out.length, "cases", tally);
