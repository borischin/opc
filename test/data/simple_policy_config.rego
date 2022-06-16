package main

env(x) { input.env == x}

_qa = env("QA")
_prod = env("PROD")

key1 = "value1-qa" { _qa }
else = "value1-prod" { _prod }
else = null

key2 = 2
key3 = "值3"

key4 = { "p1": "v1-qa", "p2": "v2-qa" } { _qa }
else = { "p1": "v1-prod", "p2": "v2-prod" } { _prod }
else = {}



key5 = [ {"p1": "v1-1-qa", "p2": "v1-2-qa"}, {"p1": "v2-1-qa", "p2": "v2-2-qa"} ] { _qa }
else = [ {"p1": "v1-1-prod", "p2": "v1-2-prod"}, {"p1": "v2-1-prod", "p2": "v2-2-prod"} ] { _prod }
else = []

