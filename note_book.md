# Note
1.当使用map存储时，key要求是可以计算hash值的，而在go中，map、slice、channel和func类型都是不能计算hash值的，因此当使用这些类型作为map的key时会发生panic，而当key是复合类型，并且包含了这四种类型，也会发生panic。

# Q&A
Q1：如果想存function 作为 handler 使用，如何解决不能在接口中存值，又不能将接口 hash 用 map 的方式挂一个setting。