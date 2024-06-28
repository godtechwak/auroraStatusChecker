How to use
===============
```
$ go mod init aurora_status_check
$ go mod tidy
$ go build ./main.go
$ ./main <api 호출 제한을 위한 sleep 밀리초>
ex) ./main 2000
```

Examples
===============
#### (1) 클러스터 or 인스턴스 목록 파일 생성
```
$ ./main 2000
================================
Aurora Version & Parameter Check
================================
Region(kr/jp/ca/uk):
kr
WorkType(cluster/instance):
instance
Do you want to create the list file? (yes/no):
yes
List saved to file: ./db_instance_list_kr.txt
```

#### (2) 리전 선택 (파일 생성 이후 재실행)
```
$ ./main 2000
================================
Aurora Version & Parameter Check
================================
Region(kr/jp/ca/uk):
kr
```
#### (3) 클러스터 또는 인스턴스 선택
```
WorkType(cluster/instance):
instance
```

Result
===============
##### cluster
```
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
                             Time│                          Duration|                           Cluster│                           Version│                            Status│                      Param Status│
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
                     02:50:20.457│                            15.08s|                  prod-silver-main│           8.0.mysql_aurora.3.04.2│                         available│                    pending-reboot│
                     02:50:20.594│                           15.021s|               prod-silver-grafana│           8.0.mysql_aurora.3.04.2│                         available│                    pending-reboot│
                     02:50:20.717│                           15.056s|                  prod-silver-blog│           8.0.mysql_aurora.3.04.2│                         available│                           in-sync│
                     02:50:20.822│                           15.033s|                  prod-silver-gold│           8.0.mysql_aurora.3.04.1│                         available│                           in-sync│
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
```
##### instance
```
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
                             Time│                          Duration│                          Instance│                           Version│                            Status│                      Param Status│
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
                     14:23:08.548│                            4.626s│                prod-silver-main-1│           8.0.mysql_aurora.3.04.2│                         available│                           in-sync│
                     14:23:08.882│                            4.729s│                prod-silver-main-2│           8.0.mysql_aurora.3.04.2│                         available│                           in-sync│
                     14:23:09.108│                            4.698s│             prod-silver-grafana-1│           8.0.mysql_aurora.3.04.2│                         available│                           in-sync│
                     14:23:09.337│                            4.651s│             prod-silver-grafana-2│           8.0.mysql_aurora.3.04.2│                         available│                           in-sync│
                     14:23:09.551│                             4.62s│                prod-silver-gold-0│           8.0.mysql_aurora.3.04.2│                         available│                           in-sync│
                     14:23:09.675│                            4.502s│                prod-silver-gold-1│           8.0.mysql_aurora.3.04.2│                         available│                           in-sync│
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
```

How to Add a Query Field
===============
- rds_util/cluster.go 또는 rds_util/instance.go 파일에서 `Describe()`와 `GetHeaders()` 함수 수정
#### DB 인스턴스의 인스턴스 클래스를 추가 조회하고 싶은 경우
(1) Describe()
```
values := []*string{
		instanceInfo.DBInstanceIdentifier,
		instanceInfo.EngineVersion,
		instanceInfo.DBInstanceStatus,
		instanceParam.ParameterApplyStatus,
		instanceInfo.DBInstanceClass, // 추가된 DB 인스턴스 클래스
	}
```
(2) GetHeaders()
```
return []string{"Time", "Duration", "Instance Name", "Version", "Status", "Param Status", "DB Instance Class"}
```
(3) 빌드 후 실행
```
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
                             Time│                          Duration│                     Instance Name│                           Version│                            Status│                      Param Status│                 DB Instance Class│
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
                     15:49:47.736│                             2.28s│                prod-silver-gold-0│           8.0.mysql_aurora.3.04.1│                         available│                           in-sync│                     db.t4g.medium│
                     15:49:47.738│                             2.28s│                prod-silver-gold-1│           8.0.mysql_aurora.3.04.1│                         available│                           in-sync│                     db.t4g.medium│
                     15:49:47.738│                             2.28s│                   prod-silver-a-3│           8.0.mysql_aurora.3.04.0│                         available│                    pending-reboot│                     db.r5.2xlarge│
                     15:49:47.739│                             2.28s│                   prod-silver-a-4│           8.0.mysql_aurora.3.04.0│                         available│                    pending-reboot│                     db.r5.2xlarge│
                     15:49:47.739│                             2.28s│             prod-silver-message-0│           8.0.mysql_aurora.3.04.0│                         available│                    pending-reboot│                      db.t3.medium│
                     15:49:47.738│                             2.28s│             prod-silver-message-1│           8.0.mysql_aurora.3.04.0│                         available│                    pending-reboot│                      db.t3.medium│
                     15:49:47.739│                             2.28s│             prod-silver-account-4│           8.0.mysql_aurora.3.04.0│                         available│                           in-sync│                    db.r6g.2xlarge│
                     15:49:47.739│                             2.28s│             prod-silver-account-5│           8.0.mysql_aurora.3.04.0│                         available│                    pending-reboot│                    db.r6g.2xlarge│
                     15:49:47.739│                             2.28s│             prod-silver-account-6│           8.0.mysql_aurora.3.04.0│                         available│                           in-sync│                    db.r6g.2xlarge│
────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ────────────────────────────────── ──────────────────────────────────
```
