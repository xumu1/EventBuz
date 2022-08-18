# 2022-0818 init
1. create some component of `eventbuz` engine including `bus` `bus publisher` `bus subscriber` `bus controller` interface.
2. create `EventBuz` that  is the implementation of `bus`.
3. add base test case and run pass it.

# 2022-0819 add feature
1. add feature: you can use a function as event handler right now rather than initialize a struct which implement `EventHandler`.
2. using reflect package to  verify the equality between with two `event handler`