package main

/*
CREATE TABLE `ecbase_card` (
`id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '自增主键',
`company_id` BIGINT NOT NULL COMMENT '所属企业Id，作为分区字段，与自增长id作为联合主键',
`code` VARCHAR(32) DEFAULT null COMMENT '车企编码 kafak导入',
`platform_id` VARCHAR(64) DEFAULT null COMMENT '平台ID',
`subs_id` VARCHAR(32) DEFAULT null COMMENT '用户id kafak导入',
`uuid` VARCHAR(64) NOT NULL COMMENT '唯一标识',
`create_time` DATETIME NOT NULL COMMENT '创建时间',
`create_by` VARCHAR(255) DEFAULT null COMMENT '创建者',
`update_time` DATETIME DEFAULT null COMMENT '更新时间',
`update_by` VARCHAR(255) DEFAULT null COMMENT '更新者',
`iccid` VARCHAR(32) NOT NULL COMMENT 'ICCID',
`imsi` VARCHAR(32) DEFAULT null COMMENT 'IMSI kafak导入',
`msisdn` VARCHAR(32) DEFAULT null COMMENT 'IMSI kafak导入',
`carrier` INT NOT NULL COMMENT '所属运营商',
`imei` VARCHAR(32) DEFAULT null COMMENT 'IMEI',
`vin` VARCHAR(32) DEFAULT null COMMENT '车辆VIN码',
`open_time` DATETIME DEFAULT null COMMENT '开户日期',
`active_date` DATETIME DEFAULT null COMMENT '激活日期',
`network_type` VARCHAR(8) DEFAULT null COMMENT '网络类型（字典项编码）',
`card_type` VARCHAR(8) DEFAULT null COMMENT '卡片物理类型（字典项编码）',
`belong_place` VARCHAR(32) DEFAULT null COMMENT '归属地',
`remark` VARCHAR(255) DEFAULT null COMMENT '备注',
`card_status` VARCHAR(8) DEFAULT null COMMENT '卡状态字典项编码，根据平台区分 kafak导入',
`status_time` DATETIME DEFAULT null COMMENT '卡号状态变更时间',
`vehicle_status` INT DEFAULT null COMMENT '车辆状态',
`vehicle_out_factory_time` DATETIME DEFAULT null COMMENT '车辆出厂时间',
`realname_status` INT DEFAULT '0' COMMENT '实名状态：0-初始化未实名，1-受理成功处理中，2-实名登记成功，3-实名注销成功',
`realname_status_t1` INT DEFAULT null COMMENT 't1实名状态：1：已实名，2：未实名',
`pay_account_id` VARCHAR(32) DEFAULT null COMMENT '付费账户id',
`pay_account_name` VARCHAR(64) DEFAULT null COMMENT '付费账户名称',
`plat_type` VARCHAR(8) DEFAULT null COMMENT '平台类型2:pb 3:ct',
`cust_id` VARCHAR(32) DEFAULT null COMMENT '客户id kafak导入',
`cust_name` VARCHAR(32) DEFAULT null COMMENT '客户名称',
`cust_type` VARCHAR(32) DEFAULT null COMMENT '客户类型 一般为C客户类型',
`be_id` VARCHAR(32) DEFAULT null COMMENT '省份编码',
`region_id` VARCHAR(32) DEFAULT null COMMENT '归属地编码',
`group_id` VARCHAR(128) DEFAULT null COMMENT '归属群组id',
`group_member_status` VARCHAR(128) DEFAULT null COMMENT '归属群组中成员状态',
`group_pool_data_usage` BIGINT DEFAULT null COMMENT '归属群组中本月池内用量',
`cust_code` VARCHAR(255) DEFAULT null COMMENT '客户集团编码',
`sync_time` DATETIME DEFAULT null COMMENT '通过kafka入库时，每次必须更新的字段，其他入口不用变动',
`boss` INT NOT NULL COMMENT '运营商BOSS系统，1-移动CT，2-移动PB，3-电信DCP，4-电信M2M，5-联通jasper，6-联通CMP，7-中国电信5GCMP平台',
`account_id` BIGINT DEFAULT null COMMENT 'BOSS系统账号id',
`vc_code` VARCHAR(32) DEFAULT null COMMENT '车企vc_code',
`source_create_time` DATETIME DEFAULT null COMMENT '源端数据创建时间',
`source_modify_time` DATETIME DEFAULT null COMMENT '源端数据修改时间',
`device_num` VARCHAR(32) DEFAULT null COMMENT '设备号',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_ecbase_card_uuid` (`uuid`),
UNIQUE KEY `uk_ecbase_card_vc_code_iccid` (`vc_code`,`iccid`) COMMENT '20231021新增',
KEY `ecbase_card_ecbase_company_id_fk` (`id`,`company_id`),
KEY `ecbase_card_ecbase_id_company_id` (`company_id`),
KEY `idx_ecbase_card_create_time` (`create_time`),
KEY `idx_account_id` (`account_id`),
KEY `ecbase_card_subs_id` (`subs_id`) COMMENT '20230529新增',
KEY `idx_ecbase_card_msisdn` (`msisdn`),
KEY `ecbase_card_vc_code_carrier_card_status_idx` (`vc_code`,`carrier`,`card_status`),
KEY `ecbase_card_vc_code_vehicle_status_idx` (`vc_code`,`vehicle_status`),
KEY `ecbase_card_vc_code_realname_status_idx` (`vc_code`,`realname_status`),
KEY `ecbase_card_vin_idx` (`vin`),
KEY `ecbase_card_iccid` (`iccid`) COMMENT '2310.7新增'
) COMMENT='卡信息表';
*/

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/rand"
)

var (
	dbUser = "root"
	dbPass = "111"
	dbName = "t"
	db     *sql.DB
	err    error
	stmt   *sql.Stmt

	usePrepare  = true
	concurrency = 256
	// insertCnt   = 2000000
	insertCnt   = 7813
)

func init() {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:6001)/%s", dbUser, dbPass, dbName))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行 drop table 语句
	_, err = db.Exec("drop table if exists ecbase_card")
	if err != nil {
		fmt.Println("执行 drop table 错误:", err)
		return
	}

	// 执行创建表语句
	_, err = db.Exec(`CREATE TABLE ecbase_card (
			id BIGINT NOT NULL AUTO_INCREMENT COMMENT '自增主键',
			company_id BIGINT NOT NULL COMMENT '所属企业Id，作为分区字段，与自增长id作为联合主键',
			code VARCHAR(32) DEFAULT null COMMENT '车企编码 kafak导入',
			platform_id VARCHAR(64) DEFAULT null COMMENT '平台ID',
			subs_id VARCHAR(32) DEFAULT null COMMENT '用户id kafak导入',
			uuid VARCHAR(64) NOT NULL COMMENT '唯一标识',
			create_time DATETIME NOT NULL COMMENT '创建时间',
			create_by VARCHAR(255) DEFAULT null COMMENT '创建者',
			update_time DATETIME DEFAULT null COMMENT '更新时间',
			update_by VARCHAR(255) DEFAULT null COMMENT '更新者',
			iccid VARCHAR(32) NOT NULL COMMENT 'ICCID',
			imsi VARCHAR(32) DEFAULT null COMMENT 'IMSI kafak导入',
			msisdn VARCHAR(32) DEFAULT null COMMENT 'IMSI kafak导入',
			carrier INT NOT NULL COMMENT '所属运营商',
			imei VARCHAR(32) DEFAULT null COMMENT 'IMEI',
			vin VARCHAR(32) DEFAULT null COMMENT '车辆VIN码',
			open_time DATETIME DEFAULT null COMMENT '开户日期',
			active_date DATETIME DEFAULT null COMMENT '激活日期',
			network_type VARCHAR(8) DEFAULT null COMMENT '网络类型（字典项编码）',
			card_type VARCHAR(8) DEFAULT null COMMENT '卡片物理类型（字典项编码）',
			belong_place VARCHAR(32) DEFAULT null COMMENT '归属地',
			remark VARCHAR(255) DEFAULT null COMMENT '备注',
			card_status VARCHAR(8) DEFAULT null COMMENT '卡状态字典项编码，根据平台区分 kafak导入',
			status_time DATETIME DEFAULT null COMMENT '卡号状态变更时间',
			vehicle_status INT DEFAULT null COMMENT '车辆状态',
			vehicle_out_factory_time DATETIME DEFAULT null COMMENT '车辆出厂时间',
			realname_status INT DEFAULT '0' COMMENT '实名状态：0-初始化未实名，1-受理成功处理中，2-实名登记成功，3-实名注销成功',
			realname_status_t1 INT DEFAULT null COMMENT 't1实名状态：1：已实名，2：未实名',
			pay_account_id VARCHAR(32) DEFAULT null COMMENT '付费账户id',
			pay_account_name VARCHAR(64) DEFAULT null COMMENT '付费账户名称',
			plat_type VARCHAR(8) DEFAULT null COMMENT '平台类型2:pb 3:ct',
			cust_id VARCHAR(32) DEFAULT null COMMENT '客户id kafak导入',
			cust_name VARCHAR(32) DEFAULT null COMMENT '客户名称',
			cust_type VARCHAR(32) DEFAULT null COMMENT '客户类型 一般为C客户类型',
			be_id VARCHAR(32) DEFAULT null COMMENT '省份编码',
			region_id VARCHAR(32) DEFAULT null COMMENT '归属地编码',
			group_id VARCHAR(128) DEFAULT null COMMENT '归属群组id',
			group_member_status VARCHAR(128) DEFAULT null COMMENT '归属群组中成员状态',
			group_pool_data_usage BIGINT DEFAULT null COMMENT '归属群组中本月池内用量',
			cust_code VARCHAR(255) DEFAULT null COMMENT '客户集团编码',
			sync_time DATETIME DEFAULT null COMMENT '通过kafka入库时，每次必须更新的字段，其他入口不用变动',
			boss INT NOT NULL COMMENT '运营商BOSS系统，1-移动CT，2-移动PB，3-电信DCP，4-电信M2M，5-联通jasper，6-联通CMP，7-中国电信5GCMP平台',
			account_id BIGINT DEFAULT null COMMENT 'BOSS系统账号id',
			vc_code VARCHAR(32) DEFAULT null COMMENT '车企vc_code',
			source_create_time DATETIME DEFAULT null COMMENT '源端数据创建时间',
			source_modify_time DATETIME DEFAULT null COMMENT '源端数据修改时间',
			device_num VARCHAR(32) DEFAULT null COMMENT '设备号',
			PRIMARY KEY (id),
			UNIQUE KEY uk_ecbase_card_uuid (uuid),
			UNIQUE KEY uk_ecbase_card_vc_code_iccid (vc_code, iccid) COMMENT '20231021新增',
			KEY ecbase_card_ecbase_company_id_fk (id, company_id),
			KEY ecbase_card_ecbase_id_company_id (company_id),
			KEY idx_ecbase_card_create_time (create_time),
			KEY idx_account_id (account_id),
			KEY ecbase_card_subs_id (subs_id) COMMENT '20230529新增',
			KEY idx_ecbase_card_msisdn (msisdn),
			KEY ecbase_card_vc_code_carrier_card_status_idx (vc_code, carrier, card_status),
			KEY ecbase_card_vc_code_vehicle_status_idx (vc_code, vehicle_status),
			KEY ecbase_card_vc_code_realname_status_idx (vc_code, realname_status),
			KEY ecbase_card_vin_idx (vin),
			KEY ecbase_card_iccid (iccid) COMMENT '2310.7新增'
		) COMMENT='卡信息表';`)
	if err != nil {
		fmt.Println("执行创建表错误:", err)
		return
	}

	if usePrepare {
		sql := `INSERT INTO ecbase_card(
			company_id, code, platform_id, subs_id, uuid, create_time, 
			create_by, update_time, update_by, iccid, imsi, msisdn, 
			carrier, imei, vin, open_time, active_date, network_type, 
			card_type, belong_place, remark, card_status, status_time, 
			vehicle_status, vehicle_out_factory_time, realname_status, realname_status_t1, pay_account_id, 
			pay_account_name, plat_type, cust_id, cust_name, cust_type, 
			be_id, region_id, group_id, group_member_status, group_pool_data_usage, 
			cust_code, sync_time, boss, account_id, vc_code, source_create_time, 
			source_modify_time, device_num) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		stmt, err = db.Prepare(sql)
		if err != nil {
			fmt.Println("Prepare error:", err)
			os.Exit(1)
		}
	}
}

type EcbaseCard struct {
	ID                    int
	CompanyID             int
	Code                  string
	PlatformID            string
	SubsID                string
	UUID                  string
	CreateTime            string
	CreateBy              string
	UpdateTime            string
	UpdateBy              string
	ICCID                 string
	IMSI                  string
	MSISDN                string
	Carrier               int
	IMEI                  string
	VIN                   string
	OpenTime              string
	ActiveDate            string
	NetworkType           string
	CardType              string
	BelongPlace           string
	Remark                string
	CardStatus            string
	StatusTime            string
	VehicleStatus         int
	VehicleOutFactoryTime string
	RealnameStatus        int
	RealnameStatusT1      int
	PayAccountID          string
	PayAccountName        string
	PlatType              string
	CustID                string
	CustName              string
	CustType              string
	BeID                  string
	RegionID              string
	GroupID               string
	GroupMemberStatus     string
	GroupPoolDataUsage    int64
	CustCode              string
	SyncTime              string
	Boss                  int
	AccountID             int64
	VcCode                string
	SourceCreateTime      string
	SourceModifyTime      string
	DeviceNum             string
}

func NewEcbaseCard() *EcbaseCard {
	card := &EcbaseCard{
		// ID:                    time.Now().UnixNano(),
		CompanyID:             int(rand.Int31()),
		Code:                  strconv.FormatInt(time.Now().UnixNano(), 10),
		PlatformID:            strconv.FormatInt(time.Now().UnixNano(), 10),
		SubsID:                strconv.FormatInt(time.Now().UnixNano(), 10),
		UUID:                  strconv.FormatInt(time.Now().UnixNano(), 10),
		CreateTime:            time.Now().Format("2006-01-02 15:04:05"),
		CreateBy:              strconv.FormatInt(time.Now().UnixNano(), 10),
		UpdateTime:            time.Now().Format("2006-01-02 15:04:05"),
		UpdateBy:              strconv.FormatInt(time.Now().UnixNano(), 10),
		ICCID:                 strconv.FormatInt(time.Now().UnixNano(), 10),
		IMSI:                  strconv.FormatInt(time.Now().UnixNano(), 10),
		MSISDN:                strconv.FormatInt(time.Now().UnixNano(), 10),
		Carrier:               int(rand.Int31()),
		IMEI:                  strconv.FormatInt(time.Now().UnixNano(), 10),
		VIN:                   strconv.FormatInt(time.Now().UnixNano(), 10),
		OpenTime:              time.Now().Format("2006-01-02 15:04:05"),
		ActiveDate:            time.Now().Format("2006-01-02 15:04:05"),
		NetworkType:           "notype",
		CardType:              "notype",
		BelongPlace:           strconv.FormatInt(time.Now().UnixNano(), 10),
		Remark:                strconv.FormatInt(time.Now().UnixNano(), 10),
		CardStatus:            strconv.FormatInt(time.Now().UnixNano(), 10),
		StatusTime:            time.Now().Format("2006-01-02 15:04:05"),
		VehicleStatus:         int(rand.Int31()),
		VehicleOutFactoryTime: time.Now().Format("2006-01-02 15:04:05"),
		RealnameStatus:        int(rand.Int31()),
		RealnameStatusT1:      int(rand.Int31()),
		PayAccountID:          strconv.FormatInt(time.Now().UnixNano(), 10),
		PayAccountName:        strconv.FormatInt(time.Now().UnixNano(), 10),
		PlatType:              "notype",
		CustID:                strconv.Itoa(rand.Intn(16)),
		CustName:              strconv.FormatInt(time.Now().UnixNano(), 10),
		CustType:              strconv.FormatInt(time.Now().UnixNano(), 10),
		BeID:                  strconv.FormatInt(time.Now().UnixNano(), 10),
		RegionID:              strconv.FormatInt(time.Now().UnixNano(), 10),
		GroupID:               strconv.FormatInt(time.Now().UnixNano(), 10),
		GroupMemberStatus:     strconv.FormatInt(time.Now().UnixNano(), 10),
		GroupPoolDataUsage:    rand.Int63(),
		CustCode:              strconv.FormatInt(time.Now().UnixNano(), 10),
		SyncTime:              time.Now().Format("2006-01-02 15:04:05"),
		Boss:                  int(rand.Int31()),
		AccountID:             time.Now().UnixNano(),
		VcCode:                strconv.FormatInt(time.Now().UnixNano(), 10),
		SourceCreateTime:      time.Now().Format("2006-01-02 15:04:05"),
		SourceModifyTime:      time.Now().Format("2006-01-02 15:04:05"),
		DeviceNum:             strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	return card
}

func insertData(db *sql.DB, usePrepare bool, concurrency int, insertCnt int) {
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var totalTime time.Duration
			for j := 0; j < insertCnt; j++ {
				card := NewEcbaseCard()
				start := time.Now()

				if usePrepare {
					_, err := stmt.Exec(card.CompanyID, card.Code, card.PlatformID, card.SubsID, card.UUID, card.CreateTime, card.CreateBy, card.UpdateTime, card.UpdateBy, card.ICCID, card.IMSI, card.MSISDN, card.Carrier, card.IMEI, card.VIN, card.OpenTime, card.ActiveDate, card.NetworkType, card.CardType, card.BelongPlace, card.Remark, card.CardStatus, card.StatusTime, card.VehicleStatus, card.VehicleOutFactoryTime, card.RealnameStatus, card.RealnameStatusT1, card.PayAccountID, card.PayAccountName, card.PlatType, card.CustID, card.CustName, card.CustType, card.BeID, card.RegionID, card.GroupID, card.GroupMemberStatus, card.GroupPoolDataUsage, card.CustCode, card.SyncTime, card.Boss, card.AccountID, card.VcCode, card.SourceCreateTime, card.SourceModifyTime, card.DeviceNum)
					if err != nil {
						log.Printf("Insert error: %v", err)
					}
				} else {
					_, err := db.Exec("INSERT INTO ecbase_card(company_id, code, platform_id, subs_id, uuid, create_time, create_by, update_time, update_by, iccid, imsi, msisdn, carrier, imei, vin, open_time, active_date, network_type, card_type, belong_place, remark, card_status, status_time, vehicle_status, vehicle_out_factory_time, realname_status, realname_status_t1, pay_account_id, pay_account_name, plat_type, cust_id, cust_name, cust_type, be_id, region_id, group_id, group_member_status, group_pool_data_usage, cust_code, sync_time, boss, account_id, vc_code, source_create_time, source_modify_time, device_num) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
						card.CompanyID, card.Code, card.PlatformID, card.SubsID, card.UUID,
						card.CreateTime, card.CreateBy, card.UpdateTime, card.UpdateBy, card.ICCID,
						card.IMSI, card.MSISDN, card.Carrier, card.IMEI, card.VIN, card.OpenTime, card.ActiveDate,
						card.NetworkType, card.CardType, card.BelongPlace, card.Remark, card.CardStatus, card.StatusTime,
						card.VehicleStatus, card.VehicleOutFactoryTime, card.RealnameStatus, card.RealnameStatusT1,
						card.PayAccountID, card.PayAccountName, card.PlatType, card.CustID, card.CustName, card.CustType,
						card.BeID, card.RegionID, card.GroupID, card.GroupMemberStatus, card.GroupPoolDataUsage, card.CustCode,
						card.SyncTime, card.Boss, card.AccountID, card.VcCode, card.SourceCreateTime, card.SourceModifyTime, card.DeviceNum)
					if err != nil {
						log.Printf("Insert error: %v", err)
					}
				}
				elapsed := time.Since(start)
				totalTime += elapsed
			}
			fmt.Printf("并发数 %d, 每个并发插入数据总量 %d, 是否使用Prepare %v, 耗时 %.2f 秒\n", concurrency, insertCnt, usePrepare, totalTime.Seconds())
		}()
	}
	wg.Wait()
}

func main() {
	insertData(db, usePrepare, concurrency, insertCnt)

	// // 测试使用预处理，高并发
	// fmt.Println("Testing with prepare, high concurrency...")
	// start = time.Now()
	// insertData(db, true, 10)
	// elapsed = time.Since(start)
	// fmt.Printf("Time taken: %v\n", elapsed)
}
