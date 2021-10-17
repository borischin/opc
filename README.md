# Open Policy Configuration

Use OPA(Open Policy Agent) to manage your configurations. You can custom your own policies(ex: env, market, region ...etc).

## Example

### Simple Policy Config

simple_policy_config.rego

```rego
package main

env(x) { input.env == x}

_qa = env("QA")
_prod = env("PROD")

key1 = "value1-qa" { _qa }
else = "value1-prod" { _prod }
else = null

key2 = 2
key3 = "值3"

```

command with input env=QA and output default json format

```bash
opc -m simple_policy_config.rego -i env=QA
```

output

```json
{
    "_qa": true,
    "key1": "value1-qa",
    "key2": 2,
    "key3": "值3"
}
```

command with input env=PROD and output env-file format

```bash
opc -m simple_policy_config.rego -i env=PROD -f env-file
```

output

```env
_prod=true
key1=value1-prod
key2=2
key3=值3
```

### Complicated Policy Config

complicated_policy_config.rego

```rego
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
```

command with input env=QA, market=TW

```bash
opc -m complicated_policy_config.rego -i env=QA -i market=TW
```

output

```json
{
    "_qa": true,
    "_tw": true,
    "key1": "value1-qa",
    "key2": 3,
    "key3": "qa-值3-tw"
}
```

command with input env=PROD, market=US

```bash
opc -m complicated_policy_config.rego -i env=PROD -i market=US
```

output

```json
{
    "_prod": true,
    "_us": true,
    "key1": "value1-us-prod",
    "key2": 2,
    "key3": "值3-us"
}
```
