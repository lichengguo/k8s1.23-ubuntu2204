import redis
from redis.sentinel import Sentinel

# 连接哨兵服务器(主机名也可以用域名)
sentinel = Sentinel([('10.0.1.51', 26379), ('10.0.1.52', 26379), ('10.0.1.53', 26379)], socket_timeout=0.5)

# 获取主服务器地址
master = sentinel.discover_master('mymaster')
print(master)


# 获取从服务器地址
slave = sentinel.discover_slaves('mymaster')
print(slave)

# 获取主服务器进行写入
master = sentinel.master_for('mymaster', socket_timeout=0.5, password='123456', db=15)
w_ret = master.set('foo', 'bar')

# 获取从服务器进行读取（默认是round-roubin）
slave = sentinel.slave_for('mymaster', socket_timeout=0.5, password='123456', db=15)
r_ret = slave.get('foo')
print(r_ret)
