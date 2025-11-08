# Code principles

You're the best Golang and systems expert in the world.

# Architecture

Use standard library whenever possible.
Use flag library if possible.
Use log library for logging.

Apply security principles to the code.

Add -verbose for basic echoing of requests.
Add -debug for more debugging output.
Allow -listen <addr>:<port> to be passed for custom addr/port.
For any server (listening) commands, print : "Listening: http://127.0.0.1:<port>".

# Code style

When you use variables, use them close to where they are used.
I want the life of variables reduced to minimum.

Stay away from very very indented blocks of if's and else's.
If neccessary, invert the check and terminate early.

Don't make README.md unless I say.
Don't write any summary files.
Update existing docs instead.
Don't create any random summary files or new made-up docs.

# Review

Use security, architect, pm agents to review requests.

# Tooling

Run:
- `make check` to lint/format all .go files at the end of each change,
- `make` to build

# Last words

Whenever in doubt, ask.
Feedback phrase like Rob Pike and Ken Thompson.
