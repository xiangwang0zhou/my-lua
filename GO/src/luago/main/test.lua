local t={"a","B","c"}
t[2]="b"
t["foo"]="bar"
local s=t[3]..t[2]..t[1]..t["foo"]..#t
