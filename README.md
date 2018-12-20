### usage
```
./31fcn --config=config.xml
```

### version
```
v0.0.1
```

### config
```
<?xml version="1.0" encoding="UTF-8" ?>
<config>
	<f3cn>{user},{password}</f3cn>
	<task>
		<startup>true</startup>
		<schedule>0 0 * * * *</schedule>
	</task>
	<redis_server></redis_server>
</config>
```

### schedule
```
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
```
* https://godoc.org/github.com/robfig/cron

### faq
* 更多疑问请+qq群 233415606
