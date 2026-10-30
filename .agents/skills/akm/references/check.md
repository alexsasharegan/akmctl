# check

Use it before every commit and after every record edit. Exit code 0 is clean or warnings only, 1 means errors. Read-only.

```
$ akmctl check
docs/work/261029-atxk-add-json-output.md: [W] CL1 active item has no claimed_by
docs/work/261029-atxk-add-json-output.md: [E] WK6 status active requires a non-empty ## Next
1 errors, 2 warnings
```

Each finding is `[E|W] <rule> <message>`; fix the file the line names, then run it again. Rule ids are in the AKM specification (`docs/spec`).

Files changed: none.

Errors: `unreadable` when the manifest or a record cannot be parsed in a way that stops the check.
