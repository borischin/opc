package main

env(x) { input.env == x}

_qa = env("QA")
_prod = env("PROD")

key1 = "value1-qa" { _qa }
else = "value1-prod" { _prod }
else = null

key2 = 2
key3 = "值3"
