# cq

csv processor

```
$ cat a.csv -o c
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
$ cat a.csv | cq -c
[{"a":"o","b":"p","c":"q","d":"r"},{"a":"w","b":"x","c":"y","d":"z"}]
```

Filter function is not yet implemented.
