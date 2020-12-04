# DMP Requirements Checker

The application is used to generate a report for checking requirements of hardware, OS, Kernel, external resources, and etc.

ReqCheck allows cases with combinations of deployments and hosts.


## Case 3 Hosts

Suppose, we have 3 separated hosts for DMP (ip 10.17.0.2), Clickhouse (ip 10.17.0.3) and Enclave (ip 10.17.0.4).


### Prepare

These steps make a stub of all checking services to check accessibility between hosts.

Run on the DMP host:

`$ ./dmp-reqcheck listen --roles dmp`


Run on the Clickhouse host:

`$ ./dmp-reqcheck listen --roles ch`


Run on the Enclave host:

`$ ./dmp-reqcheck listen --roles enclave`


### Reports

You should run check command on each hosts to get report results of matching requirements.

Run on the Dmp host:

`$ ./dmp-reqcheck check --roles dmp --hosts ch:10.17.0.3,enclave:10.17.0.4`


Run on the Clickhouse host:

`$ ./dmp-reqcheck listen --roles ch`


Run on the Enclave host:

`$ ./dmp-reqcheck listen --roles enclave --hosts ch:10.17.0.3,dmp:10.17.0.2`



## Case 2 Hosts

Suppose we have 2 separated hosts. One (ip 10.17.0.2) for both DMP and Clickhouse, other for Enclave (ip 10.17.0.4) only.

ReqCheck will summarize hardware requirements and check them on the host with several roles.

### Prepare

Run on the Dmp&Clickhouse host:

`$ ./dmp-reqcheck listen --roles dmp,ch`

Run on the Enclave host:

`$ ./dmp-reqcheck listen --roles enclave`


### Reports

You should run check command on 2 hosts to get report results of matching requirements.

Run on the Dmp&Clickhouse host:

`$ ./dmp-reqcheck check --roles dmp,ch --hosts ch:10.17.0.2,enclave:10.17.0.4`

Run on the Enclave host:

`$ ./dmp-reqcheck listen --roles enclave --hosts ch:10.17.0.2,dmp:10.17.0.2`


## Screenshots

Example report for Clickhouse role.

![ScreenShot](https://github.com/aggregion/dmp-reqcheck/blob/master/.images/ch-report-example.png)
