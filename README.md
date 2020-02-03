# Experiment: Lambda computing performance vs RAM

Execute a Lambda, with specified compute/blocking mixes, with varying
Lambda RAM settings.

## Usage

Head to [`lambda`](lambda) and look at [`Makefile`](lambda/Makefile)
to get the test Lambda going.

Consider the compute/blocking operations distributions as specified in
[`spec.json`](spec.json).  Use a [gamma
distribution](https://en.wikipedia.org/wiki/Gamma_distribution) to
specify how long (in milliseconds) a simulated blocking operation
takes.  Use [Poisson
distributions](https://en.wikipedia.org/wiki/Poisson_distribution) to
specify how many non-blocking operations each execution has and how
many work units each of those operations performs.  A work unit is
specified by `StepKeys`, which is the number of generates keys in a
map, and `StepRounds`, which is the number of map
serialization/deserialization round trips performed in a work unit.

At 512MB, the example [spec](spec.json) does about 4x more blocking
than computing (at 512MB).

After editing `spec.json`, [run](run.sh) a test that executes the
Lambda multiple times for various RAM settings.


```Shell
(cd call-lambda && make && ./run.sh)
```

Examine the output (`d.csv`):

```R
library(tidyverse)
d <- read_csv("call-lambda/d.csv")
summary(d)

d %>% group_by(mb) %>% 
	summarize(mean_work_time=mean(work_time), mean_block_time=mean(block_time))

ggplot(d, aes(worked,ms,color=mb,group=mb)) + 
	geom_point(alpha=0.4) + 
	geom_smooth(method=lm) + 
	scale_colour_gradient(low="red",high="blue") +
	labs(title="Work vs Elapsed time by RAM",
	     subtitle="Blocking 4x computing at 512MB.")
ggsave("eff-by-ram.png")


d %>% filter(mb==512) %>% 
	mutate(percent_blocked=100*block_time/ms) %>% 
	ggplot(aes(percent_blocked)) + 
	geom_histogram(bin_width=5) + 
	labs(title="Percent of time in blocking ops at 512MB",x="Percent of total time") + 
	xlim(0,100)
ggsave("blocking.png")

# Efficiency
d$e <- d$ms/d$worked

m <- lm(e ~ mb, data=d)
d1 <- data.frame(mb=seq(128,1024,128))
predict(m, d1)
d1$e <- predict(m, d1)

# For average work, what's the expected latency?
d1$ms <- mean(d$worked) * d1$e

ggplot(d1, aes(mb, ms)) + 
	geom_point() + 
	geom_smooth() + 
	ylim(0, max(d1$ms)) + 
	labs(title="Predicted mean latency by RAM")
ggsave("predict.png")
```

![graph](eff-by-ram.png)

![blocking](blocking.png)

The predicted mean latencies (milliseconds) by RAM tier (for blocking
4x more than computing):

```
    mb        e        ms
1  128 7.405195 148.81850
2  256 6.655009 133.74240
3  384 5.904824 118.66629
4  512 5.154638 103.59019
5  640 4.404453  88.51409
6  768 3.654267  73.43798
7  896 2.904082  58.36188
8 1024 2.153896  43.28578
```

Of course, the level of expected latency, in the context of a specific
application, is an important consideration.

![predict](predict.png)

Let's look quickly at costs.  Lambda has a 100ms billing minimum.

```R
costPerSecGB <- 0.000016667
dc <- d %>% 
	mutate(secs = ceiling(ms/100)*100/1000, cost = secs * costPerSecGB * (mb/1024)) %>% 
	group_by(mb) %>% 
	summarize(avgSecs = mean(secs), 
		costPerMillion = 0.20 + sum(cost)*1000*1000/n(),
		meanLatencyMS = mean(ms),
		percentageBlocking = mean(ifelse(ms > 0, block_time/ms, 0)))
```

`milCost` is the cost per million executions.

```
     mb avgSecs costPerMillion meanLatencyMS percentageBlocking
  <dbl>   <dbl>          <dbl>         <dbl>              <dbl>
1   128   0.278          0.778         229.               0.407
2   256   0.168          0.901         118.               0.518
3   384   0.137          1.06           86.2              0.658
4   512   0.129          1.27           77.4              0.715
5   640   0.127          1.52           74.7              0.740
6   768   0.125          1.77           73.8              0.746
7   896   0.125          2.03           73.8              0.746
8  1024   0.126          2.30           73.7              0.747
```

For this example job specification, which spends about 75% of time in
blocking operations, the 128MB Lambdas are 1/3 the cost of using 1GB.
_Note that this job specification has low mean latency, near the 100ms
minimum billing unit._ The mean 128MB latency of 229ms is three times
higher than the 1GB latency.  If you can tolerate 229ms mean latency,
then go with 128MB.  You spend $0.78/million instead of (say)
$2.30/million at 1GB.  If you are sensitive to higher latencies, 384MB
is the sweet spot, with 137ms mean latency and only $1.06 per million.

```R
dc %>% ggplot(aes(costPerMillion,meanLatencyMS,color=mb,size=mb)) + 
	geom_point() + 
	scale_colour_gradient(low="red",high="blue") + 
	labs(title="Cost vs latency for example Lambda") + 
	ylim(0,max(dc$meanLatencyMS)) +
	xlim(0,max(dc$costPerMillion)) + 
	guides(size=FALSE)
ggsave("costvlat.png")
```

![costvlat](costvlat.png)

```R
dq <- ds %>% 
	mutate(secs = ceiling(ms/100)*100/1000, cost = secs * costPerSecGB * (mb/1024), blocking = block_time/ms) %>% 
	group_by(steps_lam,mb) %>% 
	summarize(blocking = mean(blocking), avg_secs = mean(secs), cost_per_million = 0.20 + sum(cost)*1000*1000/n())

dq %>% arrange(mb,blocking) %>% mutate(MB = as.factor(mb)) %>% ggplot(aes(avg_secs,cost_per_million,shape=MB,group=c(MB))) + geom_point(aes(size=blocking,color=blocking)) + geom_line(alpha=0.3) + scale_shape_manual(values=1:10)

dq %>% filter(avg_secs < 0.5, blocking > 0.75) %>% arrange(cost_per_million)
```
