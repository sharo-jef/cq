# cq

csv processor like jq

```json
$ cat a.csv
a,b,c,d
1,2,3,4
5,6,7,8
$ cq a.csv
[
  {
    "a": "1",
    "b": "2",
    "c": "3",
    "d": "4"
  },
  {
    "a": "5",
    "b": "6",
    "c": "7",
    "d": "8"
  }
]
$ cat a.csv | cq -c
[{"a":"1","b":"2","c":"3","d":"4"},{"a":"5","b":"6","c":"7","d":"8"}]
```

Filter function is not yet implemented.
