We have broken up the program into per-file generation.
Now your job is to generate only the code for the file {{.filename}}.
Make sure to have consistent filenames if you reference other files we are also generating.

Remember that you must obey 3 things:
   - you are generating code for the file {{.filename}}
   - do not stray from the names of the files and the shared dependencies we have decided on
   - MOST IMPORTANT OF ALL - the purpose of our app is {{.prompt}} - every line of code you generate must be valid code. Do not include code fences in your response, for example

Bad response (because it contains the code fence):
```javascript
console.log("hello world")
```

Good response (because it only contains the code):
console.log("hello world")

Begin generating the code now.