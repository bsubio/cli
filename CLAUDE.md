# Code principles

You're the best Golang and systems expert in the world.

# Architecture

Use standard library whenever possible.
Use flag library if possible.
Use log library for logging.

Apply security principles to the code.

Add -verbose for basic echoing of requests.
Add -debug for more debugging output.
For any server (listening) commands, print : "Listening: http://127.0.0.1:<port>".

Whatever you change, test cmd/bsubio/static/*.md files for accuracy

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

Skip "Good" and "Positive" stuff.
Try to give crisp, short, review. Save words.
Give only ONE best suggestion.
Use golang-pro.md as golang agent reviewer.
Use API designer (api-designer.md) agent for API design.
Use architect agent to review (architect-reviewer.md).
Use cli-developer.md for CLI changes review.
Use code-reviewer.md to review code.
Repo follows Conventional Commits.
PR title should start with "Fix:" or "Feat:".
Don't use goto's unless some mutex/double-nested loops are involved.

# Tooling

Run:
- `make check` to lint/format all .go files at the end of each change,
- `make` to build

# Reviews

Keep reviews short and concise.

# Last words

Feedback phrase like Rob Pike and Ken Thompson.
