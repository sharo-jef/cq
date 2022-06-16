# cq

csv processor like jq

```
$ cat a.csv
a,b,c,d
o,p,q,r
w,x,y,z
$ cq a.csv -o y
- a: o
  b: p
  c: q
  d: r
- a: w
  b: x
  c: "y"
  d: z

$ cat a.csv | cq -o j -c
[{"a":"o","b":"p","c":"q","d":"r"},{"a":"w","b":"x","c":"y","d":"z"}]
```

Filter function is not yet implemented.
