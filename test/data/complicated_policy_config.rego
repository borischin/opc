package main

env(x) { input.env == x}
market(x) { input.market == x}

_qa = env("QA")
_prod = env("PROD")

_us = market("US")
_tw = market("TW")

# simple use
key1 = "value1-qa" { _qa }              # qa cross market
else = "value1-us-prod" { _prod; _us }  # prod depends on market
else = "value1-tw-prod" { _prod; _tw }  # prod depends on market
else = null

key2 = 2 { _us }  # only depends on market
else = 3 { _tw }  # only depends on market
else = null

# advanced use (when you have better naming convension)
key3 = sprintf("%v-值3-%v", [lower(input.env), lower(input.market)]) { _qa }
else = sprintf("值3-%v", [lower(input.market)]) { _prod }
else = null
