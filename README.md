# squ (eezer)

One step to scale down your test cases based on your diff. For multiple languages.

## Before vs After

Before: I have to run all the cases for even a tiny change. Which takes 40 mins.

<img width="694" alt="image" src="https://user-images.githubusercontent.com/13421694/226383324-80c878de-815b-4a80-8666-772dc44fedf0.png">

```bash
$ pytest --collect-only
============================= test session starts ==============================
platform darwin -- Python 3.7.4, pytest-7.2.2, pluggy-0.12.0
rootdir: /Users/fengzhangchi/github_workspace/stagesepx
plugins: cov-2.7.1
collected 45 items

...

========================= 45 tests collected in 15.30s =========================
```

After: Only need to run the influenced cases.

<img width="705" alt="image" src="https://user-images.githubusercontent.com/13421694/226386547-d4c88519-56b3-4113-b8fe-6cda8e1778d0.png">

```bash
$ ./squ | xargs pytest --collect-only
============================= test session starts ==============================
platform darwin -- Python 3.7.4, pytest-7.2.2, pluggy-0.12.0
plugins: cov-2.7.1
collected 5 items

...

-- Docs: https://docs.pytest.org/en/stable/how-to/capture-warnings.html
========================= 5 tests collected in 14.93s ==========================
```

All you need is a `./squ | xargs` prefix.

## Usage

You can find [pre-build release binaries](https://github.com/opensibyl/squ/releases) here.

One step on your command:

| Language | Before     | After                                  |
|----------|------------|----------------------------------------|
| Python   | `pytest`   | <code>squ &#124; xargs pytest</code>   | 
| Java     | `mvn test` | <code>squ &#124; xargs mvn test</code> |
| Golang   | `go test`  | <code>squ &#124; xargs go test</code>  |

Much benefits:

| Repo                                                                | Language | Analyze time | Before | After | Optimize |
|---------------------------------------------------------------------|----------|--------------|--------|-------|----------|
| [stagesepx (our own case)](https://github.com/williamfzc/stagesepx) | Python   | ~1 s         | 45     | 0~15  | ~75%     |
| [Jacoco](https://github.com/jacoco/jacoco)                          | Java     | ~7 s         | 1364   | 0~983 | ~50%     |
| [sibyl2](https://github.com/opensibyl/sibyl2)                       | Golang   | ~1 s         | 53     | 0~20  | ~70%     |

## How it works?

This project inspired by Google.

![](https://user-images.githubusercontent.com/13421694/207057947-894c1fb9-8ce4-4f7b-b5d3-88d220003e82.png)

Google used static analyze tech to scale down their time cost in test phase. In their case, only the influenced part
will be evaluated and tested.

<img width="572" alt="image" src="https://user-images.githubusercontent.com/13421694/226381849-0e5217f6-b3cf-4bb4-9d1e-ca427cafc81c.png">

We implement it onto [sibyl2](https://github.com/opensibyl/sibyl2). See [Contribution](#contribution) for details.

## Usage Details

### How to debug?

`squ --debug` and you will see the log in `squ-debug.log`.

```text
{"level":"info","ts":1678975693.395471,"caller":"upload/upload.go:178","msg":"upload batch: 0 - 50"}
{"level":"info","ts":1678975693.3977358,"caller":"indexer/indexer_base.go:78","msg":"indexer done"}
{"level":"info","ts":1678975693.400231,"caller":"UnitSqueezer/api.go:103","msg":"indexer ready"}
{"level":"info","ts":1678975693.412002,"caller":"extractor/extractor_git.go:56","msg":"version.go [3] => functions 0"}
{"level":"info","ts":1678975693.412069,"caller":"UnitSqueezer/api.go:113","msg":"diff calc ready: 0"}
{"level":"info","ts":1678975693.412087,"caller":"UnitSqueezer/api.go:120","msg":"case analyzer done, before: 3, after: 0"}
{"level":"info","ts":1678975695.459771,"caller":"UnitSqueezer/api.go:137","msg":"runner scope"}
{"level":"info","ts":1678975695.459826,"caller":"UnitSqueezer/api.go:139","msg":"no cases need to run"}
{"level":"info","ts":1678975695.459852,"caller":"UnitSqueezer/api.go:145","msg":"prepare stage finished, total cost: 2155 ms"}
{"level":"info","ts":1678975695.459863,"caller":"UnitSqueezer/api.go:147","msg":"runner cmd: --run=\"^$\""}
{"level":"info","ts":1678975695.460073,"caller":"server/app.go:113","msg":"sibyl server down: http: Server closed"}
```

### Set another source dir?

`squ --src ../SOMEWHERE`

### Compare another version, not `HEAD~1`

`squ --before HEAD~5`

## Still have a question?

Please use our [issue board](https://github.com/opensibyl/squ/issues).

## Contribution

> PR is always welcome.

This project consists of 4 parts.

### index (preparation)

1. upload all the files to sibyl2
2. search and tag all the test methods
3. calc and tag all the test methods influencing scope

### extract (extract data from workspace)

1. calc diff between current and previous
2. find methods influenced by diff

### mapper (mapping between cases and diff)

1. search related cases

### runner (can be implemented by different languages)

1. build test commands for different frameworks
2. call cmd

## License

[Apache License 2.0](LICENSE)
