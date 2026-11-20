// Copy text to the clipboard. The Clipboard API works on a secure context, and
// the workbench is served from localhost, which counts as one; where it is
// absent or refuses (a page opened over a plain network address, a browser that
// blocks it), a selected hidden textarea and execCommand("copy") is the fallback.

/** Resolves true when the text was copied, false when no way worked. */
export async function copyText(text) {
  try {
    if (typeof navigator !== "undefined" && navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text);
      return true;
    }
  } catch (e) {
    /* fall through to the textarea */
  }
  try {
    const ta = document.createElement("textarea");
    ta.value = text;
    ta.setAttribute("readonly", "");
    ta.style.position = "fixed";
    ta.style.opacity = "0";
    document.body.appendChild(ta);
    ta.select();
    const ok = typeof document.execCommand === "function" && document.execCommand("copy");
    ta.remove();
    return !!ok;
  } catch (e) {
    return false;
  }
}
