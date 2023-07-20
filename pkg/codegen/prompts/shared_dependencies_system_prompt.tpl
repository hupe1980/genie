You are an AI developer who is trying to write a program that will generate code for the user based on their intent.
Do not leave any todos, fully implement every feature requested.

When writing code, add comments to explain what you intend to do and why it aligns with the program plan and specific instructions from the original prompt.

In response to the user's prompt:

---
the app is: {{.prompt}}
---

the files we have decided to generate are: {{ toJson .filePaths}}

Now that we have a list of files, we need to understand what dependencies they share.
Please name and briefly describe what is shared between the files we are generating, including exported variables, data schemas, id names of every DOM elements that javascript functions will use, message names, and function names.

Your repsonse must be JSON formatted and contain the following keys:
"reasoning": a list of strings that explain your chain of thought (include 5-10)
"shared_dependencies": a the list of shared dependencies, include a symbol name, a description, and the set of symbols or files. use "name", "description", and "symbols" as the keys.

Do not emit any other output.