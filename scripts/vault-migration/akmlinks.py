"""Shared helpers for the vault migration scripts: file listing, code masking, link finding and
resolution as Obsidian does it. Standard library only. Operates on bytes so that no file is decoded
or re-encoded (Shavian text and everything else passes through untouched).
"""
import os, re, subprocess, unicodedata
from urllib.parse import unquote

def list_md(root):
    """Markdown files of the work tree (tracked or untracked, not ignored) that exist, repo-relative."""
    out = subprocess.run(["git", "-C", root, "ls-files", "-z", "-c", "-o", "--exclude-standard"],
                         capture_output=True, check=True).stdout.split(b"\0")
    res = []
    for b in out:
        if not b:
            continue
        p = b.decode("utf-8")
        if p.lower().endswith(".md") and os.path.isfile(os.path.join(root, p)):
            res.append(p)
    return sorted(set(res))

def norm(s):
    return unicodedata.normalize("NFC", s).casefold()

class Index:
    """Name index of the Markdown files, for wikilink resolution (case-insensitive, as Obsidian)."""
    def __init__(self, paths):
        self.by_name = {}
        self.paths = set(paths)
        for p in paths:
            base = p.rsplit("/", 1)[-1][:-3]
            self.by_name.setdefault(norm(base), []).append(p)
        self.lower_paths = {norm(p[:-3]): p for p in paths}

    def resolve_wiki(self, target, source):
        """Return (kind, path-or-list): kind in ok, ambiguous, none, other."""
        t = target.strip()
        if t == "":
            return ("self", source)
        if t.lower().endswith(".md"):
            t = t[:-3]
        elif "." in t.rsplit("/", 1)[-1] and not t.lower().endswith(".md"):
            ext = t.rsplit(".", 1)[-1].lower()
            if ext in ("png", "jpg", "jpeg", "gif", "svg", "pdf", "mp3", "m4a", "mp4", "webp", "base", "canvas", "json", "csv", "py", "sh", "txt"):
                return ("other", t)
        key = norm(t)
        if "/" in key:
            hits = [p for lp, p in self.lower_paths.items() if lp == key or lp.endswith("/" + key)]
        else:
            hits = self.by_name.get(key, [])
        if len(hits) == 1:
            return ("ok", hits[0])
        if not hits:
            return ("none", t)
        # Obsidian prefers the file in the same folder, then the shortest path; report it as ambiguous.
        return ("ambiguous", sorted(hits))

    def resolve_md(self, dest, source):
        d = unquote(dest.strip())
        if d.startswith("<") and d.endswith(">"):
            d = d[1:-1]
        d = d.split("#", 1)[0]
        if not d.lower().endswith(".md") or "://" in d:
            return None
        base = os.path.dirname(source)
        p = os.path.normpath(os.path.join(base, d)).replace(os.sep, "/")
        for cand in (p, d):
            if cand in self.paths:
                return cand
        return None

# ---- code masking ----
FENCE = re.compile(rb"^[ \t]{0,3}(`{3,}|~{3,})")

def frontmatter_end(b):
    """Offset just after the closing --- line, or 0."""
    if not b.startswith(b"---\n") and not b.startswith(b"---\r\n"):
        return 0
    m = re.search(rb"\n---[ \t]*\r?\n", b[3:])
    return 3 + m.end() if m else 0

def code_spans(b):
    """Byte ranges of fenced code blocks and inline code spans (outside the frontmatter)."""
    spans = []
    start = frontmatter_end(b)
    pos = start
    fence = None
    fstart = 0
    in_text = []  # (a, b) ranges outside fences
    cur = start
    for m in re.finditer(rb"[^\n]*\n|[^\n]+$", b[start:]):
        line_start = start + m.start()
        line = m.group(0)
        fm = FENCE.match(line)
        if fence is None:
            if fm:
                fence = (fm.group(1)[:1], len(fm.group(1)))
                fstart = line_start
            else:
                in_text.append((line_start, line_start + len(line)))
        else:
            if fm and fm.group(1)[:1] == fence[0] and len(fm.group(1)) >= fence[1] and line.strip(b"\r\n \t") == fm.group(1):
                spans.append((fstart, line_start + len(line)))
                fence = None
    if fence is not None:
        spans.append((fstart, len(b)))
    # inline code spans within the non-fenced text, per paragraph
    para = []
    def flush():
        if not para:
            return
        a, z = para[0][0], para[-1][1]
        seg = b[a:z]
        i = 0
        while i < len(seg):
            if seg[i:i+1] == b"`":
                j = i
                while j < len(seg) and seg[j:j+1] == b"`":
                    j += 1
                n = j - i
                k = j
                found = None
                while k < len(seg):
                    if seg[k:k+1] == b"`":
                        e = k
                        while e < len(seg) and seg[e:e+1] == b"`":
                            e += 1
                        if e - k == n:
                            found = e
                            break
                        k = e
                    else:
                        k += 1
                if found:
                    spans.append((a + i, a + found))
                    i = found
                    continue
                i = j
            else:
                i += 1
        para.clear()
    for a, z in in_text:
        if b[a:z].strip() == b"":
            flush()
        else:
            para.append((a, z))
    flush()
    return sorted(spans)

def in_spans(pos, spans):
    for a, z in spans:
        if a <= pos < z:
            return True
        if a > pos:
            break
    return False

WIKI = re.compile(rb"(!?)\[\[([^\[\]\n]+?)\]\]")
MDLINK = re.compile(rb"\]\((<[^<>\n]*>|[^()\s]+)(?:\s+\"[^\"\n]*\")?\)")

def split_inner(inner):
    """inner of [[...]] -> (target, rest) where rest starts at the first of # ^ | or \\| ."""
    m = re.search(rb"\\?\||[#^]", inner)
    if not m:
        return inner, b""
    # a backslash directly before | belongs to rest
    i = m.start()
    return inner[:i], inner[i:]

def find_links(b):
    """List of dicts {kind, start, end, tstart, tend, target(bytes)} for wikilinks, embeds and md links outside code."""
    spans = code_spans(b)
    out = []
    for m in WIKI.finditer(b):
        if in_spans(m.start(), spans):
            continue
        inner = m.group(2)
        target, rest = split_inner(inner)
        ts = m.start(2)
        out.append(dict(kind="wiki", start=m.start(), end=m.end(), tstart=ts, tend=ts + len(target), target=target, rest=rest))
    for m in MDLINK.finditer(b):
        if in_spans(m.start(), spans):
            continue
        out.append(dict(kind="md", start=m.start(), end=m.end(), tstart=m.start(1), tend=m.end(1), target=m.group(1), rest=b""))
    out.sort(key=lambda d: d["start"])
    return out

SHAVIAN = re.compile("[\U00010450-\U0001047F]".encode("utf-8"))
def has_shavian(b):
    return bool(SHAVIAN.search(b))
