## Xlang

It is a toy interpreter made while learning how to make one.

For the moment it supports:

- Variable declaration
- Arrays
- Strings
- Integers
- Functions
- Passing functions as parameters
- Helper methods like len(), push(), pop(), shift(), unshift(), reduce...
- HashMaps

## What's coming

- More examples in the example folder :)
- API Docs

# Where can I use it?

I deployed an online editor where you can execute code! The code is executed on a server, not in the JavaScript client, so you can use as well the API if you want to create your own online editor.

Web editor: https://boiling-basin-02644.herokuapp.com/
API URL: PUT https://boiling-basin-02644.herokuapp.com/api/v2

Sample body:

```
{
   "code":"// Variables usually are immutable in Xlang (except for hashmaps)\n// Example: If you declare the variable x, you cant do x = 10 later on, but you can redeclare like let x = 10\n// When you push or pop an array, you wont change directly the array.\n// You will get the new array (with the changes made) but the original wasn't mutated\n\nlet arr = [1, 2, 3];\n\n// This is already done in the standard methods of Xlang, but it's just for showing you how its implemented!\nlet map = fn(x, f) {\n  let iter = fn (arr, result) {\n    if (len(arr) == 0) {\n      return result;\n    }\n    iter(shift(arr), push(result, f(first(arr))));\n  }\n  // In Xlang there are implicit returns! This is the same as saying: return iter(x, []);\n  iter(x, []);\n}\n\n// [2, 3, 4]\nlet what = map(arr, fn (element) { element + 1 });\n// At the end of the file if you want logs call log()...\nlog(what)\n\n// [2, 3, 4, 222]\nlet what = push(what, 222);\nlog(what);\n// [2, 3, 4]\nlet what = pop(what);\nlog(what);\n// [1, 2, 3, 4]\nlet what = unshift(what, 1)\nlog(what)\n// [1000, 2, 3, 4]\nlet what = set(what, 0, 10000)\nlog(what)\n\n// Another function implemented already, but just showing you how it works!\nlet filter = fn(x, condition) {\n  let iter = fn(x, result) {\n    if (len(x) == 0) {\n      return result\n    }\n    // If the condition is met, we pass to the next iteration the updated array, but if not, we keep the same array!\n    iter(shift(x), if (condition(first(x))) { push(result, first(x)) } else { result })\n  }\n  iter(x, [])\n}\n\nlet what = set(what, 0, 10000);\nlet compar = true == \"false\";\n\nlog(what);\n\n// [2, 3, 4]\nlet filtered = filter(what, fn(n) { n < 1000 })\nlog(filtered)"
}
```

Sample response:

```
{
   "data":{
      "parse_error":{
         "line":0,
         "messages":null
      },
      "error":{
         "line":15,
         "messages":[
            "Error: Type mismatch: BOOL == STRING"
         ]
      },
      "output":[
         {
            "line":30,
            "messages":[
               "[2,3,4]"
            ]
         },
         {
            "line":33,
            "messages":[
               "[2,3,4,222]"
            ]
         },
         {
            "line":35,
            "messages":[
               "[2,3,4]"
            ]
         },
         {
            "line":37,
            "messages":[
               "[1,2,3,4]"
            ]
         },
         {
            "line":39,
            "messages":[
               "[10000,2,3,4]"
            ]
         }
      ]
   },
   "status":200
}
```

## How does it work?

Check the folder examples for integer, string, array examples...
Important points:

- Variables are immutable, check examples on how you can play around them :).
- At the moment it will only log the logs that you've made if you include log() at the end of the file, YEAH, it sucks, but I'll fix it soon.
- Have fun with it!

THIS IS ONLY FOR LEARNING PURPOSES. It isn't supposed to go in production or anything like that.
