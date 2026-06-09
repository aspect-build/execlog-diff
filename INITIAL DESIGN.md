# execlog-diff

The purpose of execlog-diff is to take execlogs and perform some kind of operation on them in order to produce useful
output. To start with, this is going to be a diff between execlogs.

## The problem

The biggest problem that execlog-diff faces is the size of execlogs. 

A compact execlog might be ~20MB deconstructed and
compressed. When reconstructed and decompressed, that same ~20MB compact execlog might expand to many, many GB.

Now some of that can be attributed to unnecessary whitespace and information we dont care about (yet at least).
But consider that we need 2 of these execlogs to be able to do a diff and suddenly we are very constrained by runtime memory. We simply cannot
keep both execlogs in memory at the same time. This was attempted and we OOM'ed constantly (even with a small repo).

It should also be clarified that we do actually need to reconstruct these execlogs to a certain degree in order
to diff them. This leaves us with only a few options

## Considered solutions

1. We use the filesystem as an intermediate and go step by step. We decompress the execlogs into new files.
We process the files into more new files. We diff the files into a diff file. 
2. We "stream" the execlogs, doing all of this decompressing, processing & diffing live as we move through
the files. Then remove what we can, also live, In hopes that we can limit the amount of data in memory at 
any one time.

## Design

We decided to go with the "streaming" solution from above. The main "steps" that a given bazel action spawn
will take as it travels through execlog-diff are these:
1. Decompression
2. Unmarshalling / transformation into an internal object
3. Processing (reconstruction, credential removal, hashing, etc)
4. Added to the tree / graph
5. Compared / diffed between execlogs
6. Removed from memory OR found to be different and saved for later

All of these steps happen asynchronously as the execlogs are streamed through execlog-diff. The reason is that
in doing this, we can remove many execlog entries from memory before others are ever read into memory. 

Therefore while some individual areas of execlog-diff are not in themselves asynchronous bits of code,
they are being used in an asynchronous way and therefore much be treated as asynchronous. Only once we 
have finished processing both execlogs and are compiling the end results, can we run synchronously. 
