# ClickHouse-log-analyzer

## About
Fast and simple ClickHouse log analyzer

## Installation
##### Prerequisites
* `Go`
* `git`
##### Clone the repo
```shell
git clone https://github.com/ChistaDATA/ChistaDATA-Profiler-for-ClickHouse.git
```
##### Get all dependencies
```shell
go get ./...
```
##### Build
```shell
go build -o profiler cmd/main.go
```

## Example
##### Run with a sample file
```shell
./profiler -n 2 ./server.log
```
##### Output
```

# Current date: 2022-10-06 13:06:39.088724 +0530 IST m=+19.943042793
# Hostname: Ahameds-MacBook-Pro.local
# Files:
	* /Users/chistadata/sandbox/ClickBench/clickhouse-local/logs/server.log
# Overall: 243122, Unique: 17 , QPS: 9.10
# Time range: 2022-10-03 18:42:08.317458 +0000 UTC to 2022-10-04 02:07:25.907276 +0000 UTC
# Attribute          total     min     max     avg     95%  stddev  median
# ============     ======= ======= ======= ======= ======= ======= =======
# Exec time        13676.0   0.00s   0.93s   0.06s   0.72s   0.18s   0.00s
# Rows read         53.44B    0.00   0.53M   0.22M   0.52M   0.16M   0.26M
# Bytes read        1.43TB   0.00B 26.02MB  6.16MB 24.91MB  7.74MB  2.27MB
# Peak Memory                0.00B  1.67GB  0.12GB  1.55GB  0.39GB  4.00MB

# Profile
# Rank Response time   Calls R/Call Query
# ==== =============== ===== ====== =====
#    1 12340.9  90.24% 16995  0.73s select distinct(city) from salary s1 LEFT JOIN salary s2 ON s1.event_id = s2.eve
#    2  43.53s   0.32% 17070  0.00s select min(amount) from salary where rand > 10000


# Query 1 : 1.243 QPS
# Time range: From 2022-10-03 18:42:08.4814 +0000 UTC to 2022-10-04 02:05:44.907986 +0000 UTC
# ====================================================================
# Attribute      total     min     max     avg     95%  stddev  median
# ============ ======= ======= ======= ======= ======= ======= =======
# Count         17.00K 
# Exec time    12340.9   0.00s   0.93s   0.73s   0.75s   0.04s   0.73s
# Rows read      8.77B    0.00   0.53M   0.52M   0.53M  21.50K   0.52M
# Bytes read    0.35TB   0.00B 22.03MB 21.33MB 21.97MB  0.89MB 21.37MB
# Peak Memory            0.00B  1.67GB  1.54GB  1.55GB 64.00MB  1.55GB
# ====================================================================
# Databases    default (16995/16995)  
# Hosts        [::1] (16995/16995)  
# Users        default (16995/16995)  
# Completion   16988/16995
# Query_time distribution
# ====================================================================
#   1us  
#  10us  
# 100us  
#   1ms  
#  10ms  
# 100ms  ###########################################################
#    1s  
#  10s+  
# ====================================================================
# Query
select distinct(city) from salary s1 LEFT JOIN salary s2 ON s1.event_id = s2.event_id 

# Query 2 : 1.248 QPS
# Time range: From 2022-10-03 18:42:11.876586 +0000 UTC to 2022-10-04 02:07:25.693396 +0000 UTC
# ====================================================================
# Attribute      total     min     max     avg     95%  stddev  median
# ============ ======= ======= ======= ======= ======= ======= =======
# Count         17.07K 
# Exec time     43.53s   0.00s   0.13s   0.00s   0.00s   0.00s   0.00s
# Rows read      4.41B    0.00   0.27M   0.26M   0.27M  10.29K   0.26M
# Bytes read   65.66GB   0.00B  4.07MB  3.94MB  4.05MB  0.16MB  3.95MB
# Peak Memory           0.53MB  4.20MB  2.88MB  2.89MB 90.90KB  2.88MB
# ====================================================================
# Databases    default (17070/17070)  
# Hosts        [::1] (17070/17070)  
# Users        default (17070/17070)  
# Completion   17066/17070
# Query_time distribution
# ====================================================================
#   1us  
#  10us  
# 100us  
#   1ms  ###########################################################
#  10ms  
# 100ms  
#    1s  
#  10s+  
# ====================================================================
# Query
select min(amount) from salary where rand > 10000

```

## Usage
```
profiler <file-paths> ...
```

Arguments:
```
<file-paths> ...    Paths of log files
```

Flags:
```
-h, --help                     Show context-sensitive help.
-n, --top-query-count=10       Count of queries for top x table
-r, --report-type="text"       Report type to be generated, types: md, text
-c, --minimum-query-count=1    Minimum no of query calls needed
```