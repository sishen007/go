package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

var (
	kv                 clientv3.KV
	putResp            *clientv3.PutResponse
	config             clientv3.Config
	client             *clientv3.Client
	err                error
	getResp            *clientv3.GetResponse
	delResp            *clientv3.DeleteResponse
	kvpair             *mvccpb.KeyValue
	lease              clientv3.Lease
	leaseGrantResp     *clientv3.LeaseGrantResponse
	leaseId            clientv3.LeaseID
	keepRespChan       <-chan *clientv3.LeaseKeepAliveResponse
	keepResp           *clientv3.LeaseKeepAliveResponse
	watchStartRevision int64
	watcher            clientv3.Watcher
	watchRespChan      <-chan clientv3.WatchResponse
	watchResp          clientv3.WatchResponse
	event              *clientv3.Event
	putOp              clientv3.Op
	opResp             clientv3.OpResponse
	getOp              clientv3.Op
	txn                clientv3.Txn
	ctx                context.Context
	cancelFunc         context.CancelFunc
	txnResp            *clientv3.TxnResponse
)

func init() {
	// config
	config = clientv3.Config{
		Endpoints:   []string{"10.70.93.216:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}
	// 创建链接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	defer func() {
		client.Close()
		fmt.Println("Main 函数终结.")
	}()

	// 写数据
	// etcdPut()
	// 读数据
	// etcdGet()
	// 按目录获取
	// etcdGetByDir()
	// 删除键值对
	// etcdDel()
	// 租约
	// etcdLease()
	// 监听
	//etcdWatch()
	// Op操作
	//etcdOp()
	// 分布式锁
	//etcdTxn()

}
func etcdTxn() {
	// lease 实现锁自动过期
	// op操作
	// txn事务,if else then

	// 1. 上锁(创建租约,自动续租,拿着租约去抢占一个key)
	// 租约实例
	lease = clientv3.NewLease(client)
	// 申请一个10s的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}
	// 获取租约ID
	leaseId = leaseGrantResp.ID
	// 租约续租
	ctx, cancelFunc = context.WithCancel(context.TODO())
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}
	// 处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else {
					fmt.Println("收到自动续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()
	// 抢key,if 不存在key ,then 设置它,else 抢锁失败
	kv = clientv3.NewKV(client)
	// 创建事务
	txn = kv.Txn(context.TODO())
	// 定义事务
	// 如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/jobs/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/jobs/job9", "xxxxxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/jobs/job9")) // 否则抢锁失败

	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}
	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 2. 处理业务
	fmt.Println("处理任务")

	// 在锁内,很安全

	time.Sleep(5 * time.Second)

	// 3. 释放锁(取消自动续租,释放租约) defer释放
}
func etcdOp() {
	// 创建put op
	putOp = clientv3.OpPut("/cron/jobs/job8", "123s")
	// 执行op
	if opResp, err = client.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}
	// 打印putResp
	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	// 创建 get op
	getOp = clientv3.OpGet("/cron/jobs/job8")
	// 执行op
	if opResp, err = client.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}
	// 打印getResp
	fmt.Println("数据revision:", opResp.Get().Kvs[0].CreateRevision, opResp.Get().Kvs[0].ModRevision)
	fmt.Println("value:", string(opResp.Get().Kvs[0].Value))

}
func etcdWatch() {
	// 用于读取etcd键值对
	kv = clientv3.NewKV(client)

	// 模拟etcd中KV变化
	go func() {
		for {
			if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job7", "Hello Job7"); err != nil {
				fmt.Println(err)
				return
			}
			if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job7"); err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
	// 先监听当前值,并监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}
	// 判断现在的key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值:", string(getResp.Kvs[0].Value))
	}

	// 当前etcd集群事务ID,单调递增
	watchStartRevision = getResp.Header.Revision
	// 创建一个watcher
	watcher = clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从版本后监听:", watchStartRevision)
	// 设置监听到期时间(可以通过context设置监听时间)
	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})
	watchRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	//处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改值为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
	fmt.Println("Over")
}
func etcdLease() {
	// 租约实例
	lease = clientv3.NewLease(client)
	// 申请一个10s的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	// 获取租约ID
	leaseId = leaseGrantResp.ID
	// 租约续租
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}
	// 处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else {
					fmt.Println("收到自动续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()
	// Put一个KV,并关联租约,从而实现10s后失效
	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "Hello", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("写入并关联租约成功.")
	}

	// 定时查看键是否到期
FORBREAK:
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
			fmt.Println(err)
		}
		if getResp.Count == 0 {
			fmt.Println("kv已过期了")
			break FORBREAK
		}
		fmt.Println("还未过期：", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Over")
}
func etcdDel() {
	// 用于读取etcd键值对
	kv = clientv3.NewKV(client)
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrevKV(), clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("delete success...")
		if len(delResp.PrevKvs) != 0 {
			for _, kvpair = range delResp.PrevKvs {
				fmt.Println("删除:", string(kvpair.Key), string(kvpair.Value))
			}
		}
	}

}
func etcdGetByDir() {
	// 用于读取etcd键值对
	kv = clientv3.NewKV(client)
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs)
	}
}
func etcdGet() {
	// 用于读取etcd键值对
	kv = clientv3.NewKV(client)
	// clientv3.WithCountOnly() 只返回Count字段返回
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1", /*clientv3.WithCountOnly()*/); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs, getResp.Count)
	}
}

func etcdPut() {
	// 用于读取etcd键值对
	kv = clientv3.NewKV(client)
	// 1. 指定超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if putResp, err = kv.Put(ctx, "/cron/jobs/job1", "hello1", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		// 版本
		fmt.Println(putResp.Header.Revision)
		// 上一个键值
		if putResp.PrevKv != nil {
			fmt.Println("PreValue:", string(putResp.PrevKv.Value))
		}
	}
	// 2. 直接设置
	/*if putResp, err = kv.Put(context.TODO1(), "/cron/jobs/job1", "hello"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(putResp.Header.Revision)
	}*/
}
